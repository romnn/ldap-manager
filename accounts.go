package ldapmanager

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	ldaphash "github.com/romnnn/ldap-manager/hash"
	log "github.com/sirupsen/logrus"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// AccountAlreadyExistsError ...
type AccountAlreadyExistsError struct {
	Username string
}

// Error ...
func (e *AccountAlreadyExistsError) Error() string {
	return fmt.Sprintf("account with username %q already exists", e.Username)
}

// ZeroOrMultipleAccountsError ...
type ZeroOrMultipleAccountsError struct {
	Username string
	Count    int
}

// Status ...
func (e *ZeroOrMultipleAccountsError) Status() int {
	if e.Count > 1 {
		return http.StatusConflict
	}
	return http.StatusNotFound
}

// Error ...
func (e *ZeroOrMultipleAccountsError) Error() string {
	if e.Count > 1 {
		return fmt.Sprintf("multiple (%d) accounts with username %q", e.Count, e.Username)
	}
	return fmt.Sprintf("no account with username %q", e.Username)
}

// AccountValidationError ...
type AccountValidationError struct {
	Invalid []string
}

// Error ...
func (e *AccountValidationError) Error() string {
	return fmt.Sprintf("invalid account request. missing or invalid: %v", e.Invalid)
}

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

func (m *LDAPManager) defaultUserFields() []string {
	return []string{m.AccountAttribute, "givenname", "sn", "mail"}
}

func parseUser(entry *ldap.Entry) *pb.User {
	user := &pb.User{Data: make(map[string]string)}
	for _, attr := range entry.Attributes {
		user.Data[attr.Name] = entry.GetAttributeValue(attr.Name)
	}
	return user
}

func (m *LDAPManager) getNewAccountGroup(username string) (string, int, error) {
	group := m.DefaultUserGroup
	if defaultGID, err := m.getGroupGID(m.DefaultUserGroup); err == nil {
		return group, defaultGID, nil
	}
	// The default user group might not yet exist
	// Note that a group can only be created with at least one member when using RFC2307BIS
	// Because we need the GID to create the user, strict checking of members remains disabled because they are added after the group
	strict := false
	if err := m.NewGroup(&pb.NewGroupRequest{Name: m.DefaultUserGroup, Members: []string{username}}, strict); err != nil {
		// Fall back to create a new group group for the user
		if err := m.NewGroup(&pb.NewGroupRequest{Name: username, Members: []string{username}}, strict); err != nil {
			if _, ok := err.(*GroupAlreadyExistsError); !ok {
				return group, 0, fmt.Errorf("failed to create group for user %q: %v", username, err)
			}
		}
		group = username
	}

	userGroupGID, err := m.getGroupGID(group)
	if err != nil {
		return group, 0, fmt.Errorf("failed to get GID for group %q: %v", group, err)
	}
	return group, userGroupGID, nil
}

// AccountNamed ...
func (m *LDAPManager) AccountNamed(name string) string {
	return fmt.Sprintf("%s=%s,%s", m.AccountAttribute, escapeDN(name), m.UserGroupDN)
}

// GetUserList ...
func (m *LDAPManager) GetUserList(req *pb.GetUserListRequest) (*pb.UserList, error) {
	if len(req.Fields) < 1 {
		req.Fields = m.defaultUserFields()
	}
	if req.GetOptions() == nil {
		req.Options = &pb.ListOptions{}
	}
	options := req.GetOptions()
	if options.GetSortKey() == "" {
		options.SortKey = m.AccountAttribute
	}
	filter := fmt.Sprintf("(&(%s=*)%s)", m.AccountAttribute, req.Filters)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		req.Fields,
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	users := make(map[string]*pb.User)
	for _, entry := range result.Entries {
		if entryKey := entry.GetAttributeValue(options.GetSortKey()); entryKey != "" {
			users[entryKey] = parseUser(entry)
		}
	}
	// Sort for deterministic clipping
	keys := make([]string, 0, len(users))
	for k := range users {
		keys = append(keys, k)
	}
	// Sort
	sort.Slice(keys, func(i, j int) bool {
		asc := keys[i] < keys[j]
		if options.GetSortOrder() == pb.SortOrder_ASCENDING {
			return !asc
		}
		return asc
	})
	// Clip
	clippedKeys := keys
	clipped := &pb.UserList{}
	if options.GetStart() >= 0 && options.GetEnd() < int32(len(keys)) && options.GetStart() < options.GetEnd() {
		clippedKeys = keys[options.GetStart():options.GetEnd()]
	}
	for _, key := range clippedKeys {
		clipped.Users = append(clipped.Users, users[key])
	}
	return clipped, nil
}

// AuthenticateUser ...
func (m *LDAPManager) AuthenticateUser(req *pb.AuthenticateUserRequest) error {
	// Validate
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return errors.New("must provide username and password")
	}
	// Search for the DN for the given username. If found, try binding with the DN and user's password.
	// If the binding succeeds, return the DN.
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(req.GetUsername())),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return err
	}
	if len(result.Entries) != 1 {
		return &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	}
	// Make sure to always re-bind as admin afterwards
	defer m.BindAdmin()
	userDN := result.Entries[0].DN
	if err := m.ldap.Bind(userDN, req.GetPassword()); err != nil {
		return fmt.Errorf("unable to bind as %q", req.GetUsername())
	}
	return nil
}

// GetAccount ...
func (m *LDAPManager) GetAccount(req *pb.GetAccountRequest) (*pb.User, error) {
	if req.GetUsername() == "" {
		return nil, errors.New("account username must not be empty")
	}
	// Check for existing user with the same username
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(req.GetUsername())),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, fmt.Errorf("failed to get account %q: %v", req.GetUsername(), err)
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	}
	return parseUser(result.Entries[0]), nil
}

// ValidAccountRequest ...
func ValidAccountRequest(req *pb.NewAccountRequest) error {
	var invalid []string
	if req.GetUsername() == "" || !validUsername(req.GetUsername()) {
		invalid = append(invalid, "username")
	}
	if req.GetPassword() == "" || !validPassword(req.GetPassword()) {
		invalid = append(invalid, "password")
	}
	if req.GetEmail() == "" || !validEmail(req.GetEmail()) {
		invalid = append(invalid, "email")
	}
	if req.GetFirstName() == "" {
		invalid = append(invalid, "first name")
	}
	if req.GetLastName() == "" {
		invalid = append(invalid, "last name")
	}
	if len(invalid) > 0 {
		return &AccountValidationError{Invalid: invalid}
	}
	return nil
}

// NewAccount ...
func (m *LDAPManager) NewAccount(req *pb.NewAccountRequest, algorithm pb.HashingAlgorithm) error {
	// Validate
	if err := ValidAccountRequest(req); err != nil {
		return err
	}
	// Check for existing user with the same username
	req.Username = escapeDN(req.GetUsername())
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", m.AccountAttribute, req.GetUsername(), m.UserGroupDN),
		[]string{},
		[]ldap.Control{},
	))
	if err != nil {
		return fmt.Errorf("failed to check for existing user %q: %v", req.GetUsername(), err)
	}
	if len(result.Entries) > 0 {
		return fmt.Errorf("account with username %q already exists", req.GetUsername())
	}
	highestUID, err := m.getHighestID(m.AccountAttribute)
	if err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			// Try to recover by running the setup
			_ = m.setupLastUID()
			highestUID, err = m.getHighestID(m.AccountAttribute)
		}
		if err != nil {
			return fmt.Errorf("failed to get highest %s: %v", m.AccountAttribute, err)
		}
	}
	newUID := highestUID + 1
	userDN := m.AccountNamed(req.GetUsername())
	group, GID, err := m.getNewAccountGroup(req.GetUsername())
	if err != nil {
		return err
	}

	if algorithm == pb.HashingAlgorithm_DEFAULT {
		algorithm = m.HashingAlgorithm
	}

	hashedPassword, err := ldaphash.Password(req.GetPassword(), algorithm)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	fullName := fmt.Sprintf("%s %s", req.GetFirstName(), req.GetLastName())
	userAttributes := []ldap.Attribute{
		{Type: "objectClass", Vals: []string{"person", "inetOrgPerson", "posixAccount"}},
		{Type: "uid", Vals: []string{req.GetUsername()}},
		{Type: "givenName", Vals: []string{req.GetFirstName()}},
		{Type: "sn", Vals: []string{req.GetLastName()}},
		{Type: "cn", Vals: []string{fullName}},
		{Type: "displayName", Vals: []string{fullName}},
		{Type: "uidNumber", Vals: []string{strconv.Itoa(newUID)}},
		{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		{Type: "loginShell", Vals: []string{m.DefaultUserShell}},
		{Type: "homeDirectory", Vals: []string{fmt.Sprintf("/home/%s", req.GetUsername())}},
		{Type: "userPassword", Vals: []string{hashedPassword}},
		{Type: "mail", Vals: []string{req.GetEmail()}},
	}

	addUserRequest := &ldap.AddRequest{
		DN:         userDN,
		Attributes: userAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(addUserRequest)
	if err := m.ldap.Add(addUserRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			return &AccountAlreadyExistsError{Username: req.GetUsername()}
		}
		return fmt.Errorf("failed to add user %q: %v", userDN, err)
	}
	allowNonExistent := false
	if err := m.AddGroupMember(&pb.GroupMember{Group: group, Username: req.GetUsername()}, allowNonExistent); err != nil {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultAttributeOrValueExists) {
			return fmt.Errorf("failed to add user %q to group %q: %v", req.GetUsername(), group, err)
		}
	}
	if err := m.updateLastID("lastUID", newUID); err != nil {
		return err
	}
	log.Infof("added new account %q (member of group %q)", req.GetUsername(), group)
	return nil
}

// DeleteAccount ...
func (m *LDAPManager) DeleteAccount(req *pb.DeleteAccountRequest, leaveGroups bool) error {
	if req.GetUsername() == "" {
		return errors.New("username must not be empty")
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("%s=%s,%s", m.AccountAttribute, escapeDN(req.GetUsername()), m.UserGroupDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	if !leaveGroups {
		// delete the account from all its groups
		groups, err := m.GetGroupList(&pb.GetGroupListRequest{})
		if err != nil {
			return fmt.Errorf("failed to get list of groups: %v", err)
		}
		for _, group := range groups.GetGroups() {
			allowDeleteOfDefaultGroups := true
			if err := m.DeleteGroupMember(&pb.GroupMember{Group: group, Username: req.GetUsername()}, allowDeleteOfDefaultGroups); err != nil {
				if _, ok := err.(*NoSuchMemberError); !ok {
					return fmt.Errorf("failed to remove deleted user %q from group %q: %v", req.GetUsername(), group, err)
				}
			}
		}
	}
	log.Infof("removed account %q", req.GetUsername())
	return nil
}
