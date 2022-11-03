package pkg

import (
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
		return fmt.Sprintf(
			"multiple (%d) accounts with username %q",
			e.Count, e.Username,
		)
	}
	return fmt.Sprintf(
		"no account with username %q",
		e.Username,
	)
}

func (e *ZeroOrMultipleUsersError) StatusError() error {
	if e.Count > 1 {
		return status.Errorf(codes.Internal, e.Error())
	}
	return status.Errorf(codes.NotFound, e.Error())
}

const (
	userUidNumber     = "uidNumber"
	userGidNumber     = "gidNumber"
	userGivenName     = "givenName"
	userDisplayName   = "displayName"
	userLoginShell    = "loginShell"
	userHomeDirectory = "homeDirectory"
	userMail          = "mail"
	userSN            = "sn"
	userCN            = "cn"
)

// ParseUser parses an ldap entry as a User
func (m *LDAPManager) ParseUser(entry *ldap.Entry) *pb.User {
	uid, _ := strconv.Atoi(entry.GetAttributeValue(userUidNumber))
	gid, _ := strconv.Atoi(entry.GetAttributeValue(userGidNumber))
	return &pb.User{
		Username:      entry.GetAttributeValue(m.AccountAttribute),
		FirstName:     entry.GetAttributeValue(userGivenName),
		LastName:      entry.GetAttributeValue(userSN),
		CN:            entry.GetAttributeValue(userCN),
		DisplayName:   entry.GetAttributeValue(userDisplayName),
		UID:           int32(uid),
		GID:           int64(gid),
		LoginShell:    entry.GetAttributeValue(userLoginShell),
		HomeDirectory: entry.GetAttributeValue(userHomeDirectory),
		Email:         entry.GetAttributeValue(userMail),
	}
}

func (m *LDAPManager) userFields() []string {
	return []string{
		m.AccountAttribute,
		userGivenName,
		userSN,
		userCN,
		userDisplayName,
		userUidNumber,
		userGidNumber,
		userLoginShell,
		userHomeDirectory,
		userMail,
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
		return nil, fmt.Errorf(
			"failed to get user %q: %v",
			username, err,
		)
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleUsersError{
			Username: username,
			Count:    len(result.Entries),
		}
	}
	return m.ParseUser(result.Entries[0]), nil
}
