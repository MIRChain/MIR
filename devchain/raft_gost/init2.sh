#!/bin/bash
DATADIR=./node2
../../build/bin/mir \
    --raft --datadir $DATADIR init raftGenesis.json