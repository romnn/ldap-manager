package pkg

import (
	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// GetUserList ...
func (m *LDAPManager) GetUserList(req *pb.GetUserListRequest) (*pb.UserList, error) {
	// if req.GetSortKey() == "" {
	// 	req.SortKey = m.AccountAttribute
	// }
	// filter := parseFilter(req.Filter)
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.UserGroupDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(&(%s=*)%s)", m.AccountAttribute, filter),
	// 	m.defaultUserFields(),
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	return nil, err
	// }
	// total, err := m.countAccounts()
	// if err != nil {
	// 	return nil, err
	// }
	// users := make(map[string]*pb.User)
	// for _, entry := range result.Entries {
	// 	if entryKey := entry.GetAttributeValue(req.GetSortKey()); entryKey != "" {
	// 		users[entryKey] = parseUser(entry)
	// 	}
	// }
	// // Sort for deterministic clipping
	// keys := make([]string, 0, len(users))
	// for k := range users {
	// 	keys = append(keys, k)
	// }
	// // Sort
	// sort.Slice(keys, func(i, j int) bool {
	// 	asc := keys[i] < keys[j]
	// 	if req.GetSortOrder() == pb.SortOrder_ASCENDING {
	// 		return !asc
	// 	}
	// 	return asc
	// })
	// // Clip
	// clippedKeys := keys
	// if req.GetStart() >= 0 && req.GetEnd() < int32(len(keys)) && req.GetStart() < req.GetEnd() {
	// 	clippedKeys = keys[req.GetStart():req.GetEnd()]
	// }
	total := 0
	clipped := &pb.UserList{Total: int64(total)}
	// for _, key := range clippedKeys {
	// 	clipped.Users = append(clipped.Users, users[key])
	// }
	return clipped, nil
}

// AuthenticateUser ...
func (m *LDAPManager) AuthenticateUser(req *pb.LoginRequest) (*ldap.Entry, error) {
	// // Validate
	// if req.GetUsername() == "" || req.GetPassword() == "" {
	// 	return nil, &ValidationError{Message: "must provide username and password"}
	// }
	// // Search for the DN for the given username. If found, try binding with the DN and user's password.
	// // If the binding succeeds, return the DN.
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.BaseDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(req.GetUsername())),
	// 	m.defaultUserFields(),
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	return nil, err
	// }
	// if len(result.Entries) != 1 {
	// 	return nil, &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	// }
	// // Make sure to always re-bind as admin afterwards
	// defer m.BindAdmin()
	// userDN := result.Entries[0].DN
	// if err := m.ldap.Bind(userDN, req.GetPassword()); err != nil {
	// 	return nil, fmt.Errorf("unable to bind as %q", req.GetUsername())
	// }
	// return result.Entries[0], nil
	return nil, nil
}

// GetAccount ...
func (m *LDAPManager) GetUser(req *pb.GetUserRequest) (*pb.UserData, error) {
	// if req.GetUsername() == "" {
	// 	return nil, errors.New("account username must not be empty")
	// }
	// // Check for existing user with the same username
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.UserGroupDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", m.AccountAttribute, escapeFilter(req.GetUsername())),
	// 	m.defaultUserFields(),
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get account %q: %v", req.GetUsername(), err)
	// }
	// if len(result.Entries) != 1 {
	// 	return nil, &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	// }
	// return parseUser(result.Entries[0]), nil
	return nil, nil
}

// NewAccount ...
func (m *LDAPManager) NewUser(req *pb.NewUserRequest, algorithm pb.HashingAlgorithm) error {
	// // Validate
	// account := req.GetAccount()
	// if err := ValidAccountRequest(account); err != nil {
	// 	return err
	// }
	// // Check for existing user with the same username
	// account.Username = escapeDN(account.GetUsername())
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.UserGroupDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", m.AccountAttribute, account.GetUsername()),
	// 	[]string{"dn"},
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
	// 		// there might be no users group, in which case this is fine
	// 		strict := false
	// 		noUserGroupErr := m.NewGroup(&pb.NewGroupRequest{Name: m.DefaultUserGroup, Members: []string{req.GetAccount().GetUsername()}}, strict)
	// 		if !ldap.IsErrorWithCode(noUserGroupErr, ldap.LDAPResultNoSuchObject) {
	// 			err = nil
	// 		}
	// 		// if there is also no users group, there must have been a problem with the setup
	// 	}
	// 	if err != nil {
	// 		return fmt.Errorf("failed to check for existing user %q: %v", account.GetUsername(), err)
	// 	}
	// } else {
	// 	if len(result.Entries) > 0 {
	// 		return &AccountAlreadyExistsError{Username: account.GetUsername()}
	// 	}
	// }

	// loginShell := account.GetLoginShell()
	// if loginShell == "" {
	// 	loginShell = m.DefaultUserShell
	// }

	// homeDirectory := account.GetHomeDirectory()
	// if homeDirectory == "" {
	// 	homeDirectory = fmt.Sprintf("/home/%s", account.GetUsername())
	// }

	// newUID := int(account.GetUid())
	// if newUID < MinUID {
	// 	highestUID, err := m.getHighestID(m.AccountAttribute)
	// 	if err != nil {
	// 		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
	// 			// Try to recover by running the setup
	// 			_ = m.setupLastUID()
	// 			highestUID, err = m.getHighestID(m.AccountAttribute)
	// 		}
	// 		if err != nil {
	// 			return fmt.Errorf("failed to get highest %s: %v", m.AccountAttribute, err)
	// 		}
	// 	}
	// 	newUID = highestUID + 1
	// }

	// var group string
	// GID := int(account.GetGid())
	// if GID < MinGID {
	// 	group, GID, err = m.getGroupForAccount(account.GetUsername())
	// } else {
	// 	group, GID, err = m.getGroupByGID(GID)
	// }
	// if err != nil {
	// 	return err
	// }

	// if algorithm == pb.HashingAlgorithm_DEFAULT {
	// 	algorithm = m.HashingAlgorithm
	// }

	// hashedPassword, err := hash.Password(account.GetPassword(), algorithm)
	// if err != nil {
	// 	return fmt.Errorf("failed to hash password: %v", err)
	// }
	// fullName := fmt.Sprintf("%s %s", account.GetFirstName(), account.GetLastName())
	// userAttributes := []ldap.Attribute{
	// 	{Type: "objectClass", Vals: []string{"person", "inetOrgPerson", "posixAccount"}},
	// 	{Type: m.AccountAttribute, Vals: []string{account.GetUsername()}},
	// 	{Type: "givenName", Vals: []string{account.GetFirstName()}},
	// 	{Type: "sn", Vals: []string{account.GetLastName()}},
	// 	{Type: "cn", Vals: []string{fullName}},
	// 	{Type: "displayName", Vals: []string{fullName}},
	// 	{Type: "uidNumber", Vals: []string{strconv.Itoa(newUID)}},
	// 	{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
	// 	{Type: "loginShell", Vals: []string{loginShell}},
	// 	{Type: "homeDirectory", Vals: []string{homeDirectory}},
	// 	{Type: "userPassword", Vals: []string{hashedPassword}},
	// 	{Type: "mail", Vals: []string{account.GetEmail()}},
	// }

	// userDN := m.AccountNamed(account.GetUsername())
	// addUserRequest := &ldap.AddRequest{
	// 	DN:         userDN,
	// 	Attributes: userAttributes,
	// 	Controls:   []ldap.Control{},
	// }
	// log.Debugf("addUserRequest=%v", addUserRequest)
	// if err := m.ldap.Add(addUserRequest); err != nil {
	// 	if ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
	// 		return &AccountAlreadyExistsError{Username: account.GetUsername()}
	// 	}
	// 	return fmt.Errorf("failed to add user %q: %v", userDN, err)
	// }
	// allowNonExistent := false
	// if err := m.AddGroupMember(&pb.GroupMember{Group: group, Username: account.GetUsername()}, allowNonExistent); err != nil {
	// 	if _, ok := err.(*MemberAlreadyExistsError); !ok {
	// 		return fmt.Errorf("failed to add user %q to group %q: %v", account.GetUsername(), group, err)
	// 	}
	// }
	// if err := m.updateLastID("lastUID", newUID); err != nil {
	// 	return err
	// }
	// log.Infof("added new account %q (member of group %q)", account.GetUsername(), group)
	return nil
}

// UpdateAccount ...
func (m *LDAPManager) UpdateUser(req *pb.UpdateUserRequest, algorithm pb.HashingAlgorithm, isAdmin bool) (string, int, error) {
	// // Check if the user even exists
	// req.Username = escapeDN(req.GetUsername())
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.UserGroupDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	fmt.Sprintf("(%s=%s)", m.AccountAttribute, req.GetUsername()),
	// 	m.defaultUserFields(),
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	return "", 0, fmt.Errorf("failed to check for existing user %q: %v", req.GetUsername(), err)
	// }
	// if len(result.Entries) != 1 {
	// 	return "", 0, &ZeroOrMultipleAccountsError{Username: req.GetUsername(), Count: len(result.Entries)}
	// }

	// user := result.Entries[0]
	// uidNumber, _ := strconv.Atoi(user.GetAttributeValue("uidNumber"))
	// username := req.GetUsername()
	// userDN := user.DN
	// update := req.GetUpdate()

	// // Check if the username was changed which requires a DN change
	// if update.GetUsername() != "" && update.GetUsername() != username {
	// 	username = update.GetUsername()
	// 	// Make sure the new username is not taken
	// 	result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 		m.UserGroupDN,
	// 		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 		fmt.Sprintf("(%s=%s)", m.AccountAttribute, username),
	// 		[]string{"dn"},
	// 		[]ldap.Control{},
	// 	))
	// 	if err != nil {
	// 		return "", 0, fmt.Errorf("failed to check for users with same username: %v", err)
	// 	}
	// 	if len(result.Entries) > 0 {
	// 		return "", 0, &AccountAlreadyExistsError{Username: username}
	// 	}

	// 	// Migrate the DN
	// 	modifyRequest := &ldap.ModifyDNRequest{
	// 		DN:           m.AccountNamed(req.GetUsername()),
	// 		NewRDN:       fmt.Sprintf("%s=%s", m.AccountAttribute, username),
	// 		DeleteOldRDN: true,
	// 		NewSuperior:  "",
	// 	}
	// 	log.Debugf("RenameAccount modifyRequest=%v", modifyRequest)
	// 	if err := m.ldap.ModifyDN(modifyRequest); err != nil {
	// 		return "", 0, err
	// 	}
	// 	log.Infof("renamed user from %q to %q", req.GetUsername(), username)
	// 	userDN = m.AccountNamed(username)

	// 	// migrate user from all his groups
	// 	groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{Username: username})
	// 	if err != nil {
	// 		return "", 0, fmt.Errorf("failed to get list of groups: %v", err)
	// 	}
	// 	for _, group := range groups.GetGroups() {
	// 		allowDeleteOfDefaultGroups := true
	// 		if err := m.DeleteGroupMember(&pb.GroupMember{Group: group, Username: req.GetUsername()}, allowDeleteOfDefaultGroups); err != nil {
	// 			if _, ok := err.(*NoSuchMemberError); !ok {
	// 				return "", 0, fmt.Errorf("failed to remove renamed user (%q -> %q) from group %q: %v", req.GetUsername(), username, group, err)
	// 			}
	// 		}
	// 		allowNonExistent := true
	// 		if err := m.AddGroupMember(&pb.GroupMember{Group: group, Username: username}, allowNonExistent); err != nil {
	// 			if _, ok := err.(*MemberAlreadyExistsError); !ok {
	// 				return "", 0, fmt.Errorf("failed to add renamed user (%q -> %q) to group %q: %v", req.GetUsername(), username, group, err)
	// 			}
	// 		}
	// 		log.Infof("Switched member %q to %q in user group %q ", req.GetUsername(), username, group)
	// 	}
	// }

	// modifyAccountRequest := ldap.NewModifyRequest(
	// 	userDN,
	// 	[]ldap.Control{},
	// )
	// firstName := user.GetAttributeValue("givenName")
	// lastName := user.GetAttributeValue("sn")
	// if update.GetFirstName() != "" {
	// 	firstName = update.GetFirstName()
	// 	modifyAccountRequest.Replace("givenName", []string{firstName})
	// }
	// if update.GetLastName() != "" {
	// 	lastName = update.GetLastName()
	// 	modifyAccountRequest.Replace("sn", []string{lastName})
	// }
	// if update.GetFirstName() != "" || update.GetLastName() != "" {
	// 	fullName := fmt.Sprintf("%s %s", firstName, lastName)
	// 	modifyAccountRequest.Replace("displayName", []string{fullName})
	// 	modifyAccountRequest.Replace("cn", []string{fullName})
	// }
	// if loginShell := update.GetLoginShell(); loginShell != "" {
	// 	modifyAccountRequest.Replace("loginShell", []string{loginShell})
	// }
	// if homeDirectory := update.GetHomeDirectory(); homeDirectory != "" {
	// 	modifyAccountRequest.Replace("homeDirectory", []string{homeDirectory})
	// }
	// if mail := update.GetEmail(); mail != "" {
	// 	modifyAccountRequest.Replace("mail", []string{mail})
	// }
	// if password := update.GetPassword(); password != "" {
	// 	hashedPassword, err := hash.Password(password, algorithm)
	// 	if err != nil {
	// 		return "", 0, fmt.Errorf("failed to hash password: %v", err)
	// 	}
	// 	modifyAccountRequest.Replace("userPassword", []string{hashedPassword})
	// }

	// if isAdmin {
	// 	// Only the admin is allowed to change these because they identify a unique user (username + uidNumber)
	// 	if uid := update.GetUid(); uid >= MinUID {
	// 		uidNumber = int(uid)
	// 		modifyAccountRequest.Replace("uidNumber", []string{strconv.Itoa(int(uid))})
	// 	}
	// 	if gid := update.GetGid(); gid >= MinGID {
	// 		modifyAccountRequest.Replace("gidNumber", []string{strconv.Itoa(int(gid))})
	// 	}
	// }

	// log.Debugf("modifyAccountRequest=%v", modifyAccountRequest)
	// if err := m.ldap.Modify(modifyAccountRequest); err != nil {
	// 	return "", 0, fmt.Errorf("failed to modify existing user: %v", err)
	// }
	// log.Infof("updated %d attributes of user %q", len(modifyAccountRequest.Changes), username)
	// return username, uidNumber, nil
	return "", 0, nil
}

// DeleteAccount ...
func (m *LDAPManager) DeleteUser(req *pb.DeleteUserRequest, keepGroups bool) error {
	// if req.GetUsername() == "" {
	// 	return errors.New("username must not be empty")
	// }
	// if !keepGroups {
	// 	// delete the account from all its groups
	// 	groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{Username: req.GetUsername()})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get list of groups: %v", err)
	// 	}
	// 	for _, group := range groups.GetGroups() {
	// 		allowDeleteOfDefaultGroups := true
	// 		if err := m.DeleteGroupMember(&pb.GroupMember{Group: group, Username: req.GetUsername()}, allowDeleteOfDefaultGroups); err != nil {
	// 			if _, ok := err.(*RemoveLastGroupMemberError); ok {
	// 				return err
	// 			}
	// 			if _, ok := err.(*NoSuchMemberError); !ok {
	// 				return fmt.Errorf("failed to remove deleted user %q from group %q: %v", req.GetUsername(), group, err)
	// 			}
	// 		}
	// 	}
	// }
	// if err := m.ldap.Del(ldap.NewDelRequest(
	// 	fmt.Sprintf("%s=%s,%s", m.AccountAttribute, escapeDN(req.GetUsername()), m.UserGroupDN),
	// 	[]ldap.Control{},
	// )); err != nil {
	// 	return err
	// }
	// log.Infof("removed account %q", req.GetUsername())
	return nil
}
