package pkg

import (
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

// InvalidUserError ...
type InvalidUserError struct {
	error
	Invalid []string
}

func (e *InvalidUserError) Error() string {
	return fmt.Sprintf("invalid account request. missing or invalid: %v", e.Invalid)
}

func (e *InvalidUserError) StatusError() error {
	return status.Errorf(codes.InvalidArgument, e.Error())
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func validEmail(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}

func validPassword(pw string) bool {
	// TODO: maybe we enforce some password length in the future
	return true
}

func validUsername(un string) bool {
	// TODO: maybe we enforce some username regex in the future
	return true
}

// ValidateUser validates a user account
func ValidateUser(acc *pb.Account) *InvalidUserError {
	var invalid []string
	if acc.GetUsername() == "" || !validUsername(acc.GetUsername()) {
		invalid = append(invalid, "username")
	}
	if acc.GetPassword() == "" || !validPassword(acc.GetPassword()) {
		invalid = append(invalid, "password")
	}
	if acc.GetEmail() == "" || !validEmail(acc.GetEmail()) {
		invalid = append(invalid, "email")
	}
	if acc.GetFirstName() == "" {
		invalid = append(invalid, "first name")
	}
	if acc.GetLastName() == "" {
		invalid = append(invalid, "last name")
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
	// note that a group can only be created with at least one member when using RFC2307BIS
	// because we need the GID to create the user, strict checking of members is disabled
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

// NewUser ...
func (m *LDAPManager) NewUser(req *pb.NewUserRequest, algorithm pb.HashingAlgorithm) error {
	account := req.GetAccount()
	if err := ValidateUser(account); err != nil {
		return err
	}

	// check for existing user with the same username
	username := EscapeDN(account.GetUsername())
	search := ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, username),
		[]string{"dn"},
		[]ldap.Control{},
	)
	result, err := m.ldap.Search(search)
	// log.Infof("result: %v error: %v", result, err)
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
		return &UserAlreadyExistsError{Username: username}
	}

	// set default values
	loginShell := account.GetLoginShell()
	if loginShell == "" {
		loginShell = m.DefaultUserShell
	}

	homeDirectory := account.GetHomeDirectory()
	if homeDirectory == "" {
		homeDirectory = fmt.Sprintf("/home/%s", username)
	}

	// skip for now to see what openLDAP does
	// newUID := int(account.GetUid())
	highestUID, err := m.GetHighestUID()
	if err != nil {

		return fmt.Errorf("failed to get highest %s: %v", m.AccountAttribute, err)
	}
  newUID := highestUID + 1
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

	fullName := fmt.Sprintf("%s %s", account.GetFirstName(), account.GetLastName())
	userAttributes := []ldap.Attribute{
		{Type: "objectClass", Vals: []string{"person", "inetOrgPerson", "posixAccount"}},
		{Type: m.AccountAttribute, Vals: []string{username}},
		{Type: "givenName", Vals: []string{account.GetFirstName()}},
		{Type: "sn", Vals: []string{account.GetLastName()}},
		{Type: "cn", Vals: []string{fullName}},
		{Type: "displayName", Vals: []string{fullName}},
		{Type: "uidNumber", Vals: []string{strconv.Itoa(newUID)}},
		// add to users group
		{Type: "gidNumber", Vals: []string{strconv.Itoa(int(userGroup.GetGid()))}},
		// {Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		{Type: "loginShell", Vals: []string{loginShell}},
		{Type: "homeDirectory", Vals: []string{homeDirectory}},
		// {Type: "userPassword", Vals: []string{hashedPassword}},
		{Type: "mail", Vals: []string{account.GetEmail()}},
	}

	// add user
	userDN := m.UserNamed(username)
	addUserRequest := &ldap.AddRequest{
		DN:         userDN,
		Attributes: userAttributes,
		Controls:   []ldap.Control{},
	}
	log.Infof("addUserRequest=%v", addUserRequest)
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
		// OldPassword: "",
		NewPassword: account.GetPassword(),
	}
	log.Infof("passwordModifyRequest=%v", passwordModifyRequest)
	resp, err := m.ldap.PasswordModify(passwordModifyRequest)
	if err != nil {
		// exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		// if exists {
		// 	return &UserAlreadyExistsError{Username: username}
		// }
		return fmt.Errorf("failed to set password of new user %q: %v", userDN, err)
	}
	// GeneratedPassword
	log.Infof("password response=%v", resp)

	allowNonExistent := false
	if err := m.AddGroupMember(&pb.GroupMember{
		Group:    m.DefaultUserGroup,
		Username: username,
	}, allowNonExistent); err != nil {
		if _, exists := err.(*MemberAlreadyExistsError); !exists {
			return fmt.Errorf("failed to add user %q to group %q: %v", username, userGroup.Name, err)
		}
	}
	// if err := m.updateLastID("lastUID", newUID); err != nil {
	// 	return err
	// }
	log.Infof("added new account %q (member of group %q)", username, userGroup.Name)
	return nil
}
