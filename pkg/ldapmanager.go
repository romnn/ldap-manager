package pkg

import (
	// "crypto/tls"
	// "strings"

	"github.com/go-ldap/ldap/v3"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
)

// Version is incremented using bump2version
// const Version = "0.0.26"

// LDAPManager ...
type LDAPManager struct {
	ldapconfig.OpenLDAPConfig
	// this is the only thing to guard with a mutex?
	ldap *ldap.Conn // Client

	GroupsDN    string
	UserGroupDN string

	GroupsOU string
	UsersOU  string

	HashingAlgorithm  pb.HashingAlgorithm
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

// NewLDAPManager ...
func NewLDAPManager(config ldapconfig.OpenLDAPConfig) *LDAPManager {
	return &LDAPManager{
		OpenLDAPConfig:           config,
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

// Close ...
func (m *LDAPManager) Close() {
	if m.ldap != nil {
		// FIXME: This will panic if the connection was not established
		m.ldap.Close()
	}
}
