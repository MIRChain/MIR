#!/bin/bash
echo "Starting node2 ..."

../../build/bin/mir --datadir node2 \
    --identity node2 \
    --syncmode full \
    --port 30312 \
    --bootnodes 'enode://21c4db540114cdb0c8308cb06da448dd233332c25a5a36ac2b2b15f0ba10f9e475e736be92ad064af764770b50332d21e9830a7e3bc24eb15e2b3c009ea69684@127.0.0.1:30310' \
    --networkid 6581 \
    --unlock 0x68fcb140304291378dbdc40fa39caeee86548971 \
    --password node2/password.txt \
    --mine \
    --verbosity 4 \
    --allow-insecure-unlock
