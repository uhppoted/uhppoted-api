package config

import ()

type WildApricot struct {
	HTTP struct {
		Timeout uint `conf:"client-timeout"`
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
			Timeout uint `conf:"client-timeout"`
		}{
			Timeout: 15,
		},

		Fields: struct {
			CardNumber string `conf:"card-number"`
		}{
			CardNumber: "Card Number",
		},
	}
}
