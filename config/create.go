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
	ques, sets := getQuestionsAndSetters(&s)
	for i, que := range ques {
		fmt.Println(que + ": ")
		anwr, _ := reader.ReadString('\n')
		err := sets[i](strings.Replace(anwr, "\n", "", -1))
		if err != nil {
			return Setup{}, err
		}
	}

	return s, nil
}

func saveSetup() {

}

func getQuestionsAndSetters(s *Setup) ([]string, []func(string) error) {
	return []string{
			"Enter absolute path to the blog folder",
			"Enter AWS access key",
			"Enter AWS access secret key",
			"Enter AWS S3 bucket name",
			"Enter AWS S3 bucket region",
		},
		[]func(string) error{
			s.setSourcePath,
			s.setAccessKey,
			s.setSecretKey,
			s.setBucket,
			s.setRegion,
		}
}
