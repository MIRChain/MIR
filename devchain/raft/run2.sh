#!/bin/bash
DATADIR=./node2
cp static-nodes.json $DATADIR/mir
../../build/bin/mir --raft --nodiscover --emitcheckpoints --raftport 50001 --port 21001 \
    --datadir $DATADIR  \
    --ws --ws.port 6001 --ws.addr 0.0.0.0  --ws.api admin,debug,eth,miner,net,personal,rpc,txpool,web3,raft --ws.origins=*  \
    --http --http.corsdomain=* --http.vhosts=* --http.api admin,debug,eth,miner,net,personal,rpc,txpool,web3,raft --http.port 20001 --http.addr 0.0.0.0 \
    --verbosity 10