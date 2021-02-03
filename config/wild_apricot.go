package config

import ()

type WildApricot struct {
	CardNumber        string `conf:"card-number"`
	GroupDisplayOrder string `conf:"groups"`
}

func NewWildApricot() *WildApricot {
	return &WildApricot{
		CardNumber: "Card Number",
	}
}
