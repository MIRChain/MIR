#!/bin/bash
DATADIR=./node1
./geth \
    --raft --datadir $DATADIR init raftGenesis.json