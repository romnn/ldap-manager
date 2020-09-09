package ldapmanager

import (
	"crypto/tls"
	"strings"

	"github.com/go-ldap/ldap"
	ldapconfig "github.com/romnnn/ldap-manager/config"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	log "github.com/sirupsen/logrus"
)

// Version is incremented using bump2version
const Version = "0.0.19"

// LDAPManager ...
type LDAPManager struct {
	ldapconfig.OpenLDAPConfig
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
func NewLDAPManager(cfg ldapconfig.OpenLDAPConfig) *LDAPManager {
	return &LDAPManager{
		OpenLDAPConfig:           cfg,
		GroupsDN:                 "ou=groups," + cfg.BaseDN,
		UserGroupDN:              "ou=users," + cfg.BaseDN,
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

// Setup ...
func (m *LDAPManager) Setup(skipSetupLDAP bool) error {
	var err error
	URI := m.OpenLDAPConfig.URI()
	log.Debugf("connecting to OpenLDAP at %s", URI)
	m.ldap, err = ldap.DialURL(URI)
	if err != nil {
		return err
	}

	// Check for TLS
	if strings.HasPrefix(URI, "ldaps:") || m.OpenLDAPConfig.TLS {
		if err := m.ldap.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			log.Warnf("failed to connect via TLS: %v", err)
			if m.OpenLDAPConfig.TLS {
				return err
			}
		}
	}

	// Bind as the admin user
	if err := m.BindAdmin(); err != nil {
		return err
	}
	if !skipSetupLDAP {
		if err := m.SetupLDAP(); err != nil {
			return err
		}
	}
	return nil
}
