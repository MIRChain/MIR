#!/bin/bash
echo "Starting node2 ..."

../../build/bin/mir --crypto gost \
    --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet \
    --datadir node2 \
    --identity soyuz-worker2 \
    --syncmode full \
    --gcmode=archive \
    --port 30312 \
    --soyuz \
    --verbosity 4
