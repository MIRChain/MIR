#!/bin/bash
set -x

sudo chmod +x run1.sh
sudo chmod +x run2.sh
sudo chmod +x run3.sh

tmux new -s "node1" -d &&
tmux send-keys -t "node1" "./run1.sh" C-m &&
tmux detach -s "node1"

tmux new -s "node2" -d &&
tmux send-keys -t "node2" "./run2.sh" C-m &&
tmux detach -s "node2"

tmux new -s "node3" -d &&
tmux send-keys -t "node3" "./run3.sh" C-m &&
tmux detach -s "node3"
