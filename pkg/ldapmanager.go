package pkg

import (
	"github.com/go-ldap/ldap/v3"
	"github.com/jwalton/go-supportscolor"
	"github.com/k0kubun/pp/v3"
	log "github.com/sirupsen/logrus"

	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
)

// LDAPManager implements the LDAP manager functionality
type LDAPManager struct {
	ldapconfig.Config
	ldap *ldap.Conn

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
	useColor := supportscolor.Stdout().SupportsColor
	pp.Default.SetColoringEnabled(useColor)
	pp.Default.SetExportedOnly(true)

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
	if m.ldap != nil {
		// FIXME: This will panic if the connection was not established
		m.ldap.Close()
	}
}
