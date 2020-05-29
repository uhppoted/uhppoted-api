package config

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"io/ioutil"
	"regexp"
	"strings"
)

type AWS struct {
	Credentials *credentials.Credentials `conf:"credentials"`
	Region      string                   `conf:"region"`
}

func NewAWS() *AWS {
	return &AWS{
		Credentials: nil,
		Region:      "us-east-1",
	}
}

func (a *AWS) UnmarshalConf(tag string, values map[string]string) (interface{}, error) {
	re := regexp.MustCompile(fmt.Sprintf(`^%s\.(.*)$`, tag))

	for key, value := range values {
		match := re.FindStringSubmatch(key)
		if len(match) > 1 {
			switch match[1] {
			case "region":
				a.Region = value

			case "credentials":
				if credentials, err := getAWSCredentials(value); err != nil {
					return a, err
				} else {
					a.Credentials = credentials
				}
			}
		}
	}

	return a, nil
}

func getAWSCredentials(file string) (*credentials.Credentials, error) {
	awsKeyID := ""
	awsSecret := ""

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`\[default\]\n+aws_access_key_id\s*=\s*(.*?)\n+aws_secret_access_key\s*=\s*(.*)`)
	if match := re.FindSubmatch(bytes); len(match) == 3 {
		awsKeyID = strings.TrimSpace(string(match[1]))
		awsSecret = strings.TrimSpace(string(match[2]))
	} else {
		re = regexp.MustCompile(`\[default\]\n+aws_secret_access_key\s*=\s*(.*?)\n+aws_access_key_id\s*=\s*(.*)`)
		if match := re.FindSubmatch(bytes); len(match) == 3 {
			awsSecret = strings.TrimSpace(string(match[1]))
			awsKeyID = strings.TrimSpace(string(match[2]))
		}
	}

	if awsKeyID == "" {
		return nil, fmt.Errorf("Invalid AWS credentials - missing 'aws_access_key_id'")
	}

	if awsSecret == "" {
		return nil, fmt.Errorf("Invalid AWS credentials - missing 'aws_secret_access_key'")
	}

	return credentials.NewStaticCredentials(awsKeyID, awsSecret, ""), nil
}
