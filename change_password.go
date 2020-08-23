package ldapmanager

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap"
	ldaphash "github.com/romnnn/ldap-manager/hash"
	log "github.com/sirupsen/logrus"
)

// ChangePassword ...
func (m *LDAPManager) ChangePassword(username, newPassword string, algorithm ldaphash.LDAPPasswordHashingAlgorithm) error {
	// Validate
	if username == "" || newPassword == "" {
		return errors.New("username and password must not be empty")
	}

	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", m.AccountAttribute, escape(username), m.UserGroupDN),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return err
	}
	if len(result.Entries) != 1 {
		return fmt.Errorf("zero or multiple (%d) accounts with username %q", len(result.Entries), username)
	}
	userDN := result.Entries[0].DN
	hashedPassword, err := ldaphash.Password(newPassword, algorithm)
	if err != nil {
		return err
	}
	modifyPasswordRequest := ldap.NewModifyRequest(
		userDN,
		[]ldap.Control{},
	)
	modifyPasswordRequest.Replace("userPassword", []string{hashedPassword})
	log.Debug(modifyPasswordRequest)
	if err := m.ldap.Modify(modifyPasswordRequest); err != nil {
		return err
	}
	log.Infof("changed password for user %q", username)
	return nil
}
