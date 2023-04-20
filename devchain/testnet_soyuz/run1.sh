#!/bin/bash
echo "Starting testnet miner ..."

../../build/bin/mir --crypto gost \
  --datadir node1 \
  --identity soyuz-worker1 \
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
  --miner.threads=4 \
  --soyuz \
  --ethstats node1:soyuz@194.87.80.101:3000 \
  --verbosity 4
