package main

import (
	"flag"
	"go-tms/pkg/config"
	"go-tms/pkg/daemon"
	"go-tms/pkg/fzf"
	"go-tms/pkg/handlers"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
)

func main() {

	daemonMode := flag.Bool("d", false, "Run in daemon mode with autosave enabled")
	flag.Parse()

	if *daemonMode {
		daemon.RunDaemon()
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		handlers.HandleError(err)
	}
	sessions, err := session.LoadSessionsFromDisk()
	if err != nil {
		handlers.HandleError(err)
	}
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		handlers.HandleError(err)
	}
	combinedSessions, err := session.CombineSessions(tmuxSessions, sessions)
	if err != nil {
		handlers.HandleError(err)
	}
	result, err := fzf.RunSessions(combinedSessions, &cfg)
	if err != nil {
		handlers.HandleError(err)
	}
	err = handlers.HandleResult(result, &combinedSessions, &cfg)
	if err != nil {
		handlers.HandleError(err)
	}
}
