package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"path/filepath"
)

var configPath string = filepath.Join(".config", "go-tms", "config.yaml")

type Config struct {
	AutoSaveEnabled         bool `yaml:"auto-save"`
	AutoSaveIntervalMinutes int  `yaml:"auto-save-interval-minutes"`
}

func GetConfigPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	configAbsPath := filepath.Join(currentUser.HomeDir, configPath)
	return configAbsPath, nil
}

func LoadConfig() (Config, error) {
	config := Config{
		AutoSaveEnabled:         true,
		AutoSaveIntervalMinutes: 5,
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
