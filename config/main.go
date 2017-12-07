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
	if exist, _ := isConfigExist(); exist == true {
		if err := s.load(); err != nil {
			return &Setup{}, err
		}
		return s, nil
	}

	fmt.Println("Configuration is missing")

	s, err := CreateConfig()
	if err != nil {
		return &Setup{}, err
	}
	return s, nil
}

func isConfigExist() (bool, error) {
	p, err := getConfigPath()
	if err != nil {
		return false, err
	}

	if _, err = os.Stat(p); err == nil {
		return true, nil
	}
	return false, err
}

func getConfigPath() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	path := filepath.Join(u.HomeDir, ".bsync", "setup.json")
	return path, nil
}
