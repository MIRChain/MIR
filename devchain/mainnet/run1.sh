#!/bin/bash
echo "Starting mainnet miner ..."

../../build/bin/mir --crypto gost \
  --datadir node1 \
  --identity mainnet-worker \
  --port 30311  \
  --ws \
  --ws.addr 0.0.0.0 \
  --ws.port 8546 \
  --ws.origins "*" \
  --http \
  --http.addr 0.0.0.0 \
  --http.port 8545 \
  --http.corsdomain "*" \
  --http.api shh,personal,db,eth,net,web3,txpool,miner,admin \
  --mine \
  --miner.threads=1 \
  --mainnet \
  --ethstats PK:buran@194.87.253.126:3000 \
  --verbosity 4
