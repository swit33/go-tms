package interfaces

import "os/exec"

type Runner interface {
	Run(cmd *exec.Cmd) error
}

type OsRunner struct{}

type MockRunner struct {
	ExecutedCommands []string
}

func (r OsRunner) Run(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (r *MockRunner) Run(cmd *exec.Cmd) error {
	var s string
	argslen := len(cmd.Args)
	for i, arg := range cmd.Args {
		if i == argslen-1 {
			s += arg
		} else {
			s += arg + " "
		}
	}
	r.ExecutedCommands = append(r.ExecutedCommands, s)
	return nil
}
