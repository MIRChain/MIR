#!/bin/bash
echo "Starting mainnet Eth sync node ..."

../../build/bin/mir --crypto nist \
  --datadir node1 \
  --identity node1 \
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
  --syncmode full \
  --rinkeby \
  --verbosity 4
