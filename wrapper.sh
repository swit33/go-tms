#!/bin/bash

if [[ -z $(tmux list-sessions -F '#S') ]]; then

	tmux new-session -d -s "go-tms-startup" "/Users/aleksejbagno/projects/go-tms/release/darwin/go-tms"

	tmux attach-session -t "go-tms-startup"
else

	tmux attach-session
fi
