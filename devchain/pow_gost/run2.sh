#!/bin/bash
echo "Starting node2 ..."

../../build/bin/mir --crypto gost \
    --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet \
    --datadir node2 \
    --identity node2 \
    --syncmode full \
    --port 30312 \
    --bootnodes 'enode://fb5f060ea4f9c3caecc9de4f7f9b1b3124373cfbb278f7c064dc68f8a5f31d16b39a3f08d549c8b0eb5399a7dce0503de4a4e83eb92d97c24f22d760e82e9304@127.0.0.1:30310' \
    --networkid 6581 \
    --verbosity 4
