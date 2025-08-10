package boot

import (
	"fmt"
	"go-tms/pkg/config"
	"go-tms/pkg/daemon"
	"go-tms/pkg/interfaces"
	"go-tms/pkg/tmux"
	"os"
	"os/exec"
	"strings"
)

func RunBoot(cfg *config.Config) error {
	cmd := exec.Command("tmux", "list-sessions")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() > 1 {
			return fmt.Errorf("failed to list sessions: %v", err)
		}
	}
	sessions := strings.TrimSpace(string(output))
	self := os.Args[0] + " -s"

	if sessions == "" {
		daemon.StartDaemon(cfg)

		_, err := tmux.CreateBootSession("go-tms-startup", self, interfaces.OsRunner{})
		if err != nil {
			return err
		}

		err = tmux.AttachSession("", interfaces.OsRunner{})
		if err != nil {
			return err
		}
		return nil
	} else {
		err = tmux.AttachSession("", interfaces.OsRunner{})
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}
