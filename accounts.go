package ldapmanager

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"google.golang.org/grpc/codes"

	"github.com/go-ldap/ldap"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	ldaphash "github.com/romnnn/ldap-manager/hash"
	log "github.com/sirupsen/logrus"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// AccountAlreadyExistsError ...
type AccountAlreadyExistsError struct {
	ApplicationError
	Username string
}

// Error ...
func (e *AccountAlreadyExistsError) Error() string {
	return fmt.Sprintf("account with username %q already exists", e.Username)
}

// Code ...
func (e *AccountAlreadyExistsError) Code() codes.Code {
	return codes.AlreadyExists
}

// ZeroOrMultipleAccountsError ...
type ZeroOrMultipleAccountsError struct {
	ApplicationError
	Username string
	Count    int
}

// Error ...
func (e *ZeroOrMultipleAccountsError) Error() string {
	if e.Count > 1 {
		return fmt.Sprintf("multiple (%d) accounts with username %q", e.Count, e.Username)
	}
	return fmt.Sprintf("no account with username %q", e.Username)
}

// Code ...
func (e *ZeroOrMultipleAccountsError) Code() codes.Code {
	if e.Count > 1 {
		return codes.Internal
	}
	return codes.NotFound
}

// AccountValidationError ...
type AccountValidationError struct {
	ApplicationError
	Invalid []string
}

// Error ...
func (e *AccountValidationError) Error() string {
	return fmt.Sprintf("invalid account request. missing or invalid: %v", e.Invalid)
}

// Code ...
func (e *AccountValidationError) Code() codes.Code {
	return codes.InvalidArgument
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
	return []string{
		m.AccountAttribute,
		"givenName", "sn", "cn", "displayName", "uidNumber", "gidNumber", "loginShell", "homeDirectory", "mail"}
}

func parseUser(entry *ldap.Entry) *pb.User {
	user := &pb.User{Data: make(map[string]string)}
	for _, attr := range entry.Attributes {
		user.Data[attr.Name] = entry.GetAttributeValue(attr.Name)
	}
	return user
}

func (m *LDAPManager) getGroupForAccount(username string) (string, int, error) {
	// First, try to get the default user group
	group := m.DefaultUserGroup
	if defaultGID, err := m.getGroupGID(m.DefaultUserGroup); err == nil {
		return group, defaultGID, nil
	}
	// The default user group might not yet exist
	// Note that a group can only be created with at least one member when using RFC2307BIS
	// Because we need the GID to create the user, strict checking of members remains disabled because they are added after the group
	strict := false
	if err := m.NewGroup(&pb.NewGroupRequest{Name: m.DefaultUserGroup, Members: []string{username}}, strict); err != nil {
		// Fall back to create a new group for the user
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

func (m *LDAPManager) countAccounts() (int, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=*)", m.AccountAttribute),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return 0, err
	}
	return len(result.Entries), nil
}

// GetUserList ...
func (m *LDAPManager) GetUserList(req *pb.GetUserListRequest) (*pb.UserList, error) {
	if req.GetSortKey() == "" {
		req.SortKey = m.AccountAttribute
	}
	filter := parseFilter(req.Filter)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(%s=*)%s)", m.AccountAttribute, filter),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	total, err := m.countAccounts()
	if err != nil {
		return nil, err
	}
	users := make(map[string]*pb.User)
	for _, entry := range result.Entries {
		if entryKey := entry.GetAttributeValue(req.GetSortKey()); entryKey != "" {
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
		if req.GetSortOrder() == pb.SortOrder_ASCENDING {
			return !asc
		}
		return asc
	})
	// Clip
	clippedKeys := keys
	if req.GetStart() >= 0 && req.GetEnd() < int32(len(keys)) && req.GetStart() < req.GetEnd() {
		clippedKeys = keys[req.GetStart():req.GetEnd()]
	}
	clipped := &pb.UserList{Total: int64(total)}
	for _, key := range clippedKeys {
		clipped.Users = append(clipped.Users, users[key])
	}
	return clipped, nil
}

// AuthenticateUser ...
func (m *LDAPManager) AuthenticateUser(req *pb.LoginRequest) (*ldap.Entry, error) {
	// Validate
	if req.GetUsername() == "" || req.GetPassword() == "" {
		return nil, &ValidationError{Message: "must provide username and password"}
	}
	// Search for the DN for the given username. If found, try binding with the DN and user's password.
	// If the binding succeeds, return the DN.
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(req.GetUsername())),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	}
	// Make sure to always re-bind as admin afterwards
	defer m.BindAdmin()
	userDN := result.Entries[0].DN
	if err := m.ldap.Bind(userDN, req.GetPassword()); err != nil {
		return nil, fmt.Errorf("unable to bind as %q", req.GetUsername())
	}
	return result.Entries[0], nil
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
func ValidAccountRequest(acc *pb.Account) error {
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
		return &AccountValidationError{Invalid: invalid}
	}
	return nil
}

// NewAccount ...
func (m *LDAPManager) NewAccount(req *pb.NewAccountRequest, algorithm pb.HashingAlgorithm) error {
	// Validate
	account := req.GetAccount()
	if err := ValidAccountRequest(account); err != nil {
		return err
	}
	// Check for existing user with the same username
	account.Username = escapeDN(account.GetUsername())
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, account.GetUsername()),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return fmt.Errorf("failed to check for existing user %q: %v", account.GetUsername(), err)
	}
	if len(result.Entries) > 0 {
		return fmt.Errorf("account with username %q already exists", account.GetUsername())
	}

	loginShell := account.GetLoginShell()
	if loginShell == "" {
		loginShell = m.DefaultUserShell
	}

	homeDirectory := account.GetHomeDirectory()
	if homeDirectory == "" {
		homeDirectory = fmt.Sprintf("/home/%s", account.GetUsername())
	}

	newUID := int(account.GetUid())
	if newUID < MinUID {
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
		newUID = highestUID + 1
	}

	var group string
	GID := int(account.GetGid())
	if GID < MinGID {
		group, GID, err = m.getGroupForAccount(account.GetUsername())
	} else {
		group, GID, err = m.getGroupByGID(GID)
	}
	if err != nil {
		return err
	}

	if algorithm == pb.HashingAlgorithm_DEFAULT {
		algorithm = m.HashingAlgorithm
	}

	hashedPassword, err := ldaphash.Password(account.GetPassword(), algorithm)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}
	fullName := fmt.Sprintf("%s %s", account.GetFirstName(), account.GetLastName())
	userAttributes := []ldap.Attribute{
		{Type: "objectClass", Vals: []string{"person", "inetOrgPerson", "posixAccount"}},
		{Type: m.AccountAttribute, Vals: []string{account.GetUsername()}},
		{Type: "givenName", Vals: []string{account.GetFirstName()}},
		{Type: "sn", Vals: []string{account.GetLastName()}},
		{Type: "cn", Vals: []string{fullName}},
		{Type: "displayName", Vals: []string{fullName}},
		{Type: "uidNumber", Vals: []string{strconv.Itoa(newUID)}},
		{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		{Type: "loginShell", Vals: []string{loginShell}},
		{Type: "homeDirectory", Vals: []string{homeDirectory}},
		{Type: "userPassword", Vals: []string{hashedPassword}},
		{Type: "mail", Vals: []string{account.GetEmail()}},
	}

	userDN := m.AccountNamed(account.GetUsername())
	addUserRequest := &ldap.AddRequest{
		DN:         userDN,
		Attributes: userAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debugf("addUserRequest=%v", addUserRequest)
	if err := m.ldap.Add(addUserRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			return &AccountAlreadyExistsError{Username: account.GetUsername()}
		}
		return fmt.Errorf("failed to add user %q: %v", userDN, err)
	}
	allowNonExistent := false
	if err := m.AddGroupMember(&pb.GroupMember{Group: group, Username: account.GetUsername()}, allowNonExistent); err != nil {
		if _, ok := err.(*MemberAlreadyExistsError); !ok {
			return fmt.Errorf("failed to add user %q to group %q: %v", account.GetUsername(), group, err)
		}
	}
	if err := m.updateLastID("lastUID", newUID); err != nil {
		return err
	}
	log.Infof("added new account %q (member of group %q)", account.GetUsername(), group)
	return nil
}

// UpdateAccount ...
func (m *LDAPManager) UpdateAccount(req *pb.UpdateAccountRequest, algorithm pb.HashingAlgorithm, isAdmin bool) (string, int, error) {
	// Check if the user even exists
	req.Username = escapeDN(req.GetUsername())
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s)", m.AccountAttribute, req.GetUsername()),
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return "", 0, fmt.Errorf("failed to check for existing user %q: %v", req.GetUsername(), err)
	}
	if len(result.Entries) != 1 {
		return "", 0, fmt.Errorf("account with username %q does not exist", req.GetUsername())
	}

	user := result.Entries[0]
	uidNumber, _ := strconv.Atoi(user.GetAttributeValue("uidNumber"))
	username := req.GetUsername()
	userDN := user.DN
	update := req.GetUpdate()

	// Check if the username was changed which requires a DN change
	if update.GetUsername() != "" {
		username = update.GetUsername()
		// Make sure the new username is not taken
		result, err := m.ldap.Search(ldap.NewSearchRequest(
			m.UserGroupDN,
			ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
			fmt.Sprintf("(%s=%s)", m.AccountAttribute, username),
			[]string{"dn"},
			[]ldap.Control{},
		))
		if err != nil {
			return "", 0, fmt.Errorf("failed to check for users with same username: %v", err)
		}
		if len(result.Entries) > 0 {
			return "", 0, &AccountAlreadyExistsError{Username: username}
		}

		// Migrate the DN
		modifyRequest := &ldap.ModifyDNRequest{
			DN:           m.AccountNamed(req.GetUsername()),
			NewRDN:       fmt.Sprintf("%s=%s", m.AccountAttribute, username),
			DeleteOldRDN: true,
			NewSuperior:  "",
		}
		log.Debugf("RenameAccount modifyRequest=%v", modifyRequest)
		if err := m.ldap.ModifyDN(modifyRequest); err != nil {
			return "", 0, err
		}
		log.Infof("renamed user from %q to %q", req.GetUsername(), username)
		userDN = m.AccountNamed(username)

		// migrate user from all his groups
		groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{Username: username})
		if err != nil {
			return "", 0, fmt.Errorf("failed to get list of groups: %v", err)
		}
		for _, group := range groups.GetGroups() {
			allowDeleteOfDefaultGroups := true
			if err := m.DeleteGroupMember(&pb.GroupMember{Group: group, Username: req.GetUsername()}, allowDeleteOfDefaultGroups); err != nil {
				if _, ok := err.(*NoSuchMemberError); !ok {
					return "", 0, fmt.Errorf("failed to remove renamed user (%q -> %q) from group %q: %v", req.GetUsername(), username, group, err)
				}
			}
			allowNonExistent := true
			if err := m.AddGroupMember(&pb.GroupMember{Group: group, Username: username}, allowNonExistent); err != nil {
				if _, ok := err.(*MemberAlreadyExistsError); !ok {
					return "", 0, fmt.Errorf("failed to add renamed user (%q -> %q) to group %q: %v", req.GetUsername(), username, group, err)
				}
			}
			log.Infof("Switched member %q to %q in user group %q ", req.GetUsername(), username, group)
		}
	}

	modifyAccountRequest := ldap.NewModifyRequest(
		userDN,
		[]ldap.Control{},
	)
	if firstName := update.GetFirstName(); firstName != "" {
		modifyAccountRequest.Replace("givenName", []string{firstName})
	}
	if lastName := update.GetLastName(); lastName != "" {
		modifyAccountRequest.Replace("sn", []string{lastName})
	}
	if loginShell := update.GetLoginShell(); loginShell != "" {
		modifyAccountRequest.Replace("loginShell", []string{loginShell})
	}
	if homeDirectory := update.GetHomeDirectory(); homeDirectory != "" {
		modifyAccountRequest.Replace("homeDirectory", []string{homeDirectory})
	}
	if mail := update.GetEmail(); mail != "" {
		modifyAccountRequest.Replace("mail", []string{mail})
	}
	if password := update.GetPassword(); password != "" {
		hashedPassword, err := ldaphash.Password(password, algorithm)
		if err != nil {
			return "", 0, fmt.Errorf("failed to hash password: %v", err)
		}
		modifyAccountRequest.Replace("userPassword", []string{hashedPassword})
	}

	if isAdmin {
		// Only the admin is allowed to change these because they identify a unique user (username + uidNumber)
		if uid := update.GetUid(); uid >= MinUID {
			uidNumber = int(uid)
			modifyAccountRequest.Replace("uidNumber", []string{strconv.Itoa(int(uid))})
		}
		if gid := update.GetGid(); gid >= MinGID {
			modifyAccountRequest.Replace("gidNumber", []string{strconv.Itoa(int(gid))})
		}
	}

	log.Debugf("modifyAccountRequest=%v", modifyAccountRequest)
	if err := m.ldap.Modify(modifyAccountRequest); err != nil {
		return "", 0, fmt.Errorf("failed to modify existing user: %v", err)
	}
	log.Infof("updated %d attributes of user %q", len(modifyAccountRequest.Changes), username)
	return username, uidNumber, nil
}

// DeleteAccount ...
func (m *LDAPManager) DeleteAccount(req *pb.DeleteAccountRequest, keepGroups bool) error {
	if req.GetUsername() == "" {
		return errors.New("username must not be empty")
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("%s=%s,%s", m.AccountAttribute, escapeDN(req.GetUsername()), m.UserGroupDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	if !keepGroups {
		// delete the account from all its groups
		groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{Username: req.GetUsername()})
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
