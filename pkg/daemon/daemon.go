package daemon

import (
	"fmt"
	"go-tms/pkg/config"
	"go-tms/pkg/session"
	"go-tms/pkg/tmux"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func RunDaemon() {
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
			// sendMsg("Autosaving sessions...")
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
