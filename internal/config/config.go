package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_user_name"`
}

func Read() (Config, error) {
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("could not unmarshal config")
	}

	return cfg, nil
}

func Write(cfg Config) error {
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("could not marshal config")
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write file")
	}

	return nil
}

func (c *Config) SetUser(username string) error {
	c.CurrentUsername = username

	err := Write(*c)
	if err != nil {
		return err
	}

	return nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory")
	}

	confDir := filepath.Join(homeDir, configFileName)

	return confDir, nil
}
