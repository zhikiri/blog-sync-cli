package main

import (
	"fmt"
	"os"

	aws "github.com/aws/aws-sdk-go/service/s3"
	"github.com/zhikiri/blog-sync-cli/config"
	s3 "github.com/zhikiri/blog-sync-cli/storage"
)

func main() {
	s, err := config.GetSetup()
	showError(err)

	con, err := s3.Connect(s)
	showError(err)

	synchronize(s, con)
}

func synchronize(s *config.Setup, con *aws.S3) {
	che, err := s3.GetStorageChecksum(con, s.AWS.Bucket)
	showError(err)

	ver, errs := syncChangedFiles(s, con, che)
	showErrorList(errs)

	errs = syncAddedFiles(s, con, ver)
	showErrorList(errs)
}

func showErrorList(errs []error) {
	for _, err := range errs {
		fmt.Printf("Error: %s", err.Error())
	}
}

func showError(err error) {
	if err != nil {
		fmt.Printf("\nError: %s\n", err.Error())
		os.Exit(1)
	}
}
