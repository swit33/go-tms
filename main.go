package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/swit33/go-tms/pkg/boot"
	"github.com/swit33/go-tms/pkg/config"
	"github.com/swit33/go-tms/pkg/daemon"
	"github.com/swit33/go-tms/pkg/fzf"
	"github.com/swit33/go-tms/pkg/interfaces"
	"github.com/swit33/go-tms/pkg/session"
	"github.com/swit33/go-tms/pkg/tmux"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	daemonMode := flag.Bool("d", false, "Run in daemon mode with autosave enabled")
	bootMode := flag.Bool("b", false, "Run in boot mode")
	switcherMode := flag.Bool("s", false, "Run in switcher mode")

	flag.Parse()

	cfg, err := config.LoadConfig()
	if err != nil {
		handleError(err)
	}

	if *daemonMode {
		daemon.RunDaemon(&cfg)
		return
	}

	if *bootMode {
		err := boot.RunBoot(&cfg)
		if err != nil {
			handleError(err)
		}
		return
	}

	if *switcherMode {
		err = runSwitcher(&cfg)
		if err != nil {
			handleError(err)
		}
	}
}

func runSwitcher(cfg *config.Config) error {
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
	err = handleResult(result, &combinedSessions, cfg)
	if err != nil {
		return err
	}
	return nil
}

func handleResult(result fzf.Result, sessions *[]session.Session, cfg *config.Config) error {
	if result.IsAction {
		switch result.Action {
		case fzf.ActionNew:
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			return handleActionNew(cwd, sessions)
		case fzf.ActionDelete:
			return handleActionDelete(result, sessions, cfg)
		case fzf.ActionInteractive:
			return handleZoxide(sessions, cfg)
		case fzf.ActionSave:
			return handleSave(sessions, cfg)
		}
	} else {
		return handleSessionLogic(false, result.SessionName, sessions, cfg)
	}
	return nil
}

func handleSave(sessions *[]session.Session, cfg *config.Config) error {
	err := session.SaveSessionsToDisk(*sessions)
	if err != nil {
		return err
	}
	return runSwitcher(cfg)
}

func handleSessionLogic(ispath bool, identifier string, sessions *[]session.Session, cfg *config.Config) error {
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
		return tmux.RestoreSession(sessionInstance, interfaces.OsRunner{}, cfg)
	}

	return handleActionNew(identifier, sessions)
}

func handleZoxide(sessions *[]session.Session, cfg *config.Config) error {
	result, err := fzf.RunZoxide(cfg)
	if err != nil {
		return err
	}
	if result.IsAction && result.Action == fzf.ActionReturn {
		return runSwitcher(cfg)
	}
	return handleSessionLogic(true, result.Arg, sessions, cfg)
}

func handleActionNew(path string, sessions *[]session.Session) error {
	runner := interfaces.OsRunner{}

	name, err := findUniqueSessionName(path, *sessions)
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

func handleActionDelete(result fzf.Result, sessions *[]session.Session, cfg *config.Config) error {
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
	err = session.SaveSessionsToDisk(*sessions)
	if err != nil {
		return err
	}

	return runSwitcher(cfg)
}

func findUniqueSessionName(startPath string, savedSessions []session.Session) (string, error) {
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

func handleError(err error) {
	if err != nil {
		fmt.Printf("\033[31mError: %v\033[0m\n", err)
		fmt.Println("Press Enter to continue...")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
	}
}
