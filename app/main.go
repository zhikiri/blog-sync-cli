package main

import (
	"flag"
	"log"
	"os"

	"github.com/zhikiri/bsync/app/config"
	"github.com/zhikiri/bsync/app/storage"
	"github.com/zhikiri/bsync/app/synchronizer"
)

func main() {
	configPath := flag.String("config", "../config.json", "configuration file path ('env' for environment)")
	isAuthRequired := flag.Bool("aws-auth", true, "is aws authorization required")
	flag.Parse()

	settings, err := config.GetSettings(*configPath)
	failOnErr(err, "Settings parsing failed")

	storage, err := getStorage(settings, *isAuthRequired)
	failOnErr(err, "Storage initialization failed")

	err = synchronizer.SyncWith(settings, storage)
	failOnErr(err, "Storage synchronization failed")
}

func getStorage(settings config.Settings, isAuthRequired bool) (storage.Storage, error) {
	if isAuthRequired {
		return storage.NewAWSAuth(settings, config.GetAccessKey(), config.GetAccessSecret())
	}
	return storage.NewAWS(settings, nil)
}

func failOnErr(err error, message string) {
	if err != nil {
		log.Printf("[ERROR] Message: %s. %+v", message, err)
		os.Exit(1)
	}
}
