package config

import ()

type WildApricot struct {
	CardNumber string `conf:"card-number"`
}

func NewWildApricot() *WildApricot {
	return &WildApricot{
		CardNumber: "Card Number",
	}
}
