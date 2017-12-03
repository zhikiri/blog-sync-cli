package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CreateConfig ask user questions about conguration
func CreateConfig() (*Setup, error) {
	s, err := createSetupStruct()
	if err != nil {
		return &Setup{}, err
	}
	return &s, err
}

func createSetupStruct() (Setup, error) {
	s := Setup{}

	fmt.Printf("Configuration initialization process started.\n\n")
	reader := bufio.NewReader(os.Stdin)
	for q, setter := range getQuestions(&s) {
		fmt.Println(q)
		anwr, _ := reader.ReadString('\n')
		err := setter(strings.Replace(anwr, "\n", "", -1))
		if err != nil {
			return Setup{}, err
		}
	}

	return s, nil
}

func saveSetup() {

}

func getQuestions(s *Setup) map[string]func(string) error {
	return map[string]func(string) error{
		"Enter absolute path to the blog folder: ": s.setSourcePath,
		"Enter AWS access key: ":                   s.setAccessKey,
		"Enter AWS access secret key: ":            s.setSecretKey,
		"Enter AWS S3 bucket name: ":               s.setBucket,
		"Enter AWS S3 bucket region: ":             s.setRegion,
	}
}
