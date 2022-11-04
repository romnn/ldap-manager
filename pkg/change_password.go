package pkg

import (
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
		return &ldaperror.ValidationError{
			Message: "username must not be empty",
		}
	}
	password := req.GetPassword()
	if password == "" {
		return &ldaperror.ValidationError{
			Message: "password must not be empty",
		}
	}

	modifyReq := ldap.PasswordModifyRequest{
		UserIdentity: m.UserDN(username),
		NewPassword:  password,
	}
	log.Debug(PrettyPrint(modifyReq))

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err := conn.PasswordModify(&modifyReq); err != nil {
		notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		if notFound {
			return &ZeroOrMultipleUsersError{
				Username: username,
				Count:    0,
			}
		}
		return fmt.Errorf(
			"failed to set password of user %q: %v",
			username, err,
		)
	}
	log.Infof("changed password for user %q", username)
	return nil
}
