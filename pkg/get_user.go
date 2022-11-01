package pkg

import (
	// "errors"
	"fmt"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// A ZeroOrMultipleUsersError is returned when zero or multiple users are found
type ZeroOrMultipleUsersError struct {
	error
	Username string
	Count    int
}

func (e *ZeroOrMultipleUsersError) Error() string {
	if e.Count > 1 {
		return fmt.Sprintf("multiple (%d) accounts with username %q", e.Count, e.Username)
	}
	return fmt.Sprintf("no account with username %q", e.Username)
}

func (e *ZeroOrMultipleUsersError) StatusError() error {
	if e.Count > 1 {
		return status.Errorf(codes.Internal, e.Error())
	}
	return status.Errorf(codes.NotFound, e.Error())
}

// ParseUser parses an ldap entry as `UserData`
func ParseUser(entry *ldap.Entry) *pb.UserData {
	user := &pb.UserData{Data: make(map[string]string)}
	for _, attr := range entry.Attributes {
		user.Data[attr.Name] = entry.GetAttributeValue(attr.Name)
	}
	return user
}

func (m *LDAPManager) defaultUserFields() []string {
	return []string{
		m.AccountAttribute,
		"givenName",
		"sn",
		"cn",
		"displayName",
		"uidNumber",
		"gidNumber",
		"loginShell",
		"homeDirectory",
		"mail",
	}
}

// GetUser gets a user
func (m *LDAPManager) GetUser(username string) (*pb.UserData, error) {
	if username == "" {
		return nil, &ldaperror.ValidationError{
      Message: "username must not be empty",
    }
	}
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(username)),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, fmt.Errorf("failed to get user %q: %v", username, err)
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleUsersError{
			Username: username,
			Count:    len(result.Entries),
		}
	}
	return ParseUser(result.Entries[0]), nil
}
