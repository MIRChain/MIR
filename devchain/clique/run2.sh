#!/bin/bash
echo "Starting node2 ..."

../../build/bin/mir --datadir node2 \
    --identity node2 \
    --syncmode full \
    --port 30312 \
    --bootnodes 'enode://71f6dbb06e962829f119ea72835a997b790de3d2ceb299394f10c343bb5ea24a836aa02419adee67a026347ab0274c72ea77143a233f0692c87a98965776c1b9@127.0.0.1:30311' \
    --networkid 6581 \
    --unlock 0x01665a4eb869efbf3af991e0b791d5347718a49d \
    --password node2/password.txt \
    --mine \
    --verbosity 4 \
    --allow-insecure-unlock
