package pkg

import (
	log "github.com/sirupsen/logrus"

	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	ldappool "github.com/romnn/ldap-manager/pkg/pool"
)

// LDAPManager implements the LDAP manager functionality
type LDAPManager struct {
	ldapconfig.Config
	Pool ldappool.Pool

	GroupsDN    string
	UserGroupDN string

	GroupsOU string
	UsersOU  string

	DefaultUserGroup  string
	DefaultAdminGroup string
	DefaultUserShell  string

	DefaultAdminUsername string
	DefaultAdminPassword string
	ForceCreateAdmin     bool

	GroupMembershipAttribute string
	AccountAttribute         string
	GroupAttribute           string

	GroupMembershipUsesUID bool
}

// NewLDAPManager creates a new LDAPManager
func NewLDAPManager(config ldapconfig.Config) *LDAPManager {
	log.SetFormatter(&log.TextFormatter{
		DisableQuote: true,
	})

	return &LDAPManager{
		Config:                   config,
		GroupsDN:                 "ou=groups," + config.BaseDN,
		UserGroupDN:              "ou=users," + config.BaseDN,
		GroupsOU:                 "groups",
		UsersOU:                  "users",
		DefaultUserGroup:         "users",
		DefaultAdminGroup:        "admins",
		DefaultUserShell:         "/bin/bash",
		GroupMembershipAttribute: "uniqueMember", // uniqueMember or memberUID
		AccountAttribute:         "uid",
		GroupAttribute:           "gid",
		GroupMembershipUsesUID:   false,
		DefaultAdminUsername:     "admin",
		DefaultAdminPassword:     "admin",
		ForceCreateAdmin:         false,
	}
}

// Close closes the LDAP connection
func (m *LDAPManager) Close() {
	if m.Pool != nil {
		m.Pool.Close()
	}
}
