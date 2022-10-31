package pkg
import ( // "github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

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
