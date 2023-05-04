tmux new -s "node1" -d &&
tmux send-keys -t "node1" "sh run1.sh" C-m &&
tmux detach -s "node1"


