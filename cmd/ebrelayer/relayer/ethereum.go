package relayer

// DONTCOVER

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"log"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/internal/symbol_translator"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ctypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	"go.uber.org/zap"

	"github.com/Sifchain/sifnode/cmd/ebrelayer/contract"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/txs"
	"github.com/Sifchain/sifnode/cmd/ebrelayer/types"
	ethbridgetypes "github.com/Sifchain/sifnode/x/ethbridge/types"
	oracletypes "github.com/Sifchain/sifnode/x/oracle/types"
)

const (
	transactionInterval = 10 * time.Second
	trailingBlocks      = 50
	ethereumWakeupTimer = 60
	maxQueryBlocks      = 5000
)

// EthereumSub is an Ethereum listener that can relay txs to Cosmos and Ethereum
type EthereumSub struct {
	EthProvider             string
	TmProvider              string
	RegistryContractAddress common.Address
	ValidatorName           string
	ValidatorAddress        sdk.ValAddress
	CliCtx                  client.Context
	PrivateKey              *ecdsa.PrivateKey
	SugaredLogger           *zap.SugaredLogger
}

// NewKeybase create a new keybase instance
func NewKeybase(validatorMoniker, mnemonic, password string) (keyring.Keyring, keyring.Info, error) {
	kr := keyring.NewInMemory()
	hdpath := *hd.NewFundraiserParams(0, sdk.CoinType, 0)
	info, err := kr.NewAccount(validatorMoniker, mnemonic, password, hdpath.String(), hd.Secp256k1)
	if err != nil {
		return nil, nil, err
	}

	return kr, info, nil
}

// NewEthereumSub initializes a new EthereumSub
func NewEthereumSub(
	cliCtx client.Context,
	nodeURL string,
	validatorMoniker,
	ethProvider string,
	registryContractAddress common.Address,
	validatorAddress sdk.ValAddress,
	sugaredLogger *zap.SugaredLogger,
) EthereumSub {

	return EthereumSub{
		EthProvider:             ethProvider,
		TmProvider:              nodeURL,
		RegistryContractAddress: registryContractAddress,
		ValidatorName:           validatorMoniker,
		ValidatorAddress:        validatorAddress,
		CliCtx:                  cliCtx,
		SugaredLogger:           sugaredLogger,
	}
}

// Start an Ethereum chain subscription
func (sub EthereumSub) Start(txFactory tx.Factory,
	completionEvent *sync.WaitGroup,
	symbolTranslator *symbol_translator.SymbolTranslator) {

	defer completionEvent.Done()
	time.Sleep(time.Second)
	ethClient, err := SetupWebsocketEthClient(sub.EthProvider)
	if err != nil {
		sub.SugaredLogger.Errorw("SetupWebsocketEthClient failed.",
			errorMessageKey, err.Error())

		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent, symbolTranslator)
		return
	}
	defer ethClient.Close()
	sub.SugaredLogger.Infow("Started Ethereum websocket with provider:",
		"Ethereum provider", sub.EthProvider)

	tmClient, err := tmclient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		return
	}

	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get network ID.",
			errorMessageKey, err.Error())
		completionEvent.Add(1)
		go sub.Start(txFactory, completionEvent, symbolTranslator)
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
	}

	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)

	// start the timer
	t := time.NewTicker(time.Second * ethereumWakeupTimer)
	for {
		select {
		// Handle any errors
		case <-quit:
			return
		case <-t.C:
			sub.CheckNonceAndProcess(txFactory,
				networkID,
				ethClient,
				tmClient,
				bridgeBankAddress,
				bridgeBankContractABI,
				symbolTranslator)
		}
	}
}

// CheckNonceAndProcess check the lock burn nonce and process the event
func (sub EthereumSub) CheckNonceAndProcess(txFactory tx.Factory,
	networkID *big.Int,
	ethClient *ethclient.Client,
	tmClient *tmclient.HTTP,
	bridgeBankAddress common.Address,
	bridgeBankContractABI abi.ABI,
	symbolTranslator *symbol_translator.SymbolTranslator) {
	// get current block height
	currentBlock, err := ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the current block from ethereum client",
			errorMessageKey, err.Error())
		return
	}
	currentBlockHeight := currentBlock.Number

	// If current block is less than 50, just return.
	if currentBlockHeight.Cmp(big.NewInt(trailingBlocks)) <= 0 {
		return
	}

	var endBlockHeight *big.Int
	endBlockHeight = endBlockHeight.Sub(currentBlockHeight, big.NewInt(trailingBlocks))

	// get lock burn nonce from cosmos
	lockBurnNonce, err := sub.GetLockBurnNonceFromCosmos(oracletypes.NetworkDescriptor(networkID.Uint64()), string(sub.ValidatorAddress))
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get the lock burn nonce from cosmos rpc",
			errorMessageKey, err.Error())
		return
	}

	topics := [][]common.Hash{}
	// add the log type as first topic, the first search filter will be lockTopic or burnTopic
	lockTopic := bridgeBankContractABI.Events[types.LogLock.String()].ID()
	burnTopic := bridgeBankContractABI.Events[types.LogBurn.String()].ID()
	topics = append(topics, []common.Hash{lockTopic, burnTopic})

	// add the lock burn nonce as second topic, combined search filter will be (lockTopic or burnTopic) and lockBurnNonceTopic
	var lockBurnNonceTopic [32]byte
	copy(lockBurnNonceTopic[:], abi.U256(big.NewInt(int64(lockBurnNonce + 1)))[:32])
	topics = append(topics, []common.Hash{lockBurnNonceTopic})

	// query the exact block number with the lock burn nonce
	ethLogs, err := ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   endBlockHeight,
		Addresses: []common.Address{bridgeBankAddress},
		Topics:    topics,
	})

	if err != nil {
		sub.SugaredLogger.Errorw("failed to filter the logs from ethereum client",
			errorMessageKey, err.Error())
		return
	}

	fromBlockNumber := uint64(0)
	if len(ethLogs) != 1 {
		sub.SugaredLogger.Errorw("the result from filter is wrong.")
		return
	}

	event, isBurnLock, err := sub.logToEvent(oracletypes.NetworkDescriptor(networkID.Uint64()),
		bridgeBankAddress,
		bridgeBankContractABI,
		ethLogs[0])

	if err != nil {
		sub.SugaredLogger.Errorw("failed to transform from log to event.",
			errorMessageKey, err.Error())
		return
	}
	if !isBurnLock {
		sub.SugaredLogger.Infow("not burn or lock event, continue events.")
		return
	}

	if event.Nonce.Uint64() != lockBurnNonce+1 {
		sub.SugaredLogger.Errorw("the lock burn nonce is not expected.")
		return
	}

	// get the block height for the specific lock burn nonce
	fromBlockNumber = ethLogs[0].BlockNumber

	events := []types.EthereumEvent{}
	// get a new topics, exclude the lock burn nonce since we already get block number
	topics = [][]common.Hash{}
	topics = append(topics, []common.Hash{lockTopic, burnTopic})

	for {
		endBlock := endBlockHeight.Uint64()
		if fromBlockNumber > endBlock {
			break
		}

		// query block scope limited to maxQueryBlocks
		if endBlock > fromBlockNumber+maxQueryBlocks {
			endBlock = fromBlockNumber + maxQueryBlocks
		}

		// query the events with block scope
		ethLogs, err = ethClient.FilterLogs(context.Background(), ethereum.FilterQuery{
			FromBlock: big.NewInt(int64(fromBlockNumber)),
			ToBlock:   big.NewInt(int64(endBlock)),
			Addresses: []common.Address{bridgeBankAddress},
			Topics:    topics,
		})

		if err != nil {
			sub.SugaredLogger.Errorw("failed to filter the logs from ethereum client",
				errorMessageKey, err.Error())
			return
		}

		// loop over ethlogs, and build an array of burn/lock events
		for _, ethLog := range ethLogs {
			event, isBurnLock, err := sub.logToEvent(oracletypes.NetworkDescriptor(networkID.Uint64()),
				bridgeBankAddress,
				bridgeBankContractABI,
				ethLog)

			if err != nil {
				sub.SugaredLogger.Errorw("failed to transform from log to event.",
					errorMessageKey, err.Error())
				continue
			}
			if !isBurnLock {
				sub.SugaredLogger.Infow("not burn or lock event, continue events.")
				continue
			}
			events = append(events, event)
		}

		if len(events) > 0 {
			var nextLockBurnNonce uint64
			if nextLockBurnNonce, err = sub.handleEthereumEvent(txFactory, events, symbolTranslator, lockBurnNonce); err != nil {
				sub.SugaredLogger.Errorw("failed to handle ethereum event.",
					errorMessageKey, err.Error())
				return
			}
			// handleEthereumEvent return the next expected lock burn nonce
			lockBurnNonce = nextLockBurnNonce - 1
			time.Sleep(transactionInterval)
		}

		// update fromBlockNumber
		fromBlockNumber += maxQueryBlocks
	}
}

// Replay the missed events
func (sub EthereumSub) Replay(txFactory tx.Factory, symbolTranslator *symbol_translator.SymbolTranslator) {

	ethClient, err := SetupWebsocketEthClient(sub.EthProvider)
	if err != nil {
		sub.SugaredLogger.Errorw("SetupWebsocketEthClient failed.",
			errorMessageKey, err.Error())

		return
	}
	defer ethClient.Close()
	sub.SugaredLogger.Infow("Started Ethereum websocket with provider:",
		"Ethereum provider", sub.EthProvider)

	tmClient, err := tmclient.New(sub.TmProvider, "/websocket")
	if err != nil {
		sub.SugaredLogger.Errorw("failed to initialize a sifchain client.",
			errorMessageKey, err.Error())
		return
	}

	networkID, err := ethClient.NetworkID(context.Background())
	if err != nil {
		sub.SugaredLogger.Errorw("failed to get network ID.",
			errorMessageKey, err.Error())

		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer close(quit)

	// get the bridgebank address from the registry contract
	bridgeBankAddress, err := txs.GetAddressFromBridgeRegistry(ethClient, sub.RegistryContractAddress, txs.BridgeBank, sub.SugaredLogger)
	if err != nil {
		log.Fatal("Error getting bridgebank address: ", err.Error())
	}

	bridgeBankContractABI := contract.LoadABI(txs.BridgeBank)

	sub.CheckNonceAndProcess(txFactory,
		networkID,
		ethClient,
		tmClient,
		bridgeBankAddress,
		bridgeBankContractABI,
		symbolTranslator)
}

// logToEvent unpacks an Ethereum event
func (sub EthereumSub) logToEvent(networkDescriptor oracletypes.NetworkDescriptor, contractAddress common.Address,
	contractABI abi.ABI, cLog ctypes.Log) (types.EthereumEvent, bool, error) {
	// Parse the event's attributes via contract ABI
	event := types.EthereumEvent{}
	eventLogLockSignature := contractABI.Events[types.LogLock.String()].ID().Hex()
	eventLogBurnSignature := contractABI.Events[types.LogBurn.String()].ID().Hex()

	var eventName string
	switch cLog.Topics[0].Hex() {
	case eventLogBurnSignature:
		eventName = types.LogBurn.String()
	case eventLogLockSignature:
		eventName = types.LogLock.String()
	default:
		eventName = ""
	}

	// If event is not expected
	if eventName == "" {
		return event, false, nil
	}

	err := contractABI.Unpack(&event, eventName, cLog.Data)
	if err != nil {
		sub.SugaredLogger.Errorw(".",
			errorMessageKey, err.Error())
		return event, false, err
	}
	event.BridgeContractAddress = contractAddress
	event.NetworkDescriptor = networkDescriptor
	if eventName == types.LogBurn.String() {
		event.ClaimType = ethbridgetypes.ClaimType_CLAIM_TYPE_BURN
	} else {
		event.ClaimType = ethbridgetypes.ClaimType_CLAIM_TYPE_LOCK
	}
	sub.SugaredLogger.Infow("receive an event.",
		"event", event)

	// Add the event to the record
	types.NewEventWrite(cLog.TxHash.Hex(), event)
	return event, true, nil
}

// handleEthereumEvent unpacks an Ethereum event, converts it to a ProphecyClaim, and relays a tx to Cosmos
func (sub EthereumSub) handleEthereumEvent(txFactory tx.Factory,
	events []types.EthereumEvent,
	symbolTranslator *symbol_translator.SymbolTranslator,
	lockBurnNonce uint64) (uint64, error) {

	var prophecyClaims []*ethbridgetypes.EthBridgeClaim
	nextLockBurnNonce := lockBurnNonce + 1
	valAddr, err := GetValAddressFromKeyring(txFactory.Keybase(), sub.ValidatorName)
	if err != nil {
		return nextLockBurnNonce, err
	}
	for _, event := range events {
		prophecyClaim, err := txs.EthereumEventToEthBridgeClaim(valAddr, event, symbolTranslator, sub.SugaredLogger)
		if err != nil {
			sub.SugaredLogger.Errorw(".",
				errorMessageKey, err.Error())
		} else {
			if prophecyClaim.EthereumLockBurnNonce == nextLockBurnNonce {
				prophecyClaims = append(prophecyClaims, &prophecyClaim)
				nextLockBurnNonce++
			} else {
				sub.SugaredLogger.Infow("ock burn nonce is not expected.",
					"expected lock burn nonce is %d", nextLockBurnNonce,
					"lock burn nonce from event is %d", prophecyClaim.EthereumLockBurnNonce)
				return nextLockBurnNonce, errors.New("lock burn nonce is not expected")
			}

		}
	}
	sub.SugaredLogger.Infow("relay prophecy claims to cosmos.",
		"prophecy claims length", len(prophecyClaims))

	if len(events) == 0 {
		return nextLockBurnNonce, nil
	}

	return nextLockBurnNonce, txs.RelayToCosmos(txFactory, prophecyClaims, sub.CliCtx, sub.SugaredLogger)
}

// GetLockBurnNonceFromCosmos via rpc
func (sub EthereumSub) GetLockBurnNonceFromCosmos(
	networkDescriptor oracletypes.NetworkDescriptor,
	relayerValAddress string) (uint64, error) {
	conn, err := grpc.Dial(sub.TmProvider)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	client := ethbridgetypes.NewQueryClient(conn)
	request := ethbridgetypes.QueryLockBurnNonceRequest{
		NetworkDescriptor: networkDescriptor,
		RelayerValAddress: relayerValAddress,
	}
	response, err := client.LockBurnNonce(ctx, &request)
	if err != nil {
		return 0, err
	}
	return response.LockBurnNonce, nil
}

// GetValAddressFromKeyring get validator address from keyring
func GetValAddressFromKeyring(k keyring.Keyring, keyname string) (sdk.ValAddress, error) {
	keyInfo, err := k.Key(keyname)
	if err != nil {
		return nil, err
	}
	return sdk.ValAddress(keyInfo.GetAddress()), nil
}
