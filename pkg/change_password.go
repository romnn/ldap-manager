package pkg

import (
	"errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	// ldaphash "github.com/romnn/ldap-manager/pkg/hash"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// ChangePassword ...
func (m *LDAPManager) ChangePassword(req *pb.ChangePasswordRequest) error {
	// Validate
	username := req.GetUsername()
	password := req.GetPassword()
	if username == "" || password == "" {
		// todo: make this application error
		return errors.New("username and password must not be empty")
	}
	if req.GetHashingAlgorithm() == pb.HashingAlgorithm_DEFAULT {
		req.HashingAlgorithm = m.HashingAlgorithm
	}

	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(username)),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		// todo: make this application error
		return fmt.Errorf("failed to find existing user: %v", err)
	}
	if len(result.Entries) != 1 {
		return &ZeroOrMultipleUsersError{Username: username, Count: len(result.Entries)}
	}
	// TODO
	// userDN := result.Entries[0].DN
	// hashedPassword, err := ldaphash.Password(req.GetPassword(), req.GetHashingAlgorithm())
	// if err != nil {
	// 	return fmt.Errorf("failed to hash password: %v", err)
	// }
	// modifyPasswordRequest := ldap.NewModifyRequest(
	// 	userDN,
	// 	[]ldap.Control{},
	// )
	// modifyPasswordRequest.Replace("userPassword", []string{hashedPassword})
	// log.Debugf("modifyPasswordRequest=%v", modifyPasswordRequest)
	// if err := m.ldap.Modify(modifyPasswordRequest); err != nil {
	// 	return fmt.Errorf("failed to modify existing user: %v", err)
	// }
	log.Infof("changed password for user %q", username)
	return nil
}
