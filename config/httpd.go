package config

import ()

type HTTPD struct {
	AuthDB        string `conf:"auth.local.db"`
	CookieMaxAge  int    `conf:"cookie.max-age"`
	SessionExpiry string `conf:"session.expiry"`
}

func NewHTTPD() *HTTPD {
	return &HTTPD{
		AuthDB:        httpdAuthDB,
		CookieMaxAge:  24,
		SessionExpiry: "60m",
	}
}
