package pkg

import (
	"fmt"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
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

	conn, err := m.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.Bind(m.UserDN(username), password); err != nil {
		log.Errorf("unable to bind as %q: %v", username, err)
		return nil, fmt.Errorf("unable to bind as %q", username)
	}
	return user, nil
}
