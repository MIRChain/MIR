#!/bin/bash
echo "Starting node1 ..."

../../build/bin/mir --crypto gost \
  --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet \
  --datadir node1 \
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
  --bootnodes 'enode://fb5f060ea4f9c3caecc9de4f7f9b1b3124373cfbb278f7c064dc68f8a5f31d16b39a3f08d549c8b0eb5399a7dce0503de4a4e83eb92d97c24f22d760e82e9304@127.0.0.1:30310' \
  --networkid 6581 \
  --unlock 0x1F1a2F8231eFe45f0fF0c2ed3c27eeBF58Fc175c \
  --password node1/password.txt \
  --mine \
  --verbosity 4 \
  --allow-insecure-unlock
