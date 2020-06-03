package config

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws/credentials"
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
				a.Credentials = credentials.NewSharedCredentials(value, "default")
			}
		}
	}

	return a, nil
}
