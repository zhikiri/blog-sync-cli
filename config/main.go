package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

// GetSetup function for initialize app configuration
func GetSetup() (*Setup, error) {
	s := &Setup{}
	exist, err := isConfigExist()
	if err != nil {
		return s, err
	}

	if exist {
		err := s.load()
		if err != nil {
			return &Setup{}, err
		}
		return s, nil
	}

	fmt.Println("Configuration is missing")
	s, err = CreateConfig()
	if err != nil {
		return s, err
	}
	return s, nil
}

func isConfigExist() (bool, error) {
	p, err := getConfigPath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(p)
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func getConfigPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	path := filepath.Join(u.HomeDir, ".bsync", "setup.json")
	return path, nil
}
