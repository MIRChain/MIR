#!/bin/bash
DATADIR=./node3
../../build/bin/mir \
    --raft --datadir $DATADIR init raftGenesis.json