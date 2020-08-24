package ldapmanager

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap"
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

// NewAccountRequest ...
type NewAccountRequest struct {
	FirstName        string `json:"first_name" form:"first_name"`
	LastName         string `json:"last_name" form:"last_name"`
	Username         string `json:"username" form:"username"`
	Password         string `json:"password" form:"password"`
	Email            string `json:"email" form:"email"`
	HashingAlgorithm ldaphash.LDAPPasswordHashingAlgorithm
}

func validEmail(e string) bool {
	if len(e) < 3 && len(e) > 254 {
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

// Validate ...
func (req *NewAccountRequest) Validate() error {
	var invalid []string
	if req.Username == "" || !validUsername(req.Username) {
		invalid = append(invalid, "username")
	}
	if req.Password == "" || !validPassword(req.Password) {
		invalid = append(invalid, "password")
	}
	if req.Email == "" || !validEmail(req.Email) {
		invalid = append(invalid, "email")
	}
	if req.FirstName == "" {
		invalid = append(invalid, "first name")
	}
	if req.LastName == "" {
		invalid = append(invalid, "last name")
	}
	if len(invalid) > 0 {
		return &AccountValidationError{Invalid: invalid}
	}
	return nil
}

func (m *LDAPManager) defaultUserFields() []string {
	return []string{m.AccountAttribute, "givenname", "sn", "mail"}
}

func parseUser(entry *ldap.Entry) map[string]string {
	user := make(map[string]string)
	for _, attr := range entry.Attributes {
		user[attr.Name] = entry.GetAttributeValue(attr.Name)
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
	if err := m.NewGroup(&NewGroupRequest{Name: m.DefaultUserGroup, Members: []string{username}}); err != nil {
		// Fall back to create a new group group for the user
		if err := m.NewGroup(&NewGroupRequest{Name: username, Members: []string{username}}); err != nil {
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

// GetUserListRequest ...
type GetUserListRequest struct {
	ListOptions
	Filters string
	Fields  []string
}

// GetUserList ...
func (m *LDAPManager) GetUserList(req *GetUserListRequest) ([]map[string]string, error) {
	if len(req.Fields) < 1 {
		req.Fields = m.defaultUserFields()
	}
	if req.SortKey == "" {
		req.SortKey = m.AccountAttribute
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
	users := make(map[string]map[string]string)
	for _, entry := range result.Entries {
		if entryKey := entry.GetAttributeValue(req.SortKey); entryKey != "" {
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
		if req.SortOrder == SortDescending {
			return !asc
		}
		return asc
	})
	// Clip
	clippedKeys := keys
	var clippedUsers []map[string]string
	if req.Start >= 0 && req.End < len(keys) && req.Start < req.End {
		clippedKeys = keys[req.Start:req.End]
	}
	for _, key := range clippedKeys {
		clippedUsers = append(clippedUsers, users[key])
	}
	return clippedUsers, nil
}

// AuthenticateUser ...
func (m *LDAPManager) AuthenticateUser(username string, password string) (string, error) {
	// Validate
	if username == "" || password == "" {
		return "", errors.New("must provide username and password")
	}
	// Search for the DN for the given username. If found, try binding with the DN and user's password.
	// If the binding succeeds, return the DN.
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(username)),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return "", err
	}
	if len(result.Entries) != 1 {
		return "", &ZeroOrMultipleAccountsError{Username: username, Count: len(result.Entries)}
	}
	// Make sure to always re-bind as admin afterwards
	defer m.BindAdmin()
	userDN := result.Entries[0].DN
	if err := m.ldap.Bind(userDN, password); err != nil {
		return "", fmt.Errorf("unable to bind as %q", username)
	}
	reg, err := regexp.Compile(fmt.Sprintf("%s=(.*?),", m.AccountAttribute))
	if err != nil {
		return "", errors.New("failed to compile regex")
	}
	matchedDN := reg.FindString(userDN)
	return matchedDN, nil
}

// GetAccount ...
func (m *LDAPManager) GetAccount(username string) (map[string]string, error) {
	if username == "" {
		return nil, errors.New("account username must not be empty")
	}
	// Check for existing user with the same username
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(username)),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, fmt.Errorf("failed to get account %q: %v", username, err)
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleAccountsError{Username: username, Count: len(result.Entries)}
	}
	return parseUser(result.Entries[0]), nil
}

// NewAccount ...
func (m *LDAPManager) NewAccount(req *NewAccountRequest) error {
	// Validate
	if err := req.Validate(); err != nil {
		return err
	}
	// Check for existing user with the same username
	req.Username = escapeDN(req.Username)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", m.AccountAttribute, req.Username, m.UserGroupDN),
		[]string{},
		[]ldap.Control{},
	))
	if err != nil {
		return fmt.Errorf("failed to check for existing user %q: %v", req.Username, err)
	}
	if len(result.Entries) > 0 {
		return fmt.Errorf("account with username %q already exists", req.Username)
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
	userDN := m.AccountNamed(req.Username)
	group, GID, err := m.getNewAccountGroup(req.Username)
	if err != nil {
		return err
	}

	if req.HashingAlgorithm == ldaphash.DEFAULT {
		req.HashingAlgorithm = m.HashingAlgorithm
	}

	hashedPassword, err := ldaphash.Password(req.Password, req.HashingAlgorithm)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	fullName := fmt.Sprintf("%s %s", req.FirstName, req.LastName)
	userAttributes := []ldap.Attribute{
		{Type: "objectClass", Vals: []string{"person", "inetOrgPerson", "posixAccount"}},
		{Type: "uid", Vals: []string{req.Username}},
		{Type: "givenName", Vals: []string{req.FirstName}},
		{Type: "sn", Vals: []string{req.LastName}},
		{Type: "cn", Vals: []string{fullName}},
		{Type: "displayName", Vals: []string{fullName}},
		{Type: "uidNumber", Vals: []string{strconv.Itoa(newUID)}},
		{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		{Type: "loginShell", Vals: []string{m.DefaultUserShell}},
		{Type: "homeDirectory", Vals: []string{fmt.Sprintf("/home/%s", req.Username)}},
		{Type: "userPassword", Vals: []string{hashedPassword}},
		{Type: "mail", Vals: []string{req.Email}},
	}

	addUserRequest := &ldap.AddRequest{
		DN:         userDN,
		Attributes: userAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(addUserRequest)
	if err := m.ldap.Add(addUserRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			return &AccountAlreadyExistsError{Username: req.Username}
		}
		return fmt.Errorf("failed to add user %q: %v", userDN, err)
	}
	if err := m.AddGroupMember(group, req.Username); err != nil && !ldap.IsErrorWithCode(err, ldap.LDAPResultAttributeOrValueExists) {
		return fmt.Errorf("failed to add user %q to group %q: %v", req.Username, group, err)
	}
	if err := m.updateLastID("lastUID", newUID); err != nil {
		return err
	}
	log.Infof("added new account %q (member of group %q)", req.Username, group)
	return nil
}

// DeleteAccount ...
func (m *LDAPManager) DeleteAccount(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("%s=%s,%s", m.AccountAttribute, escapeDN(username), m.UserGroupDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	log.Infof("removed account %q", username)
	return nil
}
