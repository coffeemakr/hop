package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type ClientConfiguration struct {
	Group   string `json:"group,omitempty"`
	BaseURL string `json:"url,omitempty"`
	Proxy   string `json:"proxy,omitempty"`
}

func GetConfigFolderPath() (string, error) {
	if runtime.GOOS == "windows" {
		appDir := os.Getenv("APPDATA")
		if appDir == "" {
			return "", errors.New("APPDATA not found")
		}
		return path.Join(appDir, "ruck"), nil
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(homedir, ".config", "ruck"), nil
}

func GetConfigPath() (string, error) {
	folder, err := GetConfigFolderPath()
	if err != nil {
		return "", err
	}
	return path.Join(folder, "config.json"), nil
}

func WriteConfig(config *ClientConfiguration) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}
	directory := filepath.Dir(path)
	os.MkdirAll(directory, os.FileMode(0700))

	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	encoder := json.NewEncoder(fp)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(config)
	return err
}

func LoadConfig() (*ClientConfiguration, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}
	_, err = os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		// create emtpy configuration
		return &ClientConfiguration{}, nil
	} else if err != nil {
		return nil, err
	}

	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var config ClientConfiguration
	err = json.NewDecoder(fp).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func PrintConfig(config *ClientConfiguration, key string) error {
	switch key {
	case "group":
		fmt.Println(config.Group)
	case "url":
		fmt.Println(config.BaseURL)
	case "proxy":
		fmt.Println(config.Proxy)
	default:
		return errors.New("Unknown key")
	}
	return nil
}

func SetConfig(config *ClientConfiguration, key string, value string) error {
	switch key {
	case "group":
		config.Group = value
	case "url":
		config.BaseURL = value
	case "proxy":
		config.Proxy = value
	default:
		return errors.New("Unknown key")
	}
	return nil
}
