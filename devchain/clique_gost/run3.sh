#!/bin/bash
echo "Starting node3 ..."

../../build/bin/mir --datadir node3 \
    --identity node3 \
    --syncmode full \
    --port 30313 \
    --bootnodes 'enode://21c4db540114cdb0c8308cb06da448dd233332c25a5a36ac2b2b15f0ba10f9e475e736be92ad064af764770b50332d21e9830a7e3bc24eb15e2b3c009ea69684@127.0.0.1:30310' \
    --networkid 6581 \
    --unlock 0x7219B2cF87ce8760d21cC1ec174402e1bcEd0425 \
    --password node3/password.txt \
    --mine \
    --verbosity 4 \
    --allow-insecure-unlock
