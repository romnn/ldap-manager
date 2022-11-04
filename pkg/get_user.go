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

// StatusError returns the GRPC status error for this error
func (e *ZeroOrMultipleUsersError) StatusError() error {
	if e.Count > 1 {
		return status.Errorf(codes.Internal, e.Error())
	}
	return status.Errorf(codes.NotFound, e.Error())
}

const (
	userUIDNumber     = "uidNumber"
	userGIDNumber     = "gidNumber"
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
	UID, _ := strconv.Atoi(entry.GetAttributeValue(userUIDNumber))
	GID, _ := strconv.Atoi(entry.GetAttributeValue(userGIDNumber))
	return &pb.User{
		Username:      entry.GetAttributeValue(m.AccountAttribute),
		FirstName:     entry.GetAttributeValue(userGivenName),
		LastName:      entry.GetAttributeValue(userSN),
		CN:            entry.GetAttributeValue(userCN),
		DisplayName:   entry.GetAttributeValue(userDisplayName),
		UID:           int32(UID),
		GID:           int64(GID),
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
		userUIDNumber,
		userGIDNumber,
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
	conn, err := m.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	result, err := conn.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, EscapeFilter(username)),
		m.userFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleUsersError{
			Username: username,
			Count:    len(result.Entries),
		}
	}
	return m.ParseUser(result.Entries[0]), nil
}
