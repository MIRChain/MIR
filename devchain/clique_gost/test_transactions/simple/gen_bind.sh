#!/bin/bash

# ../../../../build/bin/abigen --crypto gost --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet --pkg simple --sol simple.sol --out simple.go

./solc  --bin --abi ./simple.sol  --allow-paths . -o . --overwrite && \
../../../../build/bin/abigen --crypto gost --gostcurve id-GostR3410-2001-CryptoPro-A-ParamSet --pkg simple --bin=Simple.bin --abi=Simple.abi --out simple.go