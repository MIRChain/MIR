#!/bin/bash
DATADIR=./node3
cp static-nodes.json $DATADIR/geth
./geth --raft --nodiscover --emitcheckpoints --raftport 50002 --port 21002 \
    --datadir $DATADIR  \
    --ws --ws.port 6002 --ws.addr 0.0.0.0  --ws.api admin,debug,eth,miner,net,personal,rpc,txpool,web3,raft --ws.origins=*  \
    --http --http.corsdomain=* --http.vhosts=* --http.api admin,debug,eth,miner,net,personal,rpc,txpool,web3,raft --http.port 20002 --http.addr 0.0.0.0 \
    --verbosity 10