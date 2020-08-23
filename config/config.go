package config

import (
	"fmt"
	// tc "github.com/romnnn/testcontainers"
)

// OpenLDAPConfig ...
type OpenLDAPConfig struct {
	Host     string
	Port     int
	Protocol string

	Organization         string
	Domain               string
	BaseDN               string
	AdminPassword        string
	ConfigPassword       string
	ReadonlyUser         bool
	ReadonlyUserUsername string
	ReadonlyUserPassword string
	TLS                  bool
	UseRFC2307BISSchema  bool
}

// NewOpenLDAPConfig ...
func NewOpenLDAPConfig() OpenLDAPConfig {
	// populates default OpenLDAP config values
	return OpenLDAPConfig{
		Host:                 "localhost",
		Port:                 389,
		Protocol:             "ldap",
		Organization:         "Example Inc.",
		Domain:               "example.org",
		BaseDN:               "dc=example,dc=org",
		AdminPassword:        "admin",
		ConfigPassword:       "config",
		ReadonlyUser:         true,
		ReadonlyUserUsername: "readonly",
		ReadonlyUserPassword: "readonly",
		TLS:                  false,
		UseRFC2307BISSchema:  true,
	}
}

// URI ...
func (cfg *OpenLDAPConfig) URI() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Host, cfg.Port)
}
