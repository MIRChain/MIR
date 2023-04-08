#!/bin/bash
rm -rf ./node1/mir/
rm -rf ./node1/quorum-raft-state/
rm -rf ./node1/raft-snap/
rm -rf ./node1/raft-wal/
rm ./node1/mir.ipc
echo "Cleaned node1"
rm -rf ./node2/mir/
rm -rf ./node2/quorum-raft-state/
rm -rf ./node2/raft-snap/
rm -rf ./node2/raft-wal/
rm ./node2/mir.ipc
echo "Cleaned node2"
rm -rf ./node3/mir/
rm -rf ./node3/quorum-raft-state/
rm -rf ./node3/raft-snap/
rm -rf ./node3/raft-wal/
rm ./node3/mir.ipc
echo "Cleaned node3"