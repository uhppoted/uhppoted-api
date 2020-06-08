package config

import ()

type AWS struct {
	Credentials string `conf:"credentials"`
	Profile     string `conf:"profil"`
	Region      string  `conf:"region"`
}

func NewAWS() *AWS {
	return &AWS{
		Credentials: "",
		Profile: "default",
		Region:      "us-east-1",
	}
}

