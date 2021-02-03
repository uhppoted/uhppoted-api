package config

import ()

type WildApricot struct {
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
		Fields: struct {
			CardNumber string `conf:"card-number"`
		}{
			CardNumber: "Card Number",
		},
	}
}
