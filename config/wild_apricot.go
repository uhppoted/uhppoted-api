package config

import (
	"time"
)

type WildApricot struct {
	HTTP struct {
		ClientTimeout time.Duration `conf:"client-timeout"`
	} `conf:"http"`

	Fields struct {
		CardNumber string `conf:"card-number"`
	} `conf:"fields"`

	DisplayOrder struct {
		Groups string `conf:"groups"`
		Doors  string `conf:"doors"`
	} `conf:"display-order"`
}

func NewWildApricot() *WildApricot {
	return &WildApricot{
		HTTP: struct {
			ClientTimeout time.Duration `conf:"client-timeout"`
		}{
			ClientTimeout: 15 * time.Second,
		},

		Fields: struct {
			CardNumber string `conf:"card-number"`
		}{
			CardNumber: "Card Number",
		},
	}
}
