#!/usr/bin/env bash

set -x

sifnoded tx clp remove-liquidity-units \
  --from $SIF_ACT \
  --keyring-backend test \
  --symbol cusdc \
  --withdrawUnits 1559429888878798180197958710 \
  --fees 100000000000000000rowan \
  --chain-id $SIFNODE_CHAIN_ID \
  --node ${SIFNODE_NODE} \
  --broadcast-mode block \
  -y
