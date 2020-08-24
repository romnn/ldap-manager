package ldapmanager

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap"
	ldaphash "github.com/romnnn/ldap-manager/hash"
	log "github.com/sirupsen/logrus"
)

// ChangePasswordRequest ...
type ChangePasswordRequest struct {
	Username  string `json:"username" form:"username"`
	Password  string `json:"password" form:"password"`
	Algorithm ldaphash.LDAPPasswordHashingAlgorithm
}

// ChangePassword ...
func (m *LDAPManager) ChangePassword(req *ChangePasswordRequest) error {
	// Validate
	if req.Username == "" || req.Password == "" {
		return errors.New("username and password must not be empty")
	}
	if req.Algorithm == ldaphash.DEFAULT {
		req.Algorithm = m.HashingAlgorithm
	}

	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", m.AccountAttribute, escapeFilter(req.Username), m.UserGroupDN),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return err
	}
	if len(result.Entries) != 1 {
		return &ZeroOrMultipleAccountsError{Username: req.Username, Count: len(result.Entries)}
	}
	userDN := result.Entries[0].DN
	hashedPassword, err := ldaphash.Password(req.Password, req.Algorithm)
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
	log.Infof("changed password for user %q", req.Username)
	return nil
}
