package daemon

import (
	"fmt"
	"go-tms/pkg/config"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

func RunDaemon() {
	file, err := createLockFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer releaseLockFile(file)
	cfg, err := config.LoadConfig()
	if err != nil {
		sendMsg("Failed to start auto-save daemon: " + err.Error())
		return
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(cfg.AutoSaveIntervalMinutes) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case sig := <-signals:
			fmt.Printf("Received signal '%s'. Shutting down daemon...\n", sig.String())
			return
		case <-ticker.C:
			saveSessions()
		}
	}
}

func saveSessions() {
	tmuxSessions, err := tmux.ListSessions()
	if err != nil {
		sendMsg("Failed to list tmux sessions.")
		return
	}

	savedSessions, err := session.LoadSessionsFromDisk()
	if err != nil {
		sendMsg("Failed to load sessions from disk.")
		return
	}

	combinedSessions, err := session.CombineSessions(tmuxSessions, savedSessions)
	if err != nil {
		sendMsg("Failed to combine sessions.")
		return
	}

	err = session.SaveSessionsToDisk(combinedSessions)
	if err != nil {
		sendMsg("Failed to save sessions to disk.")
	} else {
		sendMsg("Sessions saved successfully.")
	}
}

func sendMsg(msg string) {
	cmd := exec.Command("tmux", "display-message", msg)
	cmd.Run()
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
