package pkg

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserAreadyExistsError ...
type UserAlreadyExistsError struct {
	error
	Username string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf("user with username %q already exists", e.Username)
}

func (e *UserAlreadyExistsError) StatusError() error {
	return status.Errorf(codes.AlreadyExists, e.Error())
}

// An InvalidUserError is returned when the user contains invalid values
type InvalidUserError struct {
	error
	Invalid map[string]error
}

func (e *InvalidUserError) Error() string {
	return fmt.Sprintf("invalid account request. missing or invalid: %v", e.Invalid)
}

func (e *InvalidUserError) StatusError() error {
	return status.Errorf(codes.InvalidArgument, e.Error())
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// ValidateEmail validates an email
func ValidateEmail(email string) error {
	if len(email) < 3 {
		return errors.New("email must contain at least 3 characters")
	}
	if len(email) > 254 {
		return errors.New("email can contain at most 254 characters")
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("%q is not a valid email", email)
	}
	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password must not be empty")
	}
	if len(password) < 5 {
		return errors.New("password must contain at least 5 characters")
	}
	return nil
}

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}
	if len(username) < 5 {
		return errors.New("username must contain at least 5 characters")
	}
	return nil
}

// ValidateFirstName validates a first name
func ValidateFirstName(name string) error {
	if name == "" {
		return errors.New("first name must not be empty")
	}
	return nil
}

// ValidateLastName validates a last name
func ValidateLastName(name string) error {
	if name == "" {
		return errors.New("last name must not be empty")
	}
	return nil
}

// ValidateNewUser validates a new user request
func ValidateNewUser(req *pb.NewUserRequest) *InvalidUserError {
	invalid := make(map[string]error)
	if err := ValidateUsername(req.GetUsername()); err != nil {
		invalid["username"] = err
	}
	if err := ValidatePassword(req.GetPassword()); err != nil {
		invalid["password"] = err
	}
	if err := ValidateEmail(req.GetEmail()); err != nil {
		invalid["email"] = err
	}
	if err := ValidateFirstName(req.GetFirstName()); err != nil {
		invalid["first name"] = err
	}
	if err := ValidateLastName(req.GetLastName()); err != nil {
		invalid["last name"] = err
	}
	if len(invalid) > 0 {
		return &InvalidUserError{Invalid: invalid}
	}
	return nil
}

func (m *LDAPManager) GetUserGroup(username string) (*pb.Group, error) {
	// fast path: the user group already exists
	groupName := m.DefaultUserGroup
	if group, err := m.GetGroupByName(groupName); err == nil {
		return group, nil
	}
	// slow path: the default user group might not yet exist
	// note that a group can only be created with at least one member
	// when using RFC2307BIS
	// because we need the GID to create the user,
	// strict checking of members is disabled
	strict := false
	if err := m.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		return nil, err
		// Fall back to create a new group for the user
		// if err := m.NewGroup(&pb.NewGroupRequest{Name: username, Members: []string{username}}, strict); err != nil {
		// 	if _, ok := err.(*GroupAlreadyExistsError); !ok {
		// 		return group, 0, fmt.Errorf("failed to create group for user %q: %v", username, err)
		// 	}
		// }
		// group = username
	}

	group, err := m.GetGroupByName(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get GID for group %q: %v", groupName, err)
	}
	return group, nil
}

// NewUser adds a new user
func (m *LDAPManager) NewUser(req *pb.NewUserRequest) error {
	if err := ValidateNewUser(req); err != nil {
		return err
	}

	// check for existing user with the same username
	username := EscapeDN(req.GetUsername())
	search := ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, username),
		[]string{"dn"},
		[]ldap.Control{},
	)
	result, err := m.ldap.Search(search)
	if err != nil {
		return err
		// notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		// if notFound {
		// there might be no users group, in which case this is fine
		// strict := false
		// err := m.NewGroup(&pb.NewGroupRequest{
		// 	Name:    m.DefaultUserGroup,
		// 	Members: []string{account.GetUsername()},
		// }, strict)
		// if err != nil {
		// 	log.Errorf("failed to create user group: %v", err)
		// }
		// already exists? -> fail
		// not found? -> fail
		// created? -> succeed
		// notFound = ldap.IsErrorWithCode(userGroupErr, ldap.LDAPResultNoSuchObject)
		// if !notFound {
		// 	err = nil
		// }
		// if there also is no users group, there must have been a problem with the setup

		// try again
		result, err = m.ldap.Search(search)
		// }
	}

	if err != nil {
		return fmt.Errorf("failed to check for existing user %q: %v", username, err)
	}
	if len(result.Entries) > 0 {
		return &UserAlreadyExistsError{
			Username: username,
		}
	}

	// set default values
	loginShell := req.GetLoginShell()
	if loginShell == "" {
		loginShell = m.DefaultUserShell
	}

	homeDirectory := req.GetHomeDirectory()
	if homeDirectory == "" {
		homeDirectory = fmt.Sprintf("/home/%s", username)
	}

	UID, err := m.GetHighestUID()
	if err != nil {
		return fmt.Errorf("failed to get highest UID (%s): %v", m.AccountAttribute, err)
	}
	// newUID := highestUID + 1
	// if newUID < MinUID {
	// 	highestUID, err := m.GetHighestID(m.AccountAttribute)
	// 	if err != nil {
	// 		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
	// 			// Try to recover by running the setup
	// 			_ = m.setupLastUID()
	// 			highestUID, err = m.GetHighestID(m.AccountAttribute)
	// 		}
	// 		if err != nil {
	// 			return fmt.Errorf("failed to get highest %s: %v", m.AccountAttribute, err)
	// 		}
	// 	}
	// 	newUID = highestUID + 1
	// }

	// GID := 0
	// group, GID, err = m.GetGroupByGID(GID)
	// group, GID, err := m.GetGroupByGID(GID)
	// var group string
	// GID := int(account.GetGid())
	// if GID < MinGID {
	// 	group, GID, err = m.getGroupForUser(account.GetUsername())
	// } else {
	// 	group, GID, err = m.GetGroupByGID(GID)
	// }
	// if err != nil {
	// 	return err
	// }

	// if algorithm == pb.HashingAlgorithm_DEFAULT {
	// 	algorithm = m.HashingAlgorithm
	// }

	// hashedPassword := account.GetPassword()
	// hashedPassword, err := account.GetPassword()
	// hash.Password(account.GetPassword(), algorithm)
	// if err != nil {
	// 	return fmt.Errorf("failed to hash password: %v", err)
	// }

	userGroup, err := m.GetUserGroup(username)
	// userGroup, err := m.GetGroupByName(m.DefaultUserGroup)
	if err != nil {
		return err
		// fmt.Errorf("failed to get GID of default user group: %v", err)
	}

	fullName := fmt.Sprintf("%s %s", req.GetFirstName(), req.GetLastName())
	userAttributes := []ldap.Attribute{
		{Type: "objectClass", Vals: []string{
			"person",
			"inetOrgPerson",
			"posixAccount",
		}},
		{Type: m.AccountAttribute, Vals: []string{username}},
		{Type: "givenName", Vals: []string{req.GetFirstName()}},
		{Type: "sn", Vals: []string{req.GetLastName()}},
		{Type: "cn", Vals: []string{fullName}},
		{Type: "displayName", Vals: []string{fullName}},
		{Type: "uidNumber", Vals: []string{strconv.Itoa(UID)}},
		{Type: "gidNumber", Vals: []string{strconv.Itoa(int(userGroup.GetGID()))}},
		{Type: "loginShell", Vals: []string{loginShell}},
		{Type: "homeDirectory", Vals: []string{homeDirectory}},
		{Type: "mail", Vals: []string{req.GetEmail()}},
	}

	// add user
	userDN := m.UserNamed(username)
	addUserRequest := &ldap.AddRequest{
		DN:         userDN,
		Attributes: userAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(PrettyPrint(addUserRequest))
	if err := m.ldap.Add(addUserRequest); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if exists {
			return &UserAlreadyExistsError{Username: username}
		}
		return fmt.Errorf("failed to add user %q: %v", userDN, err)
	}

	// change password of user
	passwordModifyRequest := &ldap.PasswordModifyRequest{
		UserIdentity: userDN,
		NewPassword:  req.GetPassword(),
	}
	log.Debug(PrettyPrint(passwordModifyRequest))
	resp, err := m.ldap.PasswordModify(passwordModifyRequest)
	if err != nil {
		// exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		// if exists {
		// 	return &UserAlreadyExistsError{Username: username}
		// }
		return fmt.Errorf("failed to set password of new user %q: %v", userDN, err)
	}
	log.Debug(PrettyPrint(resp))

	allowNonExistent := false
	if err := m.AddGroupMember(&pb.GroupMember{
		Group:    m.DefaultUserGroup,
		Username: username,
	}, allowNonExistent); err != nil {
		if _, exists := err.(*MemberAlreadyExistsError); !exists {
			return fmt.Errorf("failed to add user %q to group %q: %v", username, userGroup.Name, err)
		}
	}
	if err := m.updateLastID("lastUID", UID+1); err != nil {
		return err
	}
	log.Infof("added new user %q (member of group %q)", username, userGroup.Name)
	return nil
}
