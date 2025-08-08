package fzf

import (
	"fmt"
	"go-tms/pkg/session"
	"os"
	"os/exec"
	"strings"
)

var fzfOpts = []string{
	"--no-sort",
	"--reverse",
}

var fzfBinds = []string{
	"--bind",
	"ctrl-n:become(echo 'gotms_act_new:{}')",
	"--bind",
	"ctrl-d:become(echo 'gotms_act_delete:{}')",
	"--bind",
	"ctrl-i:become(echo 'gotms_act_interactive:{}')",
	"--bind",
	"ctrl-s:become(echo 'gotms_act_save:{}')",
}

var fzfHeader = []string{
	"--header",
	"<C-n>: new session\n<C-d>: delete session\n<C-i>: interactive search\n<C-s>: save session",
}

var fzfPrompt string = "Sessions> "

type Action string

const (
	ActionNew         Action = "gotms_act_new"
	ActionDelete      Action = "gotms_act_delete"
	ActionInteractive Action = "gotms_act_interactive"
	ActionSave        Action = "gotms_act_save"
)

type Result struct {
	Action      Action
	IsAction    bool
	Arg         string
	SessionName string
}

func Run(entries []string) (string, error) {
	args := []string{}
	args = append(args, fzfOpts...)
	args = append(args, fzfBinds...)
	args = append(args, "--prompt", fzfPrompt)
	args = append(args, fzfHeader...)
	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(strings.Join(entries, "\n"))

	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil // Return no output and no error for cancellation
		}
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return strings.TrimSpace(string(output)), nil
		}
		return "", fmt.Errorf("fzf command failed: %v", err)
	}

	if len(output) == 0 {
		return "", nil
	}

	return strings.TrimSpace(string(output)), nil
}

func RunSessions(s []session.Session) (Result, error) {
	entries := make([]string, 0)
	for _, s := range s {
		entries = append(entries, s.Name)
	}
	result, err := Run(entries)
	if err != nil {
		return Result{}, err
	}

	if result == "" {
		os.Exit(0)
	}

	if strings.HasPrefix(result, "gotms_act_") {
		parts := strings.Split(result, ":")
		return Result{IsAction: true, Action: Action(parts[0]), Arg: parts[1]}, nil
	}
	sessionName := strings.TrimSpace(string(result))
	return Result{IsAction: false, SessionName: sessionName}, nil
}

func RunZoxide() (Result, error) {
	cmd := exec.Command("zoxide", "query", "-i")
	output, err := cmd.Output()
	if err != nil {
		return Result{}, fmt.Errorf("zoxide command failed: %v", err)
	}
	path := strings.TrimSpace(string(output))
	return Result{IsAction: true, Action: ActionInteractive, Arg: path}, nil
}
