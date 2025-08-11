package fzf

import (
	"fmt"
	"go-tms/pkg/config"
	"go-tms/pkg/session"
	"os"
	"os/exec"
	"strings"
)

var fzfHeader = []string{
	"--header",
	"<C-n>: new session\n<C-d>: delete session\n<C-i>: interactive search\n<C-s>: save session",
}

var fzfPrompt string = "Sessions> "

const ActionPrefix = "gotms_act_"

type Action string

const (
	ActionNew         Action = ActionPrefix + "new"
	ActionDelete      Action = ActionPrefix + "delete"
	ActionInteractive Action = ActionPrefix + "interactive"
	ActionSave        Action = ActionPrefix + "save"
	ActionReturn      Action = ActionPrefix + "return"
)

type Result struct {
	Action      Action
	IsAction    bool
	Arg         string
	SessionName string
}

func Run(entries []string, cfg *config.Config) (string, error) {
	var binds []string
	binds = append(binds,
		fmt.Sprintf("--bind=%s:become(echo '%s:{}')", cfg.FZFBindNew, ActionNew))
	binds = append(binds,
		fmt.Sprintf("--bind=%s:become(echo '%s:{}')", cfg.FZFBindDelete, ActionDelete))
	binds = append(binds,
		fmt.Sprintf("--bind=%s:become(echo '%s:{}')", cfg.FZFBindInteractive, ActionInteractive))
	binds = append(binds,
		fmt.Sprintf("--bind=%s:become(echo '%s:{}')", cfg.FZFBindSave, ActionSave))
	var header []string
	header = append(header, "--header")
	header = append(header,
		fmt.Sprintf("<%s>: new session\n<%s>: delete session\n<%s>: interactive search\n<%s>: save session",
			cfg.FZFBindNew, cfg.FZFBindDelete, cfg.FZFBindInteractive, cfg.FZFBindSave))
	args := []string{}
	args = append(args, strings.Fields(cfg.FZFOpts)...)
	args = append(args, binds...)
	args = append(args, "--prompt", cfg.FZFPrompt)
	args = append(args, header...)
	cmd := exec.Command("fzf", args...)
	cmd.Stdin = strings.NewReader(strings.Join(entries, "\n"))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return "", nil
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

func RunSessions(s []session.Session, cfg *config.Config) (Result, error) {
	entries := make([]string, 0)
	for _, s := range s {
		entries = append(entries, s.Name)
	}
	result, err := Run(entries, cfg)
	if err != nil {
		return Result{}, err
	}

	if result == "" {
		os.Exit(0)
	}

	if strings.HasPrefix(result, ActionPrefix) {
		parts := strings.Split(result, ":")
		return Result{IsAction: true, Action: Action(parts[0]), Arg: parts[1]}, nil
	}
	sessionName := strings.TrimSpace(string(result))
	return Result{IsAction: false, SessionName: sessionName}, nil
}

func RunZoxide(cfg *config.Config) (Result, error) {
	cmd := exec.Command("zoxide", "query", "-i")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "_ZO_FZF_OPTS="+cfg.ZoxideOpts)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return Result{IsAction: true, Action: ActionReturn}, nil
		}
		return Result{}, fmt.Errorf("zoxide command failed: %v", err)
	}
	path := strings.TrimSpace(string(output))
	return Result{IsAction: true, Action: ActionInteractive, Arg: path}, nil
}
