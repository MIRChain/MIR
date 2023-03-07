#!/bin/bash
DATADIR=./node3
./geth \
    --raft --datadir $DATADIR init raftGenesis.json