package config

import (
	"fmt"
)

// Config contains the configuration of the LDAP server
type Config struct {
	Host     string
	Port     int
	Protocol string

	Organization         string
	Domain               string
	BaseDN               string
	AdminPassword        string
	ReadonlyUser         bool
	ReadonlyUserUsername string
	ReadonlyUserPassword string
	TLS                  bool
	UseRFC2307BISSchema  bool
}

// NewConfig creates a default LDAP configuration
func NewConfig() Config {
	return Config{
		Host:                 "localhost",
		Port:                 389,
		Protocol:             "ldap",
		Organization:         "Example Inc.",
		Domain:               "example.org",
		BaseDN:               "dc=example,dc=org",
		AdminPassword:        "admin",
		ReadonlyUser:         true,
		ReadonlyUserUsername: "readonly",
		ReadonlyUserPassword: "readonly",
		TLS:                  false,
		UseRFC2307BISSchema:  true,
	}
}

// URI returns the connection URI for the LDAP config
func (cfg *Config) URI() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Host, cfg.Port)
}
