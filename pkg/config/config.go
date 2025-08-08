package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

var configPath string = filepath.Join(".config", "go-tms", "config.yaml")

type Config struct {
	AutoSaveIntervalMinutes int    `yaml:"auto-save-interval-minutes"`
	FZFBindNew              string `yaml:"fzf-bind-new"`
	FZFBindDelete           string `yaml:"fzf-bind-delete"`
	FZFBindInteractive      string `yaml:"fzf-bind-interactive"`
	FZFBindSave             string `yaml:"fzf-bind-save"`
	FZFPrompt               string `yaml:"fzf-prompt"`
	FZFOpts                 string `yaml:"fzf-env"`
	ZoxideOpts              string `yaml:"zoxide-env"`
}

func GetConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	configAbsPath := filepath.Join(configDir, configPath)
	return configAbsPath, nil
}

func LoadConfig() (Config, error) {
	config := Config{
		AutoSaveIntervalMinutes: 10,
		FZFBindNew:              "ctrl-n",
		FZFBindDelete:           "ctrl-d",
		FZFBindInteractive:      "ctrl-i",
		FZFBindSave:             "ctrl-s",
		FZFPrompt:               "Sessions> ",
		FZFOpts:                 "--no-sort --reverse",
		ZoxideOpts:              "--layout=reverse --style=full --border=bold --border=rounded --margin=3%",
	}

	configFilePath, err := GetConfigPath()
	if err != nil {
		return config, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		// If the file does not exist, return the default config without creating the file.
		if os.IsNotExist(err) {
			return config, nil
		}
		// Return other errors.
		return config, err
	}
	defer configFile.Close()

	yamlDecoder := yaml.NewDecoder(configFile)
	if err = yamlDecoder.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
