package pkg

import (
	"fmt"

	// "github.com/go-ldap/ldap/v3"
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

	// // search for the DN for the given username
	// // if found, try binding with the DN and user's password.
	// // if the binding succeeds, return the DN.
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.BaseDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(req.GetUsername())),
	// 	m.defaultUserFields(),
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	return nil, err
	// }
	// if len(result.Entries) != 1 {
	// 	return nil, &ZeroOrMultipleUsersError{
	// 		Username: req.GetUsername(),
	// 		Count:    len(result.Entries),
	// 	}
	// }
	// Make sure to always re-bind as admin afterwards
	defer m.BindAdmin()

	userDN := m.UserNamed(username)
	// userDN := fmt.Sprintf("cn=admin,%s", m.OpenLDAPConfig.BaseDN)
	// userDN := result.Entries[0].DN
	if err := m.ldap.Bind(userDN, password); err != nil {
		return nil, fmt.Errorf("unable to bind as %q", username)
	}
	return user, nil
	// return result.Entries[0], nil
}
