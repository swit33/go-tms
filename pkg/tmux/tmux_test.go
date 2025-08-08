package tmux

import (
	// "fmt"
	"go-tms/pkg/interfaces"
	"go-tms/pkg/session"
	// "strings"
	"testing"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

type RestorationTestCase struct {
	Name        string
	Session     *session.Session
	ExpectedCmd []string
}

var restorationTestCases = []RestorationTestCase{
	{
		Name: "Basic restoration",
		Session: &session.Session{
			Name:        "test01",
			CurrentPath: "/home/aleksej/projects/go-tms",
			Windows: []session.Window{
				{
					Index: "1",
					Panes: []session.Pane{
						{
							Command:     "zsh",
							CurrentPath: "/home/aleksej/projects/go-tms",
							Index:       "1",
						},
					},
				},
				{
					Index: "2",
					Panes: []session.Pane{
						{
							Command:     "zsh",
							CurrentPath: "/home/aleksej/projects/go-tms",
							Index:       "1",
						},
						{
							Command:     "nvim",
							CurrentPath: "/home/aleksej/projects/go-tms",
							Index:       "2",
						},
					},
				},
			},
		},
		ExpectedCmd: []string{
			"tmux new-session -d -s test01 -c /home/aleksej/projects/go-tms",
			"tmux switch-client -t test01",
			"tmux send-keys -t 1.1 cd /home/aleksej/projects/go-tms C-m",
			"tmux send-keys -t 1.1 zsh C-m",
			"tmux new-window -t 2 -c /home/aleksej/projects/go-tms",
			"tmux send-keys -t 2.1 zsh C-m",
			"tmux split-window -t 2.1 -c /home/aleksej/projects/go-tms",
			"tmux send-keys -t 2.2 nvim C-m",
		},
	},
	{
		Name: "Restoration of a session without windows",
		Session: &session.Session{
			Name:        "empty-session",
			CurrentPath: "/home/aleksej/projects",
			Windows:     []session.Window{},
		},
		ExpectedCmd: []string{
			"tmux new-session -d -s empty-session -c /home/aleksej/projects",
			"tmux switch-client -t empty-session",
		},
	},
	{
		Name: "Restoration of a session with multiple panes in a single window",
		Session: &session.Session{
			Name:        "first-window-panes",
			CurrentPath: "/home/aleksej/projects",
			Windows: []session.Window{
				{
					Index: "1",
					Panes: []session.Pane{
						{Command: "zsh", CurrentPath: "/home/aleksej/projects"},
						{Command: "npm start", CurrentPath: "/home/aleksej/projects/go-tms"},
					},
				},
			},
		},
		ExpectedCmd: []string{
			"tmux new-session -d -s first-window-panes -c /home/aleksej/projects",
			"tmux switch-client -t first-window-panes",
			"tmux send-keys -t 1.1 cd /home/aleksej/projects C-m",
			"tmux send-keys -t 1.1 zsh C-m",
			"tmux split-window -t 1.1 -c /home/aleksej/projects/go-tms",
			"tmux send-keys -t 1.2 npm start C-m",
		},
	},
	{
		Name: "Restoration of a session with windows of different paths",
		Session: &session.Session{
			Name:        "multi-window-paths",
			CurrentPath: "/home/aleksej/projects",
			Windows: []session.Window{
				{
					Index: "1",
					Panes: []session.Pane{
						{Command: "zsh", CurrentPath: "/home/aleksej/projects/go-tms"},
					},
				},
				{
					Index: "2",
					Panes: []session.Pane{
						{Command: "docker-compose up", CurrentPath: "/home/aleksej/projects/backend"},
					},
				},
			},
		},
		ExpectedCmd: []string{
			"tmux new-session -d -s multi-window-paths -c /home/aleksej/projects",
			"tmux switch-client -t multi-window-paths",
			"tmux send-keys -t 1.1 cd /home/aleksej/projects/go-tms C-m",
			"tmux send-keys -t 1.1 zsh C-m",
			"tmux new-window -t 2 -c /home/aleksej/projects/backend",
			"tmux send-keys -t 2.1 docker-compose up C-m",
		},
	},
}

func TestRestoreSession(t *testing.T) {
	for _, testCase := range restorationTestCases {
		runner := &interfaces.MockRunner{}
		testSession := testCase.Session
		expectedCommands := testCase.ExpectedCmd
		err := RestoreSession(testSession, runner)
		if err != nil {
			t.Errorf("RestoreSession() in case %s error = %v", testCase.Name, err)
		}
		if len(runner.ExecutedCommands) != len(expectedCommands) {
			t.Errorf("RestoreSession() in case %s expected %d commands, got %d",
				testCase.Name, len(expectedCommands), len(runner.ExecutedCommands))
		}
		for i, cmd := range expectedCommands {
			if cmd != runner.ExecutedCommands[i] {
				t.Errorf("RestoreSession() in case %s expected command %s%d%s\n'%s%s%s' (expected)\n'%s%s%s' (got)",
					testCase.Name,
					colorYellow, i, colorReset,
					colorGreen, cmd, colorReset,
					colorRed, runner.ExecutedCommands[i], colorReset)
			}
		}
	}
}
