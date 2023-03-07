#!/bin/bash
DATADIR=./node1
cp static-nodes.json $DATADIR/geth
./geth --raft --nodiscover --raftport 50000 --port 21000 --emitcheckpoints \
    --datadir $DATADIR  \
    --pprof --pprof.addr 0.0.0.0 --pprof.port 6060 --metrics \
    --ws --ws.port 6000 --ws.addr 0.0.0.0  --ws.api admin,debug,eth,miner,net,personal,rpc,txpool,web3,raft --ws.origins=*  \
    --http --http.corsdomain=* --http.vhosts=* --http.api admin,debug,eth,miner,net,personal,rpc,txpool,web3,raft --http.port 8545 --http.addr 0.0.0.0 \
    --unlock 0x25FAcE3782641d4ad4C7eaf2Ee5eb8a7DCfaB465 --allow-insecure-unlock --password ./$DATADIR/keystore/accountPassword \
    --verbosity 10