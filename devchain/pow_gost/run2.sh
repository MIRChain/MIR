#!/bin/bash
echo "Starting node2 ..."

../../build/bin/mir --crypto gost \
    --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet \
    --datadir node2 \
    --identity node2 \
    --syncmode full \
    --port 30312 \
    --bootnodes 'enode://161df3a5f868dc64c0f778b2fb3c3724311b9b7f4fdec9eccac3f9a40e065ffb04932ee860d7224fc2972db93ee8a4e43d50e0dca79953ebb0c849d5083f9ab3@127.0.0.1:30310' \
    --networkid 6581 \
    --unlock 0x6B48C41ACB6B9A9ccDA279E968c2499b3e003793 \
    --password node2/password.txt \
    --mine \
    --miner.threads=1 \
    --verbosity 4 \
    --allow-insecure-unlock
