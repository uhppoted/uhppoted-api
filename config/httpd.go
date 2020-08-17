package config

import ()

type HTTPD struct {
	AuthDB                    string `conf:"auth.local.db"`
	CookieMaxAge              int    `conf:"cookie.max-age"`
	LoginExpiry             string `conf:"login.expiry"`
	SessionExpiry             string `conf:"session.expiry"`
	HttpEnabled               bool   `conf:"http.enabled"`
	HttpPort                  uint16 `conf:"http.port"`
	HttpsEnabled              bool   `conf:"https.enabled"`
	HttpsPort                 uint16 `conf:"https.port"`
	CACertificate             string `conf:"tls.ca"`
	TLSCertificate            string `conf:"tls.certificate"`
	TLSKey                    string `conf:"tls.key"`
	RequireClientCertificate bool   `conf:"tls.client.certificates.required"`
}

func NewHTTPD() *HTTPD {
	return &HTTPD{
		AuthDB:                    httpdAuthDB,
		CookieMaxAge:              24,
		LoginExpiry:             "5m",
		SessionExpiry:             "60m",
		HttpEnabled:               false,
		HttpsEnabled:              true,
		CACertificate:             httpdCACertificate,
		TLSCertificate:            httpdTLSCertificate,
		TLSKey:                    httpdTLSKey,
		RequireClientCertificate: false,
	}
}
