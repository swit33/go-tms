package handlers

import (
	"bufio"
	"fmt"
	"go-tms/pkg/config"
	"go-tms/pkg/fzf"
	"go-tms/pkg/interfaces"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func HandleResult(result fzf.Result, sessions *[]session.Session, cfg *config.Config) error {
	if result.IsAction {
		switch result.Action {
		case fzf.ActionNew:
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return handleActionNew(cwd, sessions)
		case fzf.ActionDelete:
			return handleActionDelete(result, sessions)
		case fzf.ActionInteractive:
			return handleZoxide(sessions, cfg)
		case fzf.ActionSave:
			return session.SaveSessionsToDisk(*sessions)
		}
	} else {
		return handleSessionLogic(false, result.SessionName, sessions)
	}
	return nil
}

func handleSessionLogic(ispath bool, identifier string, sessions *[]session.Session) error {
	runner := interfaces.OsRunner{}

	sessionName, err := tmux.CheckIfSessionExists(ispath, identifier)
	if err != nil {
		return err
	}
	if sessionName != "" {
		return tmux.SwitchSession(sessionName, runner)
	}
	sessionInstance, err := session.GetSessionByName(identifier, *sessions)
	if err == nil {
		return tmux.RestoreSession(sessionInstance, interfaces.OsRunner{})
	}

	return handleActionNew(identifier, sessions)
}

func handleZoxide(sessions *[]session.Session, cfg *config.Config) error {
	result, err := fzf.RunZoxide(cfg)
	if err != nil {
		return err
	}
	return handleSessionLogic(true, result.Arg, sessions)
}

func handleActionNew(path string, sessions *[]session.Session) error {
	runner := interfaces.OsRunner{}

	name, err := FindUniqueSessionName(path, *sessions)
	if err != nil {
		return err
	}
	sessionName, err := tmux.CreateNewSession(name, path, runner)
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
	*sessions, err = session.CombineSessions(tmuxSessions, *sessions)
	if err != nil {
		return err
	}
	return session.SaveSessionsToDisk(*sessions)
}

func handleActionDelete(result fzf.Result, sessions *[]session.Session) error {
	var err error
	sessionName := result.Arg
	if session.CheckIfSessionExists(sessionName, *sessions) {
		*sessions, err = session.DeleteSession(sessionName, *sessions)
		if err != nil {
			return err
		}
	}
	sessionName, err = tmux.CheckIfSessionExists(false, sessionName)
	if err != nil {
		return err
	}
	if sessionName != "" {
		err = tmux.DeleteSession(sessionName)
		if err != nil {
			return err
		}
	}

	return session.SaveSessionsToDisk(*sessions)
}

func FindUniqueSessionName(startPath string, savedSessions []session.Session) (string, error) {
	path := startPath
	var nameParts []string

	for {
		namePart := filepath.Base(path)
		nameParts = append([]string{namePart}, nameParts...)

		sessionName := strings.Join(nameParts, "-")

		re := regexp.MustCompile(`[^a-zA-Z0-9_-]`)
		sanitizedName := re.ReplaceAllString(sessionName, "_")

		tmuxName, err := tmux.CheckIfSessionExists(false, sanitizedName)
		if err != nil {
			return "", err
		}

		sessionExistsOnDisk := session.CheckIfSessionExists(sanitizedName, savedSessions)

		if tmuxName == "" && !sessionExistsOnDisk {
			return sanitizedName, nil
		}

		parentPath := filepath.Dir(path)
		if parentPath == path {
			return "", fmt.Errorf("could not find a unique session name for path: %s", startPath)
		}
		path = parentPath
	}
}

func HandleError(err error) {
	if err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		fmt.Println("Press Enter to continue...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
	}
}
