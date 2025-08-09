package switcher

import (
	"go-tms/pkg/config"
	"go-tms/pkg/fzf"
	"go-tms/pkg/handlers"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
)

func RunSwitcher(cfg *config.Config) error {
	sessions, err := session.LoadSessionsFromDisk()
	if err != nil {
		return err
	}
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}
	combinedSessions, err := session.CombineSessions(tmuxSessions, sessions)
	if err != nil {
		return err
	}
	result, err := fzf.RunSessions(combinedSessions, cfg)
	if err != nil {
		return err
	}
	err = handlers.HandleResult(result, &combinedSessions, cfg)
	if err != nil {
		return err
	}
	return nil
}
