#!/bin/bash
echo "Starting node3 ..."

../../build/bin/mir --datadir node3 \
    --identity node3 \
    --syncmode full \
    --port 30313 \
    --bootnodes 'enode://71f6dbb06e962829f119ea72835a997b790de3d2ceb299394f10c343bb5ea24a836aa02419adee67a026347ab0274c72ea77143a233f0692c87a98965776c1b9@127.0.0.1:30311' \
    --networkid 6581 \
    --unlock 0x3833067356d624e36fa8cfaf208e97263f3e0703 \
    --password node3/password.txt \
    --mine \
    --verbosity 4 \
    --allow-insecure-unlock
