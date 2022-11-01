package pkg

import (
	// "errors"
	"fmt"
	"strconv"

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

const (
	uidNumber     = "uidNumber"
	gidNumber     = "gidNumber"
	givenName     = "givenName"
	displayName   = "displayName"
	loginShell    = "loginShell"
	homeDirectory = "homeDirectory"
	mail          = "mail"
	sn            = "sn"
	cn            = "cn"
)

// ParseUser parses an ldap entry as a User
func (m *LDAPManager) ParseUser(entry *ldap.Entry) *pb.User {
	// user := &pb.User{Data: make(map[string]string)}
	// for _, attr := range entry.Attributes {
	// 	user.Data[attr.Name] = entry.GetAttributeValue(attr.Name)
	// }

	// string first_name = 1;
	// string last_name = 2;

	// int32 uid = 10;
	// int64 gid = 11;
	// string login_shell = 12;
	// string home_directory = 13;

	// string username = 20;
	// string email = 21;

	uid, _ := strconv.Atoi(entry.GetAttributeValue(uidNumber))
	gid, _ := strconv.Atoi(entry.GetAttributeValue(gidNumber))
	return &pb.User{
		Username:      entry.GetAttributeValue(m.AccountAttribute),
		FirstName:     entry.GetAttributeValue(givenName),
		LastName:      entry.GetAttributeValue(sn),
		CN:            entry.GetAttributeValue(cn),
		DisplayName:   entry.GetAttributeValue(displayName),
		UID:           int32(uid),
		GID:           int64(gid),
		LoginShell:    entry.GetAttributeValue(loginShell),
		HomeDirectory: entry.GetAttributeValue(homeDirectory),
		Email:         entry.GetAttributeValue(mail),
	}
}

func (m *LDAPManager) userFields() []string {
  return []string{
			m.AccountAttribute,
			givenName,
			sn,
			cn,
			displayName,
			uidNumber,
			gidNumber,
			loginShell,
			homeDirectory,
			mail,
		}
}

// GetUser gets a user
func (m *LDAPManager) GetUser(username string) (*pb.User, error) {
	if username == "" {
		return nil, &ldaperror.ValidationError{
			Message: "username must not be empty",
		}
	}
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(username)),
		m.userFields(),
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
	return m.ParseUser(result.Entries[0]), nil
}
