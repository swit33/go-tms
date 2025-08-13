package daemon

import (
	"fmt"
	"github.com/swit33/go-tms/pkg/config"
	"github.com/swit33/go-tms/pkg/session"
	"github.com/swit33/go-tms/pkg/tmux"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func StartDaemon(cfg *config.Config) {
	self, err := os.Executable()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	cmd := exec.Command(self, "-d")
	cmd.Start()
	return
}

func RunDaemon(cfg *config.Config) {
	file, err := createLockFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer releaseLockFile(file)
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	autosaveTicker := time.NewTicker(time.Duration(cfg.AutoSaveIntervalMinutes) * time.Minute)
	defer autosaveTicker.Stop()

	monitorTicker := time.NewTicker(10 * time.Second)
	defer monitorTicker.Stop()

	for {
		select {
		case sig := <-signals:
			fmt.Printf("Received signal '%s'. Shutting down daemon...\n", sig.String())
			return
		case <-monitorTicker.C:
			if !isTmuxServerRunning() {
				fmt.Println("Tmux server is not running. Shutting down daemon.")
				saveSessions()
				return
			}
		case <-autosaveTicker.C:
			saveSessions()
		}
	}
}

func isTmuxServerRunning() bool {
	cmd := exec.Command("tmux", "list-sessions")
	err := cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func saveSessions() {
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		tmux.SendMsg("Failed to list tmux sessions.")
		return
	}

	savedSessions, err := session.LoadSessionsFromDisk()
	if err != nil {
		tmux.SendMsg("Failed to load sessions from disk.")
		return
	}

	combinedSessions, err := session.CombineSessions(tmuxSessions, savedSessions)
	if err != nil {
		tmux.SendMsg("Failed to combine sessions.")
		return
	}

	err = session.SaveSessionsToDisk(combinedSessions)
	if err != nil {
		tmux.SendMsg("Failed to save sessions to disk.")
	} else {
		tmux.SendMsg("Sessions saved successfully.")
	}
}

const lockFileName = "go-tms.lock"

func createLockFile() (*os.File, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not get user config directory: %w", err)
	}

	lockFilePath := filepath.Join(homePath, ".tmux", "go-tms", lockFileName)

	if err := os.MkdirAll(filepath.Dir(lockFilePath), 0755); err != nil {
		return nil, fmt.Errorf("could not create config directory: %w", err)
	}

	file, err := os.OpenFile(lockFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open lock file: %w", err)
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		_ = file.Close()
		return nil, fmt.Errorf("another instance of the daemon is already running")
	}

	return file, nil
}

func releaseLockFile(file *os.File) {
	if file != nil {
		fmt.Println("Releasing daemon lock.")
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
	}
}
