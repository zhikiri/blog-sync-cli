package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Settings struct {
	Source string   `json:"src"`
	Ignore []string `json:"ignore"`
	Bucket string   `json:"bucket"`
	Region string   `json:"region"`
}

// GetAccessKey function that load AWS access key from env
func GetAccessKey() string {
	return os.Getenv("AWS_ACCESS_KEY")
}

// GetAccessSecret function that load AWS access secret from env
func GetAccessSecret() string {
	return os.Getenv("AWS_ACCESS_SECRET")
}

// GetSettings function for load configuration from file
func GetSettings(path string) (Settings, error) {
	exists, err := isFileExists(path)
	if err != nil {
		return Settings{}, err
	}
	if !exists {
		return Settings{}, errors.New("File is not exist")
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return Settings{}, err
	}
	settings := Settings{}
	err = json.Unmarshal(file, &settings)
	return settings, err
}

func isFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
