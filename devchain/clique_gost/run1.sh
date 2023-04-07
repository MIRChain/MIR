#!/bin/bash
echo "Starting node1 ..."

../../build/bin/mir --datadir node1 \
  --identity node1 \
  --syncmode full \
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
  --bootnodes 'enode://161df3a5f868dc64c0f778b2fb3c3724311b9b7f4fdec9eccac3f9a40e065ffb04932ee860d7224fc2972db93ee8a4e43d50e0dca79953ebb0c849d5083f9ab3@127.0.0.1:30310' \
  --networkid 6581 \
  --unlock 0x8ac1983A8E7656A10566c4D795f3509Ee35a41C3 \
  --password node1/password.txt \
  --mine \
  --verbosity 4 \
  --allow-insecure-unlock
