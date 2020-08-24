package ldapmanager

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap"
	"github.com/neko-neko/echo-logrus/v2/log"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	ldaphash "github.com/romnnn/ldap-manager/hash"
)

// ChangePassword ...
func (m *LDAPManager) ChangePassword(req *pb.ChangePasswordRequest) error {
	// Validate
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return errors.New("username and password must not be empty")
	}
	if req.GetHashingAlgorithm() == pb.HashingAlgorithm_DEFAULT {
		req.HashingAlgorithm = m.HashingAlgorithm
	}

	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(req.GetUsername())),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return fmt.Errorf("failed to find existing user: %v", err)
	}
	if len(result.Entries) != 1 {
		return &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	}
	userDN := result.Entries[0].DN
	hashedPassword, err := ldaphash.Password(req.GetPassword(), req.GetHashingAlgorithm())
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	modifyPasswordRequest := ldap.NewModifyRequest(
		userDN,
		[]ldap.Control{},
	)
	modifyPasswordRequest.Replace("userPassword", []string{hashedPassword})
	log.Debug(modifyPasswordRequest)
	if err := m.ldap.Modify(modifyPasswordRequest); err != nil {
		return fmt.Errorf("failed to modify existing user: %v", err)
	}
	log.Infof("changed password for user %q", req.GetUsername())
	return nil
}
