package ldapmanager

import (
	"crypto/tls"
	"strings"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// LDAPManager ...
type LDAPManager struct {
	ldap ldap.Client

	GroupsDN    string
	UserGroupDN string
	BaseDN      string

	GroupsOU string
	UsersOU  string

	AdminBindUsername string
	AdminBindPassword string

	ReadonlyBindUsername string
	ReadonlyBindPassword string

	DefaultUserGroup  string
	DefaultAdminGroup string
	DefaultUserShell  string

	GroupMembershipAttribute string
	AccountAttribute         string
	GroupAttribute           string

	GroupMembershipUsesUID bool
	UseNISSchema           bool
	RequireStartTLS        bool
}

// Close ...
func (m *LDAPManager) Close() {
	if m.ldap != nil {
		m.ldap.Close()
	}
}

// Setup ...
func (m *LDAPManager) Setup(uri string) error {
	var err error
	m.ldap, err = ldap.DialURL(uri)
	if err != nil {
		return err
	}

	// Check for TLS
	if strings.HasPrefix(uri, "ldaps:") || m.RequireStartTLS {
		if err := m.ldap.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			log.Warnf("failed to connect via TLS: %v", err)
			if m.RequireStartTLS {
				return err
			}
		}
	}

	// Bind as the admin user
	if err := m.BindAdmin(); err != nil {
		return err
	}
	if err := m.SetupLDAP(); err != nil {
		return err
	}
	return nil
}
