package keeper

import (
	"context"

	"fmt"
	"math"
	"strconv"
	"strings"

	admintypes "github.com/Sifchain/sifnode/x/admin/types"
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/Sifchain/sifnode/x/clp/types"
	tokenregistrytypes "github.com/Sifchain/sifnode/x/tokenregistry/types"
)

type msgServer struct {
	Keeper
}

func (k msgServer) SetSymmetryThreshold(goCtx context.Context, threshold *types.MsgSetSymmetryThreshold) (*types.MsgSetSymmetryThresholdResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(threshold.Signer)
	if err != nil {
		return nil, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_CLPDEX, signer) {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", threshold.Signer))
	}

	k.Keeper.SetSymmetryThreshold(sdk.UnwrapSDKContext(goCtx), threshold)

	return &types.MsgSetSymmetryThresholdResponse{}, nil
}

// NewMsgServerImpl returns an implementation of the clp MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) CancelUnlockLiquidity(goCtx context.Context, request *types.MsgCancelUnlock) (*types.MsgCancelUnlockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	lp, err := k.Keeper.GetLiquidityProvider(ctx, request.ExternalAsset.Symbol, request.Signer)
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}
	// Prune unlocks
	params := k.GetRewardsParams(ctx)
	k.PruneUnlockRecords(ctx, &lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)

	err = k.UseUnlockedLiquidity(ctx, lp, request.Units, true)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCancelUnlock,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyPool, lp.Asset.Symbol),
			sdk.NewAttribute(types.AttributeKeyUnits, request.Units.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, request.Signer),
		),
	})
	return &types.MsgCancelUnlockResponse{}, nil
}

func (k msgServer) UpdateStakingRewardParams(goCtx context.Context, msg *types.MsgUpdateStakingRewardParams) (*types.MsgUpdateStakingRewardParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return nil, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	if !(msg.Minter.AnnualProvisions.IsZero() && msg.Minter.Inflation.IsZero()) {
		k.mintKeeper.SetMinter(ctx, msg.Minter)
	}
	k.mintKeeper.SetParams(ctx, msg.Params)

	return &types.MsgUpdateStakingRewardParamsResponse{}, err

}

func (k msgServer) UpdateRewardsParams(goCtx context.Context, msg *types.MsgUpdateRewardsParamsRequest) (*types.MsgUpdateRewardsParamsResponse, error) {
	response := &types.MsgUpdateRewardsParamsResponse{}
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	params := k.GetRewardsParams(ctx)
	params.LiquidityRemovalLockPeriod = msg.LiquidityRemovalLockPeriod
	params.LiquidityRemovalCancelPeriod = msg.LiquidityRemovalCancelPeriod
	k.SetRewardParams(ctx, params)
	return response, err
}

func (k msgServer) AddRewardPeriod(goCtx context.Context, msg *types.MsgAddRewardPeriodRequest) (*types.MsgAddRewardPeriodResponse, error) {
	response := &types.MsgAddRewardPeriodResponse{}
	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	params := k.GetRewardsParams(ctx)
	params.RewardPeriods = msg.RewardPeriods
	k.SetRewardParams(ctx, params)
	return response, nil
}

func (k msgServer) AddProviderDistributionPeriod(goCtx context.Context, msg *types.MsgAddProviderDistributionPeriodRequest) (*types.MsgAddProviderDistributionPeriodResponse, error) {
	response := &types.MsgAddProviderDistributionPeriodResponse{}

	// defensive programming
	if msg == nil {
		return response, errors.Errorf("msg was nil")
	}

	if err := msg.ValidateBasic(); err != nil {
		return response, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}

	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}

	params := &types.ProviderDistributionParams{}
	params.DistributionPeriods = msg.DistributionPeriods
	k.SetProviderDistributionParams(ctx, params)

	eventMsg := CreateEventMsg(msg.Signer)
	attribute := sdk.NewAttribute(types.AttributeKeyProviderDistributionParams, params.String())
	providerDistributionPolicyEvent := CreateEventBlockHeight(ctx, types.EventTypeAddNewProviderDistributionPolicy, attribute)
	ctx.EventManager().EmitEvents(sdk.Events{providerDistributionPolicyEvent, eventMsg})

	return response, nil
}

func (k msgServer) UpdatePmtpParams(goCtx context.Context, msg *types.MsgUpdatePmtpParams) (*types.MsgUpdatePmtpParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	response := &types.MsgUpdatePmtpParamsResponse{}
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	params := k.GetPmtpParams(ctx)
	// Check to see if a policy is still running
	if k.IsInsidePmtpWindow(ctx) {
		return response, types.ErrCannotStartPolicy
	}
	// Check to make sure new policy starts in the future so that PolicyStart from begin-block can be triggered
	if msg.PmtpPeriodStartBlock <= ctx.BlockHeight() {
		return response, errors.New("Start block cannot be in the past/current block")
	}
	params.PmtpPeriodStartBlock = msg.PmtpPeriodStartBlock
	params.PmtpPeriodEndBlock = msg.PmtpPeriodEndBlock
	params.PmtpPeriodEpochLength = msg.PmtpPeriodEpochLength

	if !strings.EqualFold(msg.PmtpPeriodGovernanceRate, "") {
		rGov, err := sdk.NewDecFromStr(msg.PmtpPeriodGovernanceRate)
		if err != nil {
			return response, err
		}
		params.PmtpPeriodGovernanceRate = rGov
	}

	k.SetPmtpParams(ctx, params)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddNewPmtpPolicy,
			sdk.NewAttribute(types.AttributeKeyPmtpPolicyParams, params.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgUpdatePmtpParamsResponse{}, nil
}

func (k msgServer) ModifyPmtpRates(goCtx context.Context, msg *types.MsgModifyPmtpRates) (*types.MsgModifyPmtpRatesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	response := &types.MsgModifyPmtpRatesResponse{}
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	params := k.GetPmtpParams(ctx)
	rateParams := k.GetPmtpRateParams(ctx)

	// Set Block Rate is needed only if no policy is presently executing
	if !strings.EqualFold(msg.BlockRate, "") && !k.IsInsidePmtpWindow(ctx) {
		blockRate, err := sdk.NewDecFromStr(msg.BlockRate)
		if err != nil {
			return response, err
		}
		rateParams.PmtpPeriodBlockRate = blockRate
	}

	// Set Running Rate if Needed only if no policy is presently executing
	if !strings.EqualFold(msg.RunningRate, "") && !k.IsInsidePmtpWindow(ctx) {
		runningRate, err := sdk.NewDecFromStr(msg.RunningRate)
		if err != nil {
			return response, err
		}
		rateParams.PmtpCurrentRunningRate = runningRate
		// inter policy rate should always equal running rate between policies
		rateParams.PmtpInterPolicyRate = runningRate
	}
	k.SetPmtpRateParams(ctx, rateParams)
	events := sdk.EmptyEvents()
	// End Policy If Needed , returns if not policy is presently
	if msg.EndPolicy && k.IsInsidePmtpWindow(ctx) {
		params.PmtpPeriodEndBlock = ctx.BlockHeight()
		k.SetPmtpParams(ctx, params)
		k.SetPmtpEpoch(ctx, types.PmtpEpoch{
			EpochCounter: 0,
			BlockCounter: 0,
		})
		k.SetPmtpInterPolicyRate(ctx, rateParams.PmtpCurrentRunningRate)
		events = events.AppendEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeEndPmtpPolicy,
				sdk.NewAttribute(types.AttributeKeyPmtpPolicyParams, params.String()),
				sdk.NewAttribute(types.AttributeKeyPmtpRateParams, rateParams.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
		})
	}
	ctx.EventManager().EmitEvents(events)
	return response, nil
}

func (k msgServer) UnlockLiquidity(goCtx context.Context, request *types.MsgUnlockLiquidityRequest) (*types.MsgUnlockLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	lp, err := k.Keeper.GetLiquidityProvider(ctx, request.ExternalAsset.Symbol, request.Signer)
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}
	// Prune unlocks
	params := k.GetRewardsParams(ctx)
	k.PruneUnlockRecords(ctx, &lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)
	totalUnlocks := sdk.ZeroUint()
	for _, unlock := range lp.Unlocks {
		totalUnlocks = totalUnlocks.Add(unlock.Units)
	}
	if totalUnlocks.Add(request.Units).GT(lp.LiquidityProviderUnits) {
		return nil, types.ErrBalanceNotAvailable
	}
	lp.Unlocks = append(lp.Unlocks, &types.LiquidityUnlock{
		RequestHeight: ctx.BlockHeight(),
		Units:         request.Units,
	})
	k.Keeper.SetLiquidityProvider(ctx, &lp)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRequestUnlock,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyPool, lp.Asset.Symbol),
			sdk.NewAttribute(types.AttributeKeyUnits, request.Units.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, request.Signer),
		),
	})
	return &types.MsgUnlockLiquidityResponse{}, nil
}

func (k msgServer) DecommissionPool(goCtx context.Context, msg *types.MsgDecommissionPool) (*types.MsgDecommissionPoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, err := k.Keeper.GetPool(ctx, msg.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	addAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	// TODO : Deprecate this Admin in favor of TokenRegistry
	if !k.Keeper.ValidateAddress(ctx, addAddr) {
		return nil, sdkerrors.Wrap(types.ErrInvalid, "user does not have permission to decommission pool")
	}
	if pool.NativeAssetBalance.GTE(sdk.NewUintFromString(types.PoolThrehold)) {
		return nil, types.ErrBalanceTooHigh
	}
	// Get all LP's for the pool
	if pool.ExternalAsset == nil {
		return nil, errors.New("nill external asset")
	}
	lpList, _, err := k.Keeper.GetLiquidityProvidersForAssetPaginated(ctx, *pool.ExternalAsset, &query.PageRequest{
		Limit: uint64(math.MaxUint64),
	})
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrLiquidityProviderDoesNotExist, err.Error())
	}
	poolUnits := pool.PoolUnits
	nativeAssetBalance := pool.NativeAssetBalance
	externalAssetBalance := pool.ExternalAssetBalance
	// iterate over Lp list and refund them there tokens
	// Return both RWN and EXTERNAL ASSET
	for _, lp := range lpList {
		withdrawNativeAsset, withdrawExternalAsset, _, _ := CalculateAllAssetsForLP(pool, *lp)
		poolUnits = poolUnits.Sub(lp.LiquidityProviderUnits)
		nativeAssetBalance = nativeAssetBalance.Sub(withdrawNativeAsset)
		externalAssetBalance = externalAssetBalance.Sub(withdrawExternalAsset)

		withdrawNativeAssetInt, ok := k.Keeper.ParseToInt(withdrawNativeAsset.String())
		if !ok {
			return nil, types.ErrUnableToParseInt
		}
		withdrawExternalAssetInt, ok := k.Keeper.ParseToInt(withdrawExternalAsset.String())
		if !ok {
			return nil, types.ErrUnableToParseInt
		}
		withdrawNativeCoins := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetInt)
		withdrawExternalCoins := sdk.NewCoin(msg.Symbol, withdrawExternalAssetInt)
		refundingCoins := sdk.NewCoins(withdrawExternalCoins, withdrawNativeCoins)
		err := k.Keeper.RemoveLiquidityProvider(ctx, refundingCoins, *lp)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidityProvider, err.Error())
		}
	}
	// Pool should be empty at this point
	// Decommission the pool
	err = k.Keeper.DecommissionPool(ctx, pool)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToDecommissionPool, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDecommissionPool,
			sdk.NewAttribute(types.AttributeKeyPool, pool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgDecommissionPoolResponse{}, nil
}

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePool) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Verify min threshold
	MinThreshold := sdk.NewUintFromString(types.PoolThrehold)
	if msg.NativeAssetAmount.LT(MinThreshold) { // Need to verify
		return nil, types.ErrTotalAmountTooLow
	}
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(eAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}

	// Check if pool already exists
	if k.Keeper.ExistsPool(ctx, msg.ExternalAsset.Symbol) {
		return nil, types.ErrUnableToCreatePool
	}

	nativeBalance := msg.NativeAssetAmount
	externalBalance := msg.ExternalAssetAmount
	externalDecimals, err := Int64ToUint8Safe(eAsset.Decimals)
	if err != nil {
		return nil, err
	}

	poolUnits, lpunits, err := CalculatePoolUnits(sdk.ZeroUint(), sdk.ZeroUint(), sdk.ZeroUint(),
		nativeBalance, externalBalance, externalDecimals, k.GetSymmetryThreshold(ctx), k.GetSymmetryRatio(ctx))
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToCreatePool, err.Error())
	}

	pool, err := k.Keeper.CreatePool(ctx, poolUnits, msg)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
	}
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	lp := k.Keeper.CreateLiquidityProvider(ctx, msg.ExternalAsset, lpunits, accAddr)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreatePool,
			sdk.NewAttribute(types.AttributeKeyPool, pool.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			types.EventTypeCreateLiquidityProvider,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgCreatePoolResponse{}, nil
}

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var (
		priceImpact sdk.Uint
	)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	sAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.SentAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	rAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ReceivedAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(sAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(rAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	if k.tokenRegistryKeeper.CheckEntryPermissions(sAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_DISABLE_SELL}) {
		return nil, tokenregistrytypes.ErrNotAllowedToSellAsset
	}
	if k.tokenRegistryKeeper.CheckEntryPermissions(rAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_DISABLE_BUY}) {
		return nil, tokenregistrytypes.ErrNotAllowedToBuyAsset
	}

	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	swapFeeRate := k.GetSwapFeeRate(ctx).SwapFeeRate

	liquidityProtectionParams := k.GetLiquidityProtectionParams(ctx)
	maxRowanLiquidityThreshold := liquidityProtectionParams.MaxRowanLiquidityThreshold
	maxRowanLiquidityThresholdAsset := liquidityProtectionParams.MaxRowanLiquidityThresholdAsset
	currentRowanLiquidityThreshold := k.GetLiquidityProtectionRateParams(ctx).CurrentRowanLiquidityThreshold
	var (
		sentValue sdk.Uint
	)

	// if liquidity protection is active and selling rowan
	if liquidityProtectionParams.IsActive && strings.EqualFold(sAsset.Denom, types.NativeSymbol) {
		if strings.EqualFold(maxRowanLiquidityThresholdAsset, types.NativeSymbol) {
			sentValue = msg.SentAmount
		} else {
			pool, err := k.GetPool(ctx, maxRowanLiquidityThresholdAsset)
			if err != nil {
				return nil, types.ErrMaxRowanLiquidityThresholdAssetPoolDoesNotExist
			}
			sentValue, err = CalcRowanValue(&pool, pmtpCurrentRunningRate, msg.SentAmount)

			if err != nil {
				return nil, err
			}
		}

		if currentRowanLiquidityThreshold.LT(sentValue) {
			return nil, types.ErrReachedMaxRowanLiquidityThreshold
		}
	}

	liquidityFeeNative := sdk.ZeroUint()
	liquidityFeeExternal := sdk.ZeroUint()
	totalLiquidityFee := sdk.ZeroUint() // nolint:staticcheck
	priceImpact = sdk.ZeroUint()
	sentAmount := msg.SentAmount
	sentAsset := msg.SentAsset
	receivedAsset := msg.ReceivedAsset
	// Get native asset
	nativeAsset := types.GetSettlementAsset()
	inPool, outPool := types.Pool{}, types.Pool{}
	// If sending rowan ,deduct directly from the Native balance  instead of fetching from rowan pool
	if !msg.SentAsset.Equals(types.GetSettlementAsset()) {
		inPool, err = k.Keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	}
	sentAmountInt, ok := k.Keeper.ParseToInt(sentAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	accAddr, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return nil, err
	}
	sentCoin := sdk.NewCoin(msg.SentAsset.Symbol, sentAmountInt)
	err = k.Keeper.InitiateSwap(ctx, sentCoin, accAddr)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	// Check if its a two way swap, swapping non native fro non native .
	// If its one way we can skip this if condition and add balance to users account from outpool
	if !msg.SentAsset.Equals(nativeAsset) && !msg.ReceivedAsset.Equals(nativeAsset) {
		emitAmount, lp, ts, finalPool, err := FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx, k, *sentAsset, sentAmount, nativeAsset, inPool, pmtpCurrentRunningRate, swapFeeRate)
		if err != nil {
			return nil, err
		}
		err = k.Keeper.SetPool(ctx, &finalPool)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSetPool, err.Error())
		}
		sentAmount = emitAmount
		sentAsset = &nativeAsset
		priceImpact = priceImpact.Add(ts)
		liquidityFeeNative = liquidityFeeNative.Add(lp)
	}
	// If receiving  rowan , add directly to  Native balance  instead of fetching from rowan pool
	if msg.ReceivedAsset.Equals(types.GetSettlementAsset()) {
		outPool, err = k.Keeper.GetPool(ctx, msg.SentAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.SentAsset.String())
		}
	} else {
		outPool, err = k.Keeper.GetPool(ctx, msg.ReceivedAsset.Symbol)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrPoolDoesNotExist, msg.ReceivedAsset.String())
		}
	}
	// Calculating amount user receives
	emitAmount, lp, ts, finalPool, err := FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx, k, *sentAsset, sentAmount, *receivedAsset, outPool, pmtpCurrentRunningRate, swapFeeRate)
	if err != nil {
		return nil, err
	}
	if emitAmount.LT(msg.MinReceivingAmount) {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeSwapFailed,
				sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
				sdk.NewAttribute(types.AttributeKeyThreshold, msg.MinReceivingAmount.String()),
				sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
				sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
				sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
			),
			sdk.NewEvent(
				sdk.EventTypeMessage,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
			),
		})
		return &types.MsgSwapResponse{}, types.ErrReceivedAmountBelowExpected
	}
	// todo nil pointer deref test
	err = k.Keeper.FinalizeSwap(ctx, emitAmount.String(), finalPool, *msg)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
	}
	if liquidityFeeNative.GT(sdk.ZeroUint()) {
		liquidityFeeExternal = liquidityFeeExternal.Add(lp)
		firstSwapFeeInOutputAsset := FEATURE_TOGGLE_MARGIN_CLI_ALPHA_GetSwapFee(ctx, k, msg.ReceivedAsset, liquidityFeeNative, outPool, pmtpCurrentRunningRate, swapFeeRate)
		totalLiquidityFee = liquidityFeeExternal.Add(firstSwapFeeInOutputAsset)
	} else {
		totalLiquidityFee = liquidityFeeNative.Add(lp)
	}
	priceImpact = priceImpact.Add(ts)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(types.AttributeKeySwapAmount, emitAmount.String()),
			sdk.NewAttribute(types.AttributeKeyLiquidityFee, totalLiquidityFee.String()),
			sdk.NewAttribute(types.AttributeKeyPriceImpact, priceImpact.String()),
			sdk.NewAttribute(types.AttributeKeyInPool, inPool.String()),
			sdk.NewAttribute(types.AttributeKeyOutPool, outPool.String()),
			sdk.NewAttribute(types.AttributePmtpBlockRate, k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate.String()),
			sdk.NewAttribute(types.AttributePmtpCurrentRunningRate, pmtpCurrentRunningRate.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	if liquidityProtectionParams.IsActive {
		// if sell rowan
		if strings.EqualFold(sAsset.Denom, types.NativeSymbol) {
			// we know that sentValue < currentRowanLiquidityThreshold so we can do the
			// substitution knowing it won't panic
			currentRowanLiquidityThreshold = currentRowanLiquidityThreshold.Sub(sentValue)
			k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, currentRowanLiquidityThreshold)
		}

		// if buy rowan
		if strings.EqualFold(rAsset.Denom, types.NativeSymbol) {
			var emitValue sdk.Uint
			if strings.EqualFold(maxRowanLiquidityThresholdAsset, types.NativeSymbol) {
				emitValue = emitAmount
			} else {
				pool, err := k.GetPool(ctx, maxRowanLiquidityThresholdAsset)
				if err != nil {
					return nil, types.ErrMaxRowanLiquidityThresholdAssetPoolDoesNotExist
				}
				emitValue, err = CalcRowanValue(&pool, pmtpCurrentRunningRate, emitAmount)

				if err != nil {
					return nil, err
				}
			}

			// This is equivalent to currentRowanLiquidityThreshold := sdk.MinUint(currentRowanLiquidityThreshold.Add(emitValue), maxRowanLiquidityThreshold)
			// except it prevents any overflows when adding the emitValue
			if maxRowanLiquidityThreshold.Sub(currentRowanLiquidityThreshold).LT(emitValue) {
				currentRowanLiquidityThreshold = maxRowanLiquidityThreshold
			} else {
				currentRowanLiquidityThreshold = currentRowanLiquidityThreshold.Add(emitValue)
			}

			k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, currentRowanLiquidityThreshold)
		}
	}

	return &types.MsgSwapResponse{}, nil
}

func (k msgServer) AddLiquidity(goCtx context.Context, msg *types.MsgAddLiquidity) (*types.MsgAddLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(eAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	// Get pool
	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}

	externalDecimals, err := Int64ToUint8Safe(eAsset.Decimals)
	if err != nil {
		return nil, err
	}

	newPoolUnits, lpUnits, err := CalculatePoolUnits(
		pool.PoolUnits,
		pool.NativeAssetBalance,
		pool.ExternalAssetBalance,
		msg.NativeAssetAmount,
		msg.ExternalAssetAmount,
		externalDecimals,
		k.GetSymmetryThreshold(ctx),
		k.GetSymmetryRatio(ctx))
	if err != nil {
		return nil, err
	}
	// Get lp , if lp doesnt exist create lp
	lp, err := k.Keeper.AddLiquidity(ctx, msg, pool, newPoolUnits, lpUnits)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToAddLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAddLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyUnits, lpUnits.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})

	FEATURE_TOGGLE_MARGIN_CLI_ALPHA_ProcessRemovelQueue(ctx, k, msg, newPoolUnits)

	return &types.MsgAddLiquidityResponse{}, nil
}

func (k msgServer) RemoveLiquidityUnits(goCtx context.Context, msg *types.MsgRemoveLiquidityUnits) (*types.MsgRemoveLiquidityUnitsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(eAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := k.Keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)
	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}

	err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWithdrawUnitsAvailableForLP(ctx, k, msg, lp)
	if err != nil {
		return nil, err
	}

	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	swapFeeRate := k.GetSwapFeeRate(ctx).SwapFeeRate
	// Prune pools
	params := k.GetRewardsParams(ctx)
	k.PruneUnlockRecords(ctx, &lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft := CalculateWithdrawalFromUnits(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WithdrawUnits)

	err = k.Keeper.UseUnlockedLiquidity(ctx, lp, lp.LiquidityProviderUnits.Sub(lpUnitsLeft), false)
	if err != nil {
		return nil, err
	}

	err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWithdrawUnits(ctx, k, msg, lp, pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount, eAsset, pmtpCurrentRunningRate, swapFeeRate)
	if err != nil {
		return nil, err
	}

	withdrawExternalAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	withdrawNativeAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawNativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)
	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)

	// Check and  remove Liquidity
	err = k.Keeper.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyUnits, lp.LiquidityProviderUnits.Sub(lpUnitsLeft).String()),
			sdk.NewAttribute(types.AttributePmtpBlockRate, k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate.String()),
			sdk.NewAttribute(types.AttributePmtpCurrentRunningRate, pmtpCurrentRunningRate.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgRemoveLiquidityUnitsResponse{}, nil
}

func (k msgServer) RemoveLiquidity(goCtx context.Context, msg *types.MsgRemoveLiquidity) (*types.MsgRemoveLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	registry := k.tokenRegistryKeeper.GetRegistry(ctx)
	eAsset, err := k.tokenRegistryKeeper.GetEntry(registry, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrTokenNotSupported
	}
	if !k.tokenRegistryKeeper.CheckEntryPermissions(eAsset, []tokenregistrytypes.Permission{tokenregistrytypes.Permission_CLP}) {
		return nil, tokenregistrytypes.ErrPermissionDenied
	}
	pool, err := k.Keeper.GetPool(ctx, msg.ExternalAsset.Symbol)
	if err != nil {
		return nil, types.ErrPoolDoesNotExist
	}
	//Get LP
	lp, err := k.Keeper.GetLiquidityProvider(ctx, msg.ExternalAsset.Symbol, msg.Signer)

	if err != nil {
		return nil, types.ErrLiquidityProviderDoesNotExist
	}

	poolOriginalEB := pool.ExternalAssetBalance
	poolOriginalNB := pool.NativeAssetBalance
	pmtpCurrentRunningRate := k.GetPmtpRateParams(ctx).PmtpCurrentRunningRate
	swapFeeRate := k.GetSwapFeeRate(ctx).SwapFeeRate
	// Prune pools
	params := k.GetRewardsParams(ctx)
	k.PruneUnlockRecords(ctx, &lp, params.LiquidityRemovalLockPeriod, params.LiquidityRemovalCancelPeriod)

	if !msg.Asymmetry.IsZero() {
		return nil, types.ErrAsymmetricRemove
	}

	err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_VerifyEnoughWBasisPointsAvailableForLP(ctx, k, msg, lp)

	if err != nil {
		return nil, err
	}

	//Calculate amount to withdraw
	withdrawNativeAssetAmount, withdrawExternalAssetAmount, lpUnitsLeft, swapAmount := CalculateWithdrawal(pool.PoolUnits,
		pool.NativeAssetBalance.String(), pool.ExternalAssetBalance.String(), lp.LiquidityProviderUnits.String(),
		msg.WBasisPoints.String(), msg.Asymmetry)

	err = k.Keeper.UseUnlockedLiquidity(ctx, lp, lp.LiquidityProviderUnits.Sub(lpUnitsLeft), false)
	if err != nil {
		return nil, err
	}

	err = FEATURE_TOGGLE_MARGIN_CLI_ALPHA_QueueRemovalWithWBasisPoints(ctx, k, msg, lp, pool, withdrawNativeAssetAmount, withdrawExternalAssetAmount, eAsset, pmtpCurrentRunningRate, swapFeeRate)
	if err != nil {
		return nil, err
	}

	withdrawExternalAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawExternalAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	withdrawNativeAssetAmountInt, ok := k.Keeper.ParseToInt(withdrawNativeAssetAmount.String())
	if !ok {
		return nil, types.ErrUnableToParseInt
	}
	externalAssetCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, withdrawExternalAssetAmountInt)
	nativeAssetCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, withdrawNativeAssetAmountInt)
	// Subtract Value from pool
	pool.PoolUnits = pool.PoolUnits.Sub(lp.LiquidityProviderUnits).Add(lpUnitsLeft)
	pool.NativeAssetBalance = pool.NativeAssetBalance.Sub(withdrawNativeAssetAmount)
	pool.ExternalAssetBalance = pool.ExternalAssetBalance.Sub(withdrawExternalAssetAmount)
	// Check if withdrawal makes pool too shallow , checking only for asymetric withdraw.
	if !msg.Asymmetry.IsZero() && (pool.ExternalAssetBalance.IsZero() || pool.NativeAssetBalance.IsZero()) {
		return nil, sdkerrors.Wrap(types.ErrPoolTooShallow, "pool balance nil before adjusting asymmetry")

	}
	// Swapping between Native and External based on Asymmetry
	if msg.Asymmetry.IsPositive() {
		swapResult, _, _, swappedPool, err := FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx, k, types.GetSettlementAsset(), swapAmount, *msg.ExternalAsset, pool, pmtpCurrentRunningRate, swapFeeRate)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapResultInt, ok := k.Keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.Keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapResultInt)
			swapAmountInCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapAmountInt)
			externalAssetCoin = externalAssetCoin.Add(swapCoin)
			nativeAssetCoin = nativeAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	if msg.Asymmetry.IsNegative() {
		swapResult, _, _, swappedPool, err := FEATURE_TOGGLE_MARGIN_CLI_ALPHA_SwapOne(ctx, k, *msg.ExternalAsset, swapAmount, types.GetSettlementAsset(), pool, pmtpCurrentRunningRate, swapFeeRate)
		if err != nil {
			return nil, sdkerrors.Wrap(types.ErrUnableToSwap, err.Error())
		}
		if !swapResult.IsZero() {
			swapInt, ok := k.Keeper.ParseToInt(swapResult.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapAmountInt, ok := k.Keeper.ParseToInt(swapAmount.String())
			if !ok {
				return nil, types.ErrUnableToParseInt
			}
			swapCoin := sdk.NewCoin(types.GetSettlementAsset().Symbol, swapInt)
			swapAmountInCoin := sdk.NewCoin(msg.ExternalAsset.Symbol, swapAmountInt)
			nativeAssetCoin = nativeAssetCoin.Add(swapCoin)
			externalAssetCoin = externalAssetCoin.Sub(swapAmountInCoin)
		}
		pool = swappedPool
	}
	// Check and  remove Liquidity
	err = k.Keeper.RemoveLiquidity(ctx, pool, externalAssetCoin, nativeAssetCoin, lp, lpUnitsLeft, poolOriginalEB, poolOriginalNB)
	if err != nil {
		return nil, sdkerrors.Wrap(types.ErrUnableToRemoveLiquidity, err.Error())
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRemoveLiquidity,
			sdk.NewAttribute(types.AttributeKeyLiquidityProvider, lp.String()),
			sdk.NewAttribute(types.AttributeKeyUnits, lp.LiquidityProviderUnits.Sub(lpUnitsLeft).String()),
			sdk.NewAttribute(types.AttributePmtpBlockRate, k.GetPmtpRateParams(ctx).PmtpPeriodBlockRate.String()),
			sdk.NewAttribute(types.AttributePmtpCurrentRunningRate, pmtpCurrentRunningRate.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgRemoveLiquidityResponse{}, nil
}

func (k msgServer) UpdateLiquidityProtectionParams(goCtx context.Context, msg *types.MsgUpdateLiquidityProtectionParams) (*types.MsgUpdateLiquidityProtectionParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	response := &types.MsgUpdateLiquidityProtectionParamsResponse{}
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_CLPDEX, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	params := k.GetLiquidityProtectionParams(ctx)
	params.MaxRowanLiquidityThreshold = msg.MaxRowanLiquidityThreshold
	params.MaxRowanLiquidityThresholdAsset = msg.MaxRowanLiquidityThresholdAsset
	params.EpochLength = msg.EpochLength
	params.IsActive = msg.IsActive
	k.SetLiquidityProtectionParams(ctx, params)
	k.SetLiquidityProtectionCurrentRowanLiquidityThreshold(ctx, params.MaxRowanLiquidityThreshold)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateLiquidityProtectionParams,
			sdk.NewAttribute(types.AttributeKeyLiquidityProtectionParams, params.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	return &types.MsgUpdateLiquidityProtectionParamsResponse{}, nil
}

func (k msgServer) ModifyLiquidityProtectionRates(goCtx context.Context, msg *types.MsgModifyLiquidityProtectionRates) (*types.MsgModifyLiquidityProtectionRatesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	response := &types.MsgModifyLiquidityProtectionRatesResponse{}
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}
	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_CLPDEX, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}
	rateParams := k.GetLiquidityProtectionRateParams(ctx)
	rateParams.CurrentRowanLiquidityThreshold = msg.CurrentRowanLiquidityThreshold
	k.SetLiquidityProtectionRateParams(ctx, rateParams)
	events := sdk.EmptyEvents()
	events = events.AppendEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUpdateLiquidityProtectionRateParams,
			sdk.NewAttribute(types.AttributeKeyLiquidityProtectionRateParams, rateParams.String()),
			sdk.NewAttribute(types.AttributeKeyHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer),
		),
	})
	ctx.EventManager().EmitEvents(events)
	return response, nil
}

func (k msgServer) UpdateSwapFeeRate(goCtx context.Context, msg *types.MsgUpdateSwapFeeRateRequest) (*types.MsgUpdateSwapFeeRateResponse, error) {
	response := &types.MsgUpdateSwapFeeRateResponse{}

	// defensive programming
	if msg == nil {
		return response, errors.Errorf("msg was nil")
	}

	if err := msg.ValidateBasic(); err != nil {
		return response, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	signer, err := sdk.AccAddressFromBech32(msg.Signer)
	if err != nil {
		return response, err
	}

	if !k.adminKeeper.IsAdminAccount(ctx, admintypes.AdminType_PMTPREWARDS, signer) {
		return response, errors.Wrap(types.ErrNotEnoughPermissions, fmt.Sprintf("Sending Account : %s", msg.Signer))
	}

	k.SetSwapFeeRate(ctx, &types.SwapFeeRate{SwapFeeRate: msg.SwapFeeRate})

	return response, nil
}
