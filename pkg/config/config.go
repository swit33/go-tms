package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var configPath string = filepath.Join("go-tms", "config.yaml")

type Config struct {
	AutoSaveIntervalMinutes int    `yaml:"auto-save-interval-minutes"`
	FZFBindNew              string `yaml:"fzf-bind-new"`
	FZFBindDelete           string `yaml:"fzf-bind-delete"`
	FZFBindInteractive      string `yaml:"fzf-bind-interactive"`
	FZFBindSave             string `yaml:"fzf-bind-save"`
	FZFBindKill             string `yaml:"fzf-bind-kill"`
	FZFPrompt               string `yaml:"fzf-prompt"`
	FZFOpts                 string `yaml:"fzf-opts"`
	ZoxideOpts              string `yaml:"zoxide-opts"`
	ProgramWhitelist        string `yaml:"program-whitelist"`
	NvimCustomCommand       string `yaml:"nvim-custom-command"`
	SelectFirst             bool   `yaml:"select-first"`
	CloseOnNew              bool   `yaml:"close-on-new"`
	ActiveSessionPrefix     string `yaml:"active-session-prefix"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	xdgConfigPath := filepath.Join(homeDir, ".config", configPath)
	return xdgConfigPath, nil
}

func LoadConfig() (Config, error) {
	config := Config{
		AutoSaveIntervalMinutes: 10,
		FZFBindNew:              "ctrl-n",
		FZFBindDelete:           "ctrl-d",
		FZFBindInteractive:      "ctrl-i",
		FZFBindSave:             "ctrl-s",
		FZFBindKill:             "ctrl-k",
		FZFPrompt:               "Sessions> ",
		FZFOpts:                 "--no-sort --reverse",
		ZoxideOpts:              "--layout=reverse --style=full --border=bold --border=rounded --margin=3%",
		ProgramWhitelist:        "btop,vim,nvim,yazi",
		NvimCustomCommand:       "",
		SelectFirst:             true,
		CloseOnNew:              true,
		ActiveSessionPrefix:     " ",
	}

	configFilePath, err := getConfigPath()
	if err != nil {
		return config, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}
	defer configFile.Close()

	yamlDecoder := yaml.NewDecoder(configFile)
	if err = yamlDecoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
