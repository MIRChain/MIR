#!/bin/bash
set -x

tmux kill-ses -t "node1"
tmux kill-ses -t "node2"
tmux kill-ses -t "node3"