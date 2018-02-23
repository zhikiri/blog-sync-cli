package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// CreateConfig ask user questions about conguration
func CreateConfig() (*Setup, error) {
	s, err := createSetup()
	if err != nil {
		return &Setup{}, err
	}

	err = s.save()
	if err != nil {
		return &Setup{}, err
	}

	return &s, nil
}

func createSetup() (Setup, error) {
	s := Setup{}

	fmt.Printf("Configuration initialization process started.\n\n")
	reader := bufio.NewReader(os.Stdin)
	ques, sets := getQuestionsAndSetters(&s)

	for i, que := range ques {
		fmt.Println(que + ": ")
		asw, _ := reader.ReadString('\n')
		fmt.Println()
		err := sets[i](strings.Replace(asw, "\n", "", -1))
		if err != nil {
			return Setup{}, err
		}
	}
	fmt.Println()

	return s, nil
}

func getQuestionsAndSetters(s *Setup) ([]string, []func(string) error) {
	return []string{
			"Enter absolute path to the blog static content folder (public folder for hugo)",
			"Enter AWS access key",
			"Enter AWS access secret key",
			"Enter AWS S3 bucket name",
			"Enter AWS S3 bucket region",
			"Enter ignore extensions (separate by comma)",
		},
		[]func(string) error{
			s.setSourcePath,
			s.setAccessKey,
			s.setSecretKey,
			s.setBucket,
			s.setRegion,
			s.setIgnoreExtensions,
		}
}
