package main

import (
	"log"
	"os"
	"flag"

	"github.com/zhikiri/blog-sync-cli/app/config"
	"github.com/zhikiri/blog-sync-cli/app/storage"
	"github.com/zhikiri/blog-sync-cli/app/synchronizer"
)

func main() {
	configPath := flag.String("c", "../config.json", "configuration file path ('env' for environment)")
	flag.Parse()

	settings, err := config.GetSettings(*configPath)
	failOnErr(err, "Settings parsing failed")

	access := config.GetAccessKey()
	secret := config.GetAccessSecret()

	storage, err := storage.NewAWS(settings, access, secret)
	failOnErr(err, "Storage initialization failed")

	err = synchronizer.SyncWith(settings, storage)
	failOnErr(err, "Storage synchronization failed")
}

func failOnErr(err error, message string) {
	if err != nil {
		log.Printf("[ERROR] Message: %s. %+v", message, err)
		os.Exit(1)
	}
}
