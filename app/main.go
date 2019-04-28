package main

import (
	"log"
	"os"
	"path/filepath"
	"flag"

	"github.com/zhikiri/blog-sync-cli/app/config"
	"github.com/zhikiri/blog-sync-cli/app/storage"
	"github.com/zhikiri/blog-sync-cli/app/synchronizer"
)

func main() {
	settings, err := config.GetSettings(getConfigFilePath())
	failOnErr(err, "Settings parsing failed")

	access := config.GetAccessKey()
	secret := config.GetAccessSecret()

	storage, err := storage.NewAWS(settings, access, secret)
	failOnErr(err, "Storage initialization failed")

	err = synchronizer.SyncWith(settings, storage)
	failOnErr(err, "Storage synchronization failed")
}

func getConfigFilePath() string {
	configPath := flag.String("config", "../config.json", "path to the configuration file")
	flag.Parse()

	path, err := filepath.Abs(*configPath)
	failOnErr(err, "Configuration path cannot be retrieved")
	return path
}

func failOnErr(err error, message string) {
	if err != nil {
		log.Printf("[ERROR] Message: %s. %+v", message, err)
		os.Exit(1)
	}
}
