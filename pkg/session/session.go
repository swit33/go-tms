package session

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"slices"
)

type Pane struct {
	Command     string `yaml:"command"`
	CurrentPath string `yaml:"workdir"`
	Index       string `yaml:"index"`
}

type Window struct {
	Panes []Pane `yaml:"panes"`
	Index string `yaml:"index"`
}

type Session struct {
	Name        string   `yaml:"name"`
	Windows     []Window `yaml:"windows"`
	CurrentPath string   `yaml:"current-path"`
}

var sessionStorePath string = filepath.Join(".tmux", "go-tms", "sessions.yaml")

func GetSessionStorePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	sessionStoreAbsPath := filepath.Join(configDir, sessionStorePath)
	return sessionStoreAbsPath, nil
}

func SaveSessionsToDisk(sessions []Session) error {
	if len(sessions) == 0 {
		return nil
	}

	sessionStorePath, err := GetSessionStorePath()
	if err != nil {
		return err
	}

	sessionStoreDir := filepath.Dir(sessionStorePath)

	if err := os.MkdirAll(sessionStoreDir, 0755); err != nil {
		return err
	}

	file, err := os.Create(sessionStorePath)
	if err != nil {
		return err
	}
	defer file.Close()

	yamlEncoder := yaml.NewEncoder(file)
	yamlEncoder.SetIndent(2)
	if err = yamlEncoder.Encode(sessions); err != nil {
		return err
	}

	return nil
}

func LoadSessionsFromDisk() ([]Session, error) {
	sessionStorePath, err := GetSessionStorePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(sessionStorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Session{}, nil
		}
		return nil, err
	}
	defer file.Close()

	yamlDecoder := yaml.NewDecoder(file)
	var sessions []Session
	if err = yamlDecoder.Decode(&sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

func GetSessionByName(name string, s []Session) (*Session, error) {
	for _, session := range s {
		if session.Name == name {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

func GetSessionByPath(path string, s []Session) (*Session, error) {
	for _, session := range s {
		if session.CurrentPath == path {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("session not found")
}

func CombineSessions(s1 []Session, s2 []Session) ([]Session, error) {
	sessions := make([]Session, 0, len(s1)+len(s2))

	contains := func(list []Session, name string) bool {
		for _, sess := range list {
			if sess.Name == name {
				return true
			}
		}
		return false
	}

	for _, s := range s1 {
		sessions = append(sessions, s)
	}

	for _, s := range s2 {
		if !contains(s1, s.Name) {
			sessions = append(sessions, s)
		}
	}

	return sessions, nil
}

func DeleteSession(name string, s []Session) ([]Session, error) {
	for i, session := range s {
		if session.Name == name {
			return slices.Delete(s, i, i+1), nil
		}
	}
	return s, fmt.Errorf("session not found")
}

func CheckIfSessionExists(name string, s []Session) bool {
	for _, session := range s {
		if session.Name == name {
			return true
		}
	}
	return false
}
