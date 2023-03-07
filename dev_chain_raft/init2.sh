#!/bin/bash
DATADIR=./node2
./geth \
    --raft --datadir $DATADIR init raftGenesis.json