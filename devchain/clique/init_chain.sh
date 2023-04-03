#!/bin/bash
rm -rf ./node1/mir/
echo "Cleaned node1"
rm -rf ./node2/mir/
echo "Cleaned node2"
rm -rf ./node3/mir/
echo "Cleaned node3"

../../build/bin/mir --datadir node1/ init genesis.json
../../build/bin/mir --datadir node2/ init genesis.json
../../build/bin/mir --datadir node3/ init genesis.json
