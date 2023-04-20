#!/bin/bash
echo "Starting node3 ..."

../../build/bin/mir --crypto gost \
    --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet \
    --datadir node3 \
    --identity node3 \
    --syncmode full \
    --port 30313 \
    --bootnodes 'enode://fb5f060ea4f9c3caecc9de4f7f9b1b3124373cfbb278f7c064dc68f8a5f31d16b39a3f08d549c8b0eb5399a7dce0503de4a4e83eb92d97c24f22d760e82e9304@127.0.0.1:30310' \
    --networkid 6581 \
    --unlock 0x0eB3D79A3fD6ff8F711b418CcCa19E3ed6c56Ae6 \
    --password node3/password.txt \
    --mine \
    --verbosity 4 \
    --allow-insecure-unlock
