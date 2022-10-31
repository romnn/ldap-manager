package pkg

import (
  "fmt"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// AuthenticateUser authenticates a user
func (m *LDAPManager) AuthenticateUser(req *pb.LoginRequest) (*ldap.Entry, error) {
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, &ldaperror.ValidationError{Message: "must provide username and password"}
	}
	// search for the DN for the given username
	// if found, try binding with the DN and user's password.
	// if the binding succeeds, return the DN.
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(req.GetUsername())),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleUsersError{
			Username: req.GetUsername(),
			Count:    len(result.Entries),
		}
	}
	// Make sure to always re-bind as admin afterwards
	defer m.BindAdmin()
	userDN := result.Entries[0].DN
	if err := m.ldap.Bind(userDN, req.GetPassword()); err != nil {
		return nil, fmt.Errorf("unable to bind as %q", req.GetUsername())
	}
	return result.Entries[0], nil
}
