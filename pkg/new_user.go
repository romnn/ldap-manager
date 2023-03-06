package pkg

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/badoux/checkmail"
	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserAlreadyExistsError is returned when a username already exists
type UserAlreadyExistsError struct {
	ldaperror.ApplicationError
	Username string
}

func (e *UserAlreadyExistsError) Error() string {
	return fmt.Sprintf(
		"user with username %q already exists",
		e.Username,
	)
}

// StatusError returns the GRPC status error for this error
func (e *UserAlreadyExistsError) StatusError() error {
	return status.Errorf(codes.AlreadyExists, e.Error())
}

// InvalidUserError is returned when the user contains invalid values
type InvalidUserError struct {
	ldaperror.ApplicationError
	Invalid map[string]error
}

func (e *InvalidUserError) Error() string {
	return fmt.Sprintf(
		"invalid new user: missing or invalid: %v",
		e.Invalid,
	)
}

// StatusError returns the GRPC status error for this error
func (e *InvalidUserError) StatusError() error {
	message := "Invalid new user request:"
	for field, reason := range e.Invalid {
		message += fmt.Sprintf("\n-> %s: %s", field, reason)
	}
	return status.Errorf(codes.InvalidArgument, message)
}

// ValidateEmail validates an email
func ValidateEmail(email string) error {
	if len(email) < 3 {
		return errors.New("must contain at least 3 characters")
	}
	if len(email) > 254 {
		return errors.New("can contain at most 254 characters")
	}
	if err := checkmail.ValidateFormat(email); err != nil {
		return fmt.Errorf("%q is not a valid email: %v", email, err)
	}
	return nil
}

// ValidatePassword validates a password
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("must not be empty")
	}
	if len(password) < 5 {
		return errors.New("must contain at least 5 characters")
	}
	return nil
}

// ValidateUsername validates a username
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("must not be empty")
	}
	if len(username) < 5 {
		return errors.New("must contain at least 5 characters")
	}
	return nil
}

// ValidateFirstName validates a first name
func ValidateFirstName(name string) error {
	if name == "" {
		return errors.New("must not be empty")
	}
	return nil
}

// ValidateLastName validates a last name
func ValidateLastName(name string) error {
	if name == "" {
		return errors.New("must not be empty")
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

//// GetUserGroup gets or creates the user group
////
//// If there exist no users yet, the default user groups is created with
//// the given user as the initial member
//func (m *LDAPManager) GetUserGroup(username string) (*pb.Group, error) {
//	// fast path: the user group already exists
//	groupName := m.DefaultUserGroup
//	if group, err := m.GetGroupByName(groupName); err == nil {
//		return group, nil
//	}
//	// slow path: the default user group might not yet exist
//	// note that a group can only be created with at least one member
//	// when using RFC2307BIS
//	strict := false
//	if err := m.NewGroup(&pb.NewGroupRequest{
//		Name:    groupName,
//    // note: user must already exist for memberOf overlay to work correctly
//		Members: []string{username},
//	}, strict); err != nil {
//		return nil, err
//	}

//	group, err := m.GetGroupByName(groupName)
//	if err != nil {
//		return nil, fmt.Errorf(
//			"failed to get group %q: %v",
//			groupName, err,
//		)
//	}
//	return group, nil
//}

func (m *LDAPManager) checkUserExists(username string) error {
	_, err := m.GetUser(username)
	if err != nil {
		notFoundErr, notFound := err.(*ZeroOrMultipleUsersError)
		if notFound {
			if notFoundErr.Count > 1 {
				return &UserAlreadyExistsError{
					Username: username,
				}
			}
		} else {
			return fmt.Errorf(
				"failed to check for existing user %q: %v",
				username, err,
			)
		}
	} else {
		return &UserAlreadyExistsError{
			Username: username,
		}
	}
	return nil
}

// NewUser adds a new user
func (m *LDAPManager) NewUser(req *pb.NewUserRequest) error {
	if err := ValidateNewUser(req); err != nil {
		return err
	}

	username := req.GetUsername()
	if err := m.checkUserExists(username); err != nil {
		return err
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
		return fmt.Errorf(
			"failed to get highest UID (%s): %v",
			m.AccountAttribute, err,
		)
	}

	fullName := fmt.Sprintf(
		"%s %s",
		req.GetFirstName(),
		req.GetLastName(),
	)
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
		{Type: "gidNumber", Vals: []string{
			strconv.Itoa(MinGID),
			// strconv.Itoa(int(userGroup.GetGID())),
		}},
		{Type: "loginShell", Vals: []string{loginShell}},
		{Type: "homeDirectory", Vals: []string{homeDirectory}},
		{Type: "mail", Vals: []string{req.GetEmail()}},
	}

	// add user
	addUserRequest := &ldap.AddRequest{
		DN:         m.UserDN(username),
		Attributes: userAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(PrettyPrint(addUserRequest))

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Add(addUserRequest); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if exists {
			return &UserAlreadyExistsError{
				Username: username,
			}
		}
		return fmt.Errorf(
			"failed to add user %q: %v",
			username, err,
		)
	}

	// change password of user
	passwordModifyRequest := &ldap.PasswordModifyRequest{
		UserIdentity: m.UserDN(username),
		NewPassword:  req.GetPassword(),
	}
	log.Debug(PrettyPrint(passwordModifyRequest))
	_, err = conn.PasswordModify(passwordModifyRequest)
	if err != nil {
		return fmt.Errorf(
			"failed to set password of new user %q: %v",
			username, err,
		)
	}

	groupName := m.DefaultUserGroup

	strict := false
	if err := m.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		if _, exists := err.(*GroupAlreadyExistsError); !exists {
			return fmt.Errorf(
				"failed to create %q group: %v",
				groupName, err,
			)
		}
	}

	allowNonExistent := false
	if err := m.AddGroupMember(&pb.GroupMember{
		Group:    groupName,
		Username: username,
	}, allowNonExistent); err != nil {
		if _, exists := err.(*MemberAlreadyExistsError); !exists {
			return fmt.Errorf(
				"failed to add user %q to group %q: %v",
				username, groupName, err,
			)
		}
	}
	if err := m.updateLastID("lastUID", UID+1); err != nil {
		return err
	}
	log.Infof(
		"added new user %q (member of group %q)",
		username, groupName,
	)
	return nil
}
