package config

import (
	"path/filepath"
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
	return os.Getenv("AWS_ACCESS_KEY_ID")
}

// GetAccessSecret function that load AWS access secret from env
func GetAccessSecret() string {
	return os.Getenv("AWS_SECRET_ACCESS_KEY")
}

// GetSettings function for load configuration
func GetSettings(config string) (Settings, error) {
	if config == "env" {
		return getSettingsFromEnv()
	}
	path, err := getConfigFilePath(config)
	if err != nil {
		return Settings{}, err
	}
	return getSettingsFromFile(path)
}

func getSettingsFromEnv() (Settings, error) {
	source := os.Getenv("WEBSITE_SOURCE_PATH")
	bucket := os.Getenv("WEBSITE_BUCKET")
	region := os.Getenv("WEBSITE_REGION")

	if source == "" || bucket == "" || region == "" {
		return Settings{}, errors.New("Env variables are not set properly")
	}
	return Settings{source, make([]string, 0), bucket, region}, nil
}

func getConfigFilePath(config string) (string, error) {
	path, err := filepath.Abs(config)
	if err != nil {
		return "", err
	}
	return path, nil
}

func getSettingsFromFile(path string) (Settings, error) {
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
