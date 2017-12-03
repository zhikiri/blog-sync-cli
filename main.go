package main

import (
	"fmt"
	"os"

	"github.com/zhikiri/blog-sync-cli/config"
)

func main() {
	fmt.Println("Started...")
	s, err := config.GetSetup()
	showUserError(err)

	fmt.Printf("%v", s)
}

func showUserError(err error) {
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
		os.Exit(1)
	}
}
