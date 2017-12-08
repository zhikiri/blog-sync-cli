package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Setup app configuration parameters
type Setup struct {
	SourcePath string `json:"source_path"`
	Ignore     struct {
		Ext []string `json:"ext"`
	} `json:"ignore"`
	AWS struct {
		AccessKey string `json:"access_key"`
		SecretKey string `json:"secret_key"`
		Bucket    string `json:"bucket"`
		Region    string `json:"region"`
	} `json:"aws"`
}

func (s *Setup) setAccessKey(v string) error {
	s.AWS.AccessKey = v
	return nil
}

func (s *Setup) setSecretKey(v string) error {
	s.AWS.SecretKey = v
	return nil
}

func (s *Setup) setBucket(v string) error {
	s.AWS.Bucket = v
	return nil
}

func (s *Setup) setRegion(v string) error {
	s.AWS.Region = v
	return nil
}

func (s *Setup) setIgnoreExtensions(v string) error {
	ext := strings.Split(v, ",")
	for i, e := range ext {
		ext[i] = strings.Replace(strings.Trim(e, " "), ".", "", -1)
	}
	s.Ignore.Ext = ext
	return nil
}

func (s *Setup) setSourcePath(v string) error {
	_, err := os.Stat(v)
	if os.IsNotExist(err) {
		return errors.New("blog path doesn't exist")
	}
	s.SourcePath = v
	return nil
}

func (s *Setup) load() error {
	fmt.Println("Loading your config...")
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	// Read file content
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	// Set configuration properties
	return json.Unmarshal(f, s)
}

func (s *Setup) save() error {
	fmt.Println("Saving your config...")
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	// Creare a folder
	err = os.MkdirAll(filepath.Dir(path), 0766)
	if err != nil {
		return err
	}
	// Create a file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	// Encode the struct
	r, err := s.marshal()
	if err != nil {
		return err
	}
	// Save configuration
	_, err = io.Copy(f, r)
	return err
}

func (s *Setup) marshal() (io.Reader, error) {
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
