package main

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap"
	"github.com/romnnn/ldap-manager/hash"
	log "github.com/sirupsen/logrus"
)

const (
	// MinUID for POSIX accounts
	MinUID = 2000
	// MinGID for POSIX accounts
	MinGID = 2000

	// SortAscending ...
	SortAscending = "asc"
	// SortDescending ...
	SortDescending = "desc"
)

// DeleteGroup ...
func (s *LDAPManager) DeleteGroup(groupName string) error {
	if err := s.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("cn=%s,%s", escape(groupName), s.GroupsDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	log.Infof("removed group %q", groupName)
	return nil
}

// DeleteAccount ...
func (s *LDAPManager) DeleteAccount(username string) error {
	if username == "" {
		return errors.New("username must not be empty")
	}
	if err := s.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("%s=%s,%s", s.AccountAttribute, escape(username), s.UserGroupDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	log.Infof("removed account %q", username)
	return nil
}

// GetGroupGID ...
func (s *LDAPManager) GetGroupGID(groupName string) (int, error) {
	result, err := s.findGroup(groupName, []string{"gidNumber"})
	if err != nil {
		return 0, err
	}
	if len(result.Entries) != 1 {
		return 0, fmt.Errorf("group %q does not exist or too many entries returned", groupName)
	}
	gidNumbers := result.Entries[0].GetAttributeValues("gidNumber")
	if len(gidNumbers) != 1 {
		return 0, fmt.Errorf("group %q does not have gidNumber or multiple", groupName)
	}
	return strconv.Atoi(gidNumbers[0])
}

// IsGroupMember ...
func (s *LDAPManager) IsGroupMember(username, groupName string) (bool, error) {
	result, err := s.findGroup(groupName, []string{"dn", s.GroupMembershipAttribute})
	if err != nil {
		return false, err
	}
	if len(result.Entries) != 1 {
		return false, fmt.Errorf("user %q does not exist or too many entries returned", username)
	}
	if !s.GroupMembershipUsesUID {
		// "${LDAP['account_attribute']}=$username,${LDAP['user_dn']}";
		username = fmt.Sprintf("%s=%s,%s", s.AccountAttribute, username, s.UserGroupDN)
	}
	// preg_grep ("/^${username}$/i", $result[0][$LDAP['group_membership_attribute']])
	for _, member := range result.Entries[0].GetAttributeValues(s.GroupMembershipAttribute) { // uniqueMember or memberUID
		if member == username {
			return true, nil
		}
	}
	return false, nil
}

// NewAccountRequest ...
type NewAccountRequest struct {
	FirstName, LastName, Username, Password, Email string
}

// Validate ...
func (req *NewAccountRequest) Validate() error {
	if req.Username == "" {
		return errors.New("Must specify username")
	}
	if req.Password == "" {
		return errors.New("Must specify password")
	}
	return nil
	/* Optional
	if req.Email == "" {
		return errors.New("Must specify email")
	}
	*/
}

// GetHighestID ...
func (s *LDAPManager) GetHighestID(attribute string) (int, error) {
	var highestID int
	var entryBaseDN, entryFilter, entryAttribute string
	switch attribute {
	case s.GroupAttribute:
		highestID = MinGID
		entryBaseDN = s.GroupsDN
		entryFilter = "(objectClass=posixGroup)"
		entryAttribute = "gidNumber"
	case s.AccountAttribute:
		highestID = MinUID
		entryBaseDN = s.UserGroupDN
		entryFilter = fmt.Sprintf("(%s=*)", s.AccountAttribute)
		entryAttribute = "uidNumber"
	default:
		return highestID, fmt.Errorf("unknown id attribute %q", attribute)
	}

	filter := fmt.Sprintf("(&(objectclass=device)(cn=last%s))", attribute)
	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"serialNumber"},
		[]ldap.Control{},
	))
	if err != nil {
		return highestID, err
	}
	// Check for cached lastUID / lastGID value first
	if len(result.Entries) > 0 {
		if fetchedID, err := strconv.Atoi(result.Entries[0].GetAttributeValue("serialNumber")); err == nil && fetchedID >= highestID {
			return fetchedID, nil
		}
	}

	// cache miss requires traversing all entries
	result, err = s.ldap.Search(ldap.NewSearchRequest(
		entryBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		entryFilter,
		[]string{entryAttribute},
		[]ldap.Control{},
	))
	if err != nil {
		return highestID, err
	}
	for _, entry := range result.Entries {
		if entryAttrValue := entry.GetAttributeValue(entryAttribute); entryAttrValue != "" {
			if entryAttrNumericValue, err := strconv.Atoi(entryAttrValue); err == nil {
				if entryAttrNumericValue > highestID {
					highestID = entryAttrNumericValue
				}
			}
		}
	}
	return highestID, nil
}

// GroupExistsError ...
type GroupExistsError struct {
	Group string
}

// GroupExistsError ...
func (e *GroupExistsError) Error() string {
	return fmt.Sprintf("group %q already exists", e.Group)
}

// NewGroup ...
func (s *LDAPManager) NewGroup(name string, members []string) error {
	if name == "" {
		return errors.New("group name can not be empty")
	}
	result, err := s.findGroup(name, []string{"dn", s.GroupMembershipAttribute})
	if err != nil {
		return err
	}
	if len(result.Entries) > 0 {
		return &GroupExistsError{Group: name}
	}
	highestGID, err := s.GetHighestID(s.GroupAttribute)
	if err != nil {
		return err
	}
	newGID := highestGID + 1

	var groupAttributes []ldap.Attribute
	if s.UseNISSchema {
		groupAttributes = []ldap.Attribute{
			{"objectClass", []string{"top", "posixGroup"}},
			{"cn", []string{name}},
			{"gidNumber", []string{strconv.Itoa(newGID)}},
		}
	} else {
		if len(members) < 1 {
			return errors.New("when using RFC2307BIS (not NIS), you must specify at least one group member")
		}
		groupAttributes = []ldap.Attribute{
			{"objectClass", []string{"top", "groupOfUniqueNames", "posixGroup"}},
			{"cn", []string{name}},
			{"gidNumber", []string{strconv.Itoa(newGID)}},
			{s.GroupMembershipAttribute, members},
		}
	}
	addGroupRequest := &ldap.AddRequest{
		DN:         fmt.Sprintf("cn=%s,%s", name, s.GroupsDN),
		Attributes: groupAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(addGroupRequest)
	if err := s.ldap.Add(addGroupRequest); err != nil {
		return err
	}
	if err := s.updateLastID("lastGID", newGID); err != nil {
		return err
	}
	log.Infof("added new group %q (gid=%d)", name, newGID)
	return nil
}

// GetGroupMembers ...
func (s *LDAPManager) GetGroupMembers(groupName string, start, end int, sortOrder string) ([]string, error) {
	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escape(groupName)),
		[]string{s.GroupMembershipAttribute},
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, fmt.Errorf("zero or multiple groups with name %q", groupName)
	}
	var members []string
	group := result.Entries[0]
	for _, member := range group.GetAttributeValues(s.GroupMembershipAttribute) {
		log.Info(member)
		// TODO
		/*
			reg, err := regexp.Compile(fmt.Sprintf("%s=(.*?),", s.AccountAttribute))
			if err != nil {
				return "", errors.New("failed to compile regex")
			}
			matchedDN := reg.FindString(userDN)
		*/

		// if member.Key != "count" && member.Value != "" {
		// $this_member = preg_replace("/^.*?=(.*?),.*/", "$1", $value);
		// array_push($records, $this_member);
		// }
	}

	// Sort
	sort.Slice(members, func(i, j int) bool {
		asc := members[i] < members[j]
		if sortOrder == SortDescending {
			return !asc
		}
		return asc
	})
	// Clip
	if start >= 0 && end < len(members) && start < end {
		return members[start:end], nil
	}
	return members, nil
}

// GetGroupList ...
func (s *LDAPManager) GetGroupList(start, end int, sortOrder string, filters []string) ([]string, error) {
	filter := fmt.Sprintf("(&(objectclass=*)%s)", filters)
	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{},
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	var groups []string
	for _, group := range result.Entries {
		if cn := group.GetAttributeValue("cn"); cn != "" {
			groups = append(groups, cn)
		}
	}
	// Sort
	sort.Slice(groups, func(i, j int) bool {
		asc := groups[i] < groups[j]
		if sortOrder == SortDescending {
			return !asc
		}
		return asc
	})
	// Clip
	if start >= 0 && end < len(groups) && start < end {
		return groups[start:end], nil
	}
	return groups, nil
}

// GetUserList ...
func (s *LDAPManager) GetUserList(start, end int, sortOrder string, sortKey string, filters string, fields []string) ([]map[string]string, error) {
	if len(fields) < 1 {
		fields = []string{s.AccountAttribute, "givenname", "sn", "mail"}
	}
	if sortKey == "" {
		sortKey = s.AccountAttribute
	}
	filter := fmt.Sprintf("(&(%s=*)%s)", s.AccountAttribute, filters)
	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		fields,
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	var users map[string]map[string]string
	for _, entry := range result.Entries {
		log.Info(entry)
		if entryKey := entry.GetAttributeValue(sortKey); entryKey != "" {
			user := make(map[string]string)
			for _, attr := range entry.Attributes {
				// if attr.Name != sortKey {
				user[attr.Name] = entry.GetAttributeValue(attr.Name)
				// }
			}
			users[entryKey] = user
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
		if sortOrder == SortDescending {
			return !asc
		}
		return asc
	})
	// Clip
	clippedKeys := keys
	var clippedUsers []map[string]string
	if start >= 0 && end < len(keys) && start < end {
		clippedKeys = keys[start:end]
	}
	for _, key := range clippedKeys {
		clippedUsers = append(clippedUsers, users[key])
	}
	return clippedUsers, nil
}

// AddGroupMember ...
func (s *LDAPManager) AddGroupMember(groupName string, username string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", escape(groupName), s.GroupsDN)
	if !s.GroupMembershipUsesUID {
		username = fmt.Sprintf("%s=%s,%s", s.AccountAttribute, username, s.UserGroupDN)
	}

	addGroupMemberRequest := &ldap.AddRequest{
		DN: groupDN,
		Attributes: []ldap.Attribute{
			{s.GroupMembershipAttribute, []string{username}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addGroupMemberRequest)
	if err := s.ldap.Add(addGroupMemberRequest); err != nil {
		return err
	}
	log.Infof("added user %q to group %q", username, groupName)
	return nil
}

// DeleteGroupMember ...
func (s *LDAPManager) DeleteGroupMember(groupName string, username string) error {
	groupDN := fmt.Sprintf("cn=%s,%s", escape(groupName), s.GroupsDN)
	if !s.GroupMembershipUsesUID {
		username = fmt.Sprintf("%s=%s,%s", s.AccountAttribute, username, s.UserGroupDN)
	}

	modifyRequest := ldap.NewModifyRequest(
		groupDN,
		[]ldap.Control{},
	)
	modifyRequest.Delete(s.GroupMembershipAttribute, []string{username})
	log.Debug(modifyRequest)
	if err := s.ldap.Modify(modifyRequest); err != nil {
		return err
	}
	log.Infof("removed user %q from group %q", username, groupName)
	return nil
}

// AuthenticateUser ...
func (s *LDAPManager) AuthenticateUser(username string, password string) (string, error) {
	// Validate
	if username == "" || password == "" {
		return "", errors.New("must provide username and password")
	}
	// Search for the DN for the given username. If found, try binding with the DN and user's password.
	// If the binding succeeds, return the DN.
	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", s.AccountAttribute, escape(username), s.UserGroupDN),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return "", err
	}
	if len(result.Entries) != 1 {
		return "", fmt.Errorf("zero or multiple accounts with username %q", username)
	}
	userDN := result.Entries[0].DN
	if err := s.ldap.Bind(userDN, password); err != nil {
		return "", fmt.Errorf("unable to bind as %q", username)
	}
	reg, err := regexp.Compile(fmt.Sprintf("%s=(.*?),", s.AccountAttribute))
	if err != nil {
		return "", errors.New("failed to compile regex")
	}
	matchedDN := reg.FindString(userDN)
	return matchedDN, nil
}

// NewAccount ...
func (s *LDAPManager) NewAccount(req *NewAccountRequest) error {
	// Validate
	if err := req.Validate(); err != nil {
		return err
	}
	// Check for existing user with the same username
	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", s.AccountAttribute, escape(req.Username), s.UserGroupDN),
		[]string{},
		[]ldap.Control{},
	))
	if err != nil {
		return err
	}
	if len(result.Entries) > 0 {
		return fmt.Errorf("account with username %q already exists", req.Username)
	}
	highestUID, err := s.GetHighestID(s.AccountAttribute)
	if err != nil {
		return err
	}
	newUID := highestUID + 1
	group := s.DefaultUserGroup
	userDN := fmt.Sprintf("%s=%s,%s", s.AccountAttribute, req.Username, s.UserGroupDN)

	var GID int
	if defaultGID, err := s.GetGroupGID(s.DefaultUserGroup); err == nil {
		GID = defaultGID
	} else {
		// The default user group might not yet exist
		// Note that a group can only be created with at least one member when using RFC2307BIS
		if err := s.NewGroup(s.DefaultUserGroup, []string{userDN}); err != nil {
			// Fall back to create a new group group for the user
			if err := s.NewGroup(req.Username, []string{userDN}); err != nil {
				if _, ok := err.(*GroupExistsError); !ok {
					return err
				}
			}
			group = req.Username
		}

		userGroupGID, err := s.GetGroupGID(group)
		if err != nil {
			return err
		}
		GID = userGroupGID
	}

	hashedPassword, err := hash.Password(req.Password, hash.SHA512CRYPT)
	if err != nil {
		return err
	}
	log.Info(hashedPassword)

	addUserRequest := &ldap.AddRequest{
		DN: userDN,
		Attributes: []ldap.Attribute{
			{"objectClass", []string{"person", "inetOrgPerson", "posixAccount"}},
			{"uid", []string{req.Username}},
			{"givenName", []string{req.FirstName}},
			{"sn", []string{req.LastName}},
			{"cn", []string{fmt.Sprintf("%s %s", req.FirstName, req.LastName)}},
			{"displayName", []string{fmt.Sprintf("%s %s", req.FirstName, req.LastName)}},
			{"uidNumber", []string{strconv.Itoa(newUID)}},
			{"gidNumber", []string{strconv.Itoa(GID)}},
			{"loginShell", []string{s.DefaultUserShell}},
			{"homeDirectory", []string{fmt.Sprintf("/home/%s", req.Username)}},
			{"userPassword", []string{hashedPassword}},
			{"mail", []string{req.Email}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addUserRequest)
	// TODO
	/*
		if err := s.ldap.Add(addUserRequest); err != nil {
			return err
		}
	*/
	if err := s.AddGroupMember(group, req.Username); err != nil {
		return err
	}
	if err := s.updateLastID("lastUID", newUID); err != nil {
		return err
	}
	log.Infof("added new account %q (member of group %q)", req.Username, group)
	return nil
}

// ChangePassword ...
func (s *LDAPManager) ChangePassword(username, newPassword string) error {
	// Validate
	if username == "" || newPassword == "" {
		return errors.New("username and password must not be empty")
	}

	result, err := s.ldap.Search(ldap.NewSearchRequest(
		s.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=%s,%s)", s.AccountAttribute, escape(username), s.UserGroupDN),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return err
	}
	if len(result.Entries) != 1 {
		return fmt.Errorf("zero or multiple (%d) accounts with username %q", len(result.Entries), username)
	}
	userDN := result.Entries[0].DN
	hashedPassword, err := hash.Password(newPassword, hash.Default)
	if err != nil {
		return err
	}
	modifyPasswordRequest := ldap.NewModifyRequest(
		userDN,
		[]ldap.Control{},
	)
	modifyPasswordRequest.Replace("userPassword", []string{hashedPassword})
	log.Debug(modifyPasswordRequest)
	if err := s.ldap.Modify(modifyPasswordRequest); err != nil {
		return err
	}
	log.Infof("changed password for user %q", username)
	return nil
}
