#!/bin/bash
echo "Starting mainnet miner ..."

../../build/bin/mir --crypto gost \
  --datadir node2 \
  --identity mainnet-worker \
  --port 30312  \
  --mine \
  --miner.threads=1 \
  --mainnet.mir \
  --ethstats PK2:buran@194.87.253.126:3000 \
  --verbosity 4
