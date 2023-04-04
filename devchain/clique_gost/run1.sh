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
  --bootnodes 'enode://21c4db540114cdb0c8308cb06da448dd233332c25a5a36ac2b2b15f0ba10f9e475e736be92ad064af764770b50332d21e9830a7e3bc24eb15e2b3c009ea69684@127.0.0.1:30310' \
  --networkid 6581 \
  --unlock 0x44E22398e9A9C39c7a1fD10cB067d38a0D0b9167 \
  --password node1/password.txt \
  --mine \
  --verbosity 4 \
  --allow-insecure-unlock
