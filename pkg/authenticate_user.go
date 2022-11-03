package pkg

import (
	"fmt"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// AuthenticateUser authenticates a user
func (m *LDAPManager) AuthenticateUser(req *pb.LoginRequest) (*pb.User, error) {
	username := req.GetUsername()
	password := req.GetPassword()
	if username == "" {
		return nil, &ldaperror.ValidationError{
			Message: "username must not be empty",
		}
	}
	if password == "" {
		return nil, &ldaperror.ValidationError{
			Message: "password must not be empty",
		}
	}

	user, err := m.GetUser(username)
	if err != nil {
		return nil, err
	}
	// re-bind as admin afterwards
	defer m.BindAdmin()

	if err := m.ldap.Bind(m.UserDN(username), password); err != nil {
		return nil, fmt.Errorf("unable to bind as %q", username)
	}
	return user, nil
}
