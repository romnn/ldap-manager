package config

import (
	"fmt"
)

// Config contains the configuration of the LDAP server
type Config struct {
	Host     string
	Port     int
	Protocol string

	Organization        string
	Domain              string
	BaseDN              string
	AdminUsername       string
	AdminPassword       string
	ReadOnlyUser        bool
	ReadOnlyUsername    string
	ReadOnlyPassword    string
	ConfigPassword      string
	TLS                 bool
	UseRFC2307BISSchema bool
}

// NewConfig creates a default LDAP configuration
func NewConfig() Config {
	return Config{
		Host:                "localhost",
		Port:                389,
		Protocol:            "ldap",
		Organization:        "Example Inc.",
		Domain:              "example.org",
		BaseDN:              "dc=example,dc=org",
		AdminUsername:       "admin",
		AdminPassword:       "admin",
		ReadOnlyUser:        true,
		ReadOnlyUsername:    "readonly",
		ReadOnlyPassword:    "readonly",
		TLS:                 false,
		UseRFC2307BISSchema: true,
	}
}

// URI returns the connection URI for the LDAP config
func (cfg *Config) URI() string {
	return fmt.Sprintf(
		"%s://%s:%d",
		cfg.Protocol, cfg.Host, cfg.Port,
	)
}
