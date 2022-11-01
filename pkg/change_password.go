package pkg

import (
	// "errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// ChangePassword changes the password of a user
func (m *LDAPManager) ChangePassword(req *pb.ChangePasswordRequest) error {
	username := req.GetUsername()
	if username == "" {
		return &ldaperror.ValidationError{Message: "username must not be empty"}
	}
	password := req.GetPassword()
	if password == "" {
		return &ldaperror.ValidationError{Message: "password must not be empty"}
	}
	user, err := m.GetUser(username)
	log.Infof("user: %v", user)
	if err != nil {
		return err
	}
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.BaseDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(username)),
	// 	[]string{"dn"},
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// return
	// return &ZeroOrMultipleUsersError{
	// Username: username,
	// Count: len(result.Entries),
	// }
	// }
	// if len(result.Entries) != 1 {
	// 	return &ZeroOrMultipleUsersError{
	// Username: username,
	// Count: len(result.Entries),
	// }
	// }
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
	// // log.Debugf("modifyPasswordRequest=%v", modifyPasswordRequest)
	// if err := m.ldap.Modify(modifyPasswordRequest); err != nil {
	// 	return fmt.Errorf("failed to modify existing user: %v", err)
	// }

	// change password of user
	passwordModifyRequest := &ldap.PasswordModifyRequest{
		UserIdentity: username,
		NewPassword:  password,
	}
	// log.Infof("passwordModifyRequest=%v", passwordModifyRequest)
	_, err = m.ldap.PasswordModify(passwordModifyRequest)
	if err != nil {
		return fmt.Errorf("failed to set password of user %q: %v", username, err)
	}

	// log.Infof("changed password for user %q", username)
	return nil
}
