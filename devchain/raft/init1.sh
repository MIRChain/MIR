#!/bin/bash
DATADIR=./node1
../../build/bin/mir \
    --raft --datadir $DATADIR init raftGenesis.json