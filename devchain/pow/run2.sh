#!/bin/bash
echo "Starting node2 ..."

../../build/bin/mir --crypto nist \
    --datadir node2 \
    --identity node2 \
    --syncmode full \
    --port 30312 \
    --bootnodes 'enode://fd8b7d623070867bd0458369f5e9f6f4031d105fe559180719846d4a2a82f96d5a5cb987047e86b55b0dafcca786349173f18a3565db9d7ba8c2aecbdfd1ea8d@127.0.0.1:30310' \
    --networkid 6581 \
    --unlock 0x01665a4eb869efbf3af991e0b791d5347718a49d \
    --password node2/password.txt \
    --mine \
    --miner.threads=1 \
    --verbosity 4 \
    --allow-insecure-unlock
