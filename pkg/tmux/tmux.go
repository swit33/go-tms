package tmux

import (
	"fmt"
	"github.com/swit33/go-tms/pkg/config"
	"github.com/swit33/go-tms/pkg/interfaces"
	"github.com/swit33/go-tms/pkg/session"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

func ListSessions() ([]session.Session, error) {
	var output []byte
	var err error

	cmd := exec.Command("tmux", "list-panes", "-a", "-F", "#{session_name}|#{session_path}|#{window_index}|#{pane_index}|#{pane_current_command}|#{pane_current_path}")
	output, err = cmd.Output()
	if err != nil {
		if strings.Contains(err.Error(), "no server running") {
			return []session.Session{}, nil
		}
		return nil, fmt.Errorf("failed to list sessions: %v", err)
	}

	sessionsMap := make(map[string]*session.Session)
	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")

	for line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 6 {
			continue
		}

		sessionName := parts[0]
		sessionPath := parts[1]
		windowIndex := parts[2]
		paneIndex := parts[3]
		paneCommand := parts[4]
		panePath := parts[5]

		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home dir: %v", err)
		}

		if sessionPath == userHomeDir {
			continue
		}
		if sessionPath == "/dev/null" {
			continue
		}
		if sessionPath == "/tmp" {
			continue
		}

		if sessionName == "go-tms-startup" {
			continue
		}

		sessionInst, ok := sessionsMap[sessionName]
		if !ok {
			sessionInst = &session.Session{
				Name:        sessionName,
				CurrentPath: sessionPath,
			}
			sessionsMap[sessionName] = sessionInst
		}

		var windowInst *session.Window
		for i := range sessionInst.Windows {
			if sessionInst.Windows[i].Index == windowIndex {
				windowInst = &sessionInst.Windows[i]
				break
			}
		}
		if windowInst == nil {
			newWindow := session.Window{
				Index: windowIndex,
				Panes: make([]session.Pane, 0),
			}
			sessionInst.Windows = append(sessionInst.Windows, newWindow)
			windowInst = &sessionInst.Windows[len(sessionInst.Windows)-1]
		}

		paneInst := session.Pane{
			Command:     paneCommand,
			CurrentPath: panePath,
			Index:       paneIndex,
		}
		windowInst.Panes = append(windowInst.Panes, paneInst)
	}
	sessions := make([]session.Session, 0, len(sessionsMap))
	for _, s := range sessionsMap {
		sessions = append(sessions, *s)
	}
	return sessions, nil
}

func CreateNewSession(sessionName string, directory string, runner interfaces.Runner) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", directory)
	if err := runner.Run(cmd); err != nil {
		return "", fmt.Errorf("failed to create new session: %v", err)
	}
	return sessionName, nil
}

func SwitchSession(sessionName string, runner interfaces.Runner) error {
	cmd := exec.Command("tmux", "switch-client", "-t", sessionName)
	if err := runner.Run(cmd); err != nil {
		return fmt.Errorf("failed to switch to session: %v sessionname: %s", err, sessionName)
	}
	return nil
}

func CreateBootSession(sessionName string, executable string, runner interfaces.Runner) (string, error) {
	var cmd *exec.Cmd
	cmd = exec.Command("tmux", "new-session", "-d", "-s", sessionName, executable)
	if err := runner.Run(cmd); err != nil {
		return "", fmt.Errorf("failed to create new session: %v", err)
	}
	return sessionName, nil
}

func AttachSession(sessionName string, runner interfaces.Runner) error {
	var cmd *exec.Cmd
	if sessionName == "" {
		cmd = exec.Command("tmux", "attach-session")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd = exec.Command("tmux", "attach-session", "-t", sessionName)
	}
	fmt.Println(cmd)
	if err := runner.Run(cmd); err != nil {
		return fmt.Errorf("failed to attach to session: %v sessionname: %s", err, sessionName)
	}
	return nil
}

func RestoreSession(s *session.Session, runner interfaces.Runner, cfg *config.Config) error {
	if _, err := CreateNewSession(s.Name, s.CurrentPath, runner); err != nil {
		return err
	}
	err := SwitchSession(s.Name, runner)
	if err != nil {
		return err
	}

	for i, window := range s.Windows {
		if i != 0 {
			cmd := exec.Command("tmux", "new-window",
				"-t", s.Name+":"+strconv.Itoa(i+1), "-c", window.Panes[0].CurrentPath)
			if err := runner.Run(cmd); err != nil {
				return fmt.Errorf("failed to create new window: %v", err)
			}
		}
		for j, pane := range window.Panes {
			if j == 0 {
				if i == 0 {
					cmd := exec.Command("tmux", "send-keys",
						"-t", s.Name+":"+strconv.Itoa(i+1)+"."+strconv.Itoa(j+1), "cd "+pane.CurrentPath, "C-m")
					if err := runner.Run(cmd); err != nil {
						return fmt.Errorf("failed to set pane path: %v", err)
					}
				}
			} else {
				cmd := exec.Command("tmux", "split-window",
					"-t", s.Name+":"+strconv.Itoa(i+1)+"."+strconv.Itoa(j), "-c", pane.CurrentPath)
				if err := runner.Run(cmd); err != nil {
					return fmt.Errorf("failed to split window: %v", err)
				}

			}
			if slices.Contains(strings.Split(cfg.ProgramWhitelist, ","), pane.Command) {
				var cmd *exec.Cmd
				if pane.Command == "nvim" && cfg.NvimCustomCommand != "" {
					cmd = exec.Command("tmux", "send-keys",
						"-t", s.Name+":"+strconv.Itoa(i+1)+"."+strconv.Itoa(j+1), cfg.NvimCustomCommand, "C-m")
				} else {
					cmd = exec.Command("tmux", "send-keys",
						"-t", s.Name+":"+strconv.Itoa(i+1)+"."+strconv.Itoa(j+1), pane.Command, "C-m")
				}
				if err := runner.Run(cmd); err != nil {
					return fmt.Errorf("failed to run pane command: %v", err)
				}
			}
		}
		if cfg.SelectFirst {
			cmd := exec.Command("tmux", "select-window", "-t", "1")
			if err := runner.Run(cmd); err != nil {
				return fmt.Errorf("failed to select window: %v", err)
			}
		}
	}
	return nil
}

func CheckIfSessionExists(ispath bool, identifier string) (string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}|#{session_path}")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
	for line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 2 {
			continue
		}

		sessionName := parts[0]
		sessionPath := parts[1]
		if ispath {
			if sessionPath == identifier {
				return sessionName, nil
			}
		} else {
			if sessionName == identifier {
				return sessionName, nil
			}
		}
	}
	return "", nil
}

func DeleteSession(sessionName string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete session: %v", err)
	}
	return nil

}
