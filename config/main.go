package config

import (
	"errors"
	"fmt"
)

// GetSetup function for initialize app configuration
func GetSetup() (*Setup, error) {
	if err := isConfigExist(); err != nil {
		fmt.Println("Configuration is missing")
		s, err := CreateConfig()
		if err != nil {
			return s, err
		}
		return s, nil
	}
	return &Setup{}, nil
}

func isConfigExist() error {
	return errors.New("Configuration file is not exist")
}
