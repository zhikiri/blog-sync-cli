package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zhikiri/blog-sync-cli/config"
)

func main() {
	s, err := config.GetSetup()
	showUserError(err)

	fmt.Printf("%v", s)
}

func showUserError(err error) {
	if err != nil {
		log.Panic(err)
		//fmt.Printf("\nError: %s\n", err.Error())
		os.Exit(1)
	}
}
