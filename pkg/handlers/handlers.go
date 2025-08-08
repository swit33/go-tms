package handlers

import (
	"bufio"
	"fmt"
	"go-tms/pkg/fzf"
	"go-tms/pkg/interfaces"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
	"os"
)

func HandleResult(result fzf.Result, sessions *[]session.Session) error {
	if result.IsAction {
		switch result.Action {
		case fzf.ActionNew:
			return handleActionNew(sessions)
		case fzf.ActionDelete:
			return handleActionDelete(result, sessions)
		case fzf.ActionInteractive:
			return handleZoxide(sessions)
		case fzf.ActionSave:
			return session.SaveSessionsToDisk(*sessions)
		}
	} else {
		return handleSessionLogic(result.SessionName, sessions)
	}
	return nil
}

func handleSessionLogic(identifier string, sessions *[]session.Session) error {
	runner := interfaces.OsRunner{}
	if tmux.CheckIfSessionExists(identifier) {
		return tmux.SwitchSession(identifier, runner)
	}
	sessionInstance, err := session.GetSessionByName(identifier, *sessions)
	if err == nil {
		if err := tmux.RestoreSession(sessionInstance, interfaces.OsRunner{}); err != nil {
			return err
		}
		tmuxSessions, err := tmux.ListSessions()
		if err != nil {
			return err
		}
		*sessions, err = session.CombineSessions(*sessions, tmuxSessions)
		if err != nil {
			return err
		}
		return session.SaveSessionsToDisk(*sessions)
	}
	sessionName, err := tmux.CreateNewSession("", identifier, runner)
	if err != nil {
		return err
	}
	if err := tmux.SwitchSession(sessionName, runner); err != nil {
		return err
	}
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}
	*sessions, err = session.CombineSessions(*sessions, tmuxSessions)
	if err != nil {
		return err
	}
	return session.SaveSessionsToDisk(*sessions)
}

func handleZoxide(sessions *[]session.Session) error {
	result, err := fzf.RunZoxide()
	if err != nil {
		return err
	}
	return handleSessionLogic(result.Arg, sessions)
}

func handleActionNew(sessions *[]session.Session) error {
	runner := interfaces.OsRunner{}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	name, err := tmux.CreateNewSession("", cwd, runner)
	if err != nil {
		return err
	}
	if err := tmux.SwitchSession(name, runner); err != nil {
		return err
	}
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		return err
	}
	*sessions, err = session.CombineSessions(*sessions, tmuxSessions)
	if err != nil {
		return err
	}
	return session.SaveSessionsToDisk(*sessions)
}

func handleActionDelete(result fzf.Result, sessions *[]session.Session) error {
	var err error
	if session.CheckIfSessionExists(result.Arg, *sessions) {
		*sessions, err = session.DeleteSession(result.Arg, *sessions)
		if err != nil {
			return err
		}
	}
	if tmux.CheckIfSessionExists(result.Arg) {
		err = tmux.DeleteSession(result.Arg)
		if err != nil {
			return err
		}
	}
	return session.SaveSessionsToDisk(*sessions)
}

func HandleError(err error) {
	if err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		fmt.Println("Press Enter to continue...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
	}
}
