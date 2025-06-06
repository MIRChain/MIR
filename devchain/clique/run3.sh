#!/bin/bash
echo "Starting node3 ..."

../../build/bin/mir --crypto nist \
    --datadir node3 \
    --identity node3 \
    --syncmode full \
    --port 30313 \
    --bootnodes 'enode://fd8b7d623070867bd0458369f5e9f6f4031d105fe559180719846d4a2a82f96d5a5cb987047e86b55b0dafcca786349173f18a3565db9d7ba8c2aecbdfd1ea8d@127.0.0.1:30310' \
    --networkid 6581 \
    --unlock 0x3833067356d624e36fa8cfaf208e97263f3e0703 \
    --password node3/password.txt \
    --mine \
    --verbosity 4 \
    --allow-insecure-unlock
