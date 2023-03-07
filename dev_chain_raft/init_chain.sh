#!/bin/bash
rm -rf ./node1/masterchain/
rm -rf ./node1/quorum-raft-state/
rm -rf ./node1/raft-snap/
rm -rf ./node1/raft-wal/
rm ./node1/geth.ipc
echo "Cleaned node1"
rm -rf ./node2/masterchain/
rm -rf ./node2/quorum-raft-state/
rm -rf ./node2/raft-snap/
rm -rf ./node2/raft-wal/
rm ./node2/geth.ipc
echo "Cleaned node2"
rm -rf ./node3/masterchain/
rm -rf ./node3/quorum-raft-state/
rm -rf ./node3/raft-snap/
rm -rf ./node3/raft-wal/
rm ./node3/geth.ipc
echo "Cleaned node3"


./init1.sh
./init2.sh
./init3.sh