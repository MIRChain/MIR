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
  --bootnodes 'enode://b68b4762248b8cb21e83fb16c14a8f2e360bc20d2b811d5b15a899d9fb9d59b698a47a0d335883db9318cf99f19b60f4426d9e4633d8d9129b59b9ece3e1b599@127.0.0.1:30312' \
  --networkid 6581 \
  --unlock 0xb47f736b9b15dcc888ab790c38a6ad930217cbee \
  --password node1/password.txt \
  --mine \
  --verbosity 4 \
  --allow-insecure-unlock
