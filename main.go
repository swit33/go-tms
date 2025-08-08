package main

import (
	"go-tms/pkg/fzf"
	"go-tms/pkg/handlers"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
)

func main() {

	sessions, err := session.LoadSessionsFromDisk()
	if err != nil {
		handlers.HandleError(err)
	}
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		handlers.HandleError(err)
	}
	combinedSessions, err := session.CombineSessions(sessions, tmuxSessions)
	if err != nil {
		handlers.HandleError(err)
	}
	result, err := fzf.RunSessions(combinedSessions)
	if err != nil {
		handlers.HandleError(err)
	}
	err = handlers.HandleResult(result, &combinedSessions)
	if err != nil {
		handlers.HandleError(err)
	}
}
