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

	LDAPOrganization         string
	LDAPDomain               string
	LDAPBaseDN               string
	LDAPAdminPassword        string
	LDAPConfigPassword       string
	LDAPReadonlyUser         bool
	LDAPReadonlyUserUsername string
	LDAPReadonlyUserPassword string
	LDAPTLS                  bool
	LDAPRFC2307BISSchema     bool
}

// NewOpenLDAPConfig ...
func NewOpenLDAPConfig() OpenLDAPConfig {
	// populates default OpenLDAP config values
	return OpenLDAPConfig{
		Host:                     "localhost",
		Port:                     389,
		Protocol:                 "ldap",
		LDAPOrganization:         "Example Inc.",
		LDAPDomain:               "example.org",
		LDAPBaseDN:               "dc=example,dc=org",
		LDAPAdminPassword:        "admin",
		LDAPConfigPassword:       "config",
		LDAPReadonlyUser:         true,
		LDAPReadonlyUserUsername: "readonly",
		LDAPReadonlyUserPassword: "readonly",
		LDAPTLS:                  false,
		LDAPRFC2307BISSchema:     true,
	}
}

// URI ...
func (cfg *OpenLDAPConfig) URI() string {
	return fmt.Sprintf("%s://%s:%d", cfg.Protocol, cfg.Host, cfg.Port)
}
