package ldapmanager

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// GroupAlreadyExistsError ...
type GroupAlreadyExistsError struct {
	Group string
}

// GroupAlreadyExistsError ...
func (e *GroupAlreadyExistsError) Error() string {
	return fmt.Sprintf("group %q already exists", e.Group)
}

// GroupValidationError ...
type GroupValidationError struct {
	Message string
}

// Error ...
func (e *GroupValidationError) Error() string {
	return e.Message
}

// ZeroOrMultipleGroupsError ...
type ZeroOrMultipleGroupsError struct {
	Group string
	Count int
}

// Status ...
func (e *ZeroOrMultipleGroupsError) Status() int {
	if e.Count > 1 {
		return http.StatusConflict
	}
	return http.StatusNotFound
}

// Error ...
func (e *ZeroOrMultipleGroupsError) Error() string {
	if e.Count > 1 {
		return fmt.Sprintf("multiple (%d) groups with name %q", e.Count, e.Group)
	}
	return fmt.Sprintf("no group with name %q", e.Group)
}

func (m *LDAPManager) getGroupGID(groupName string) (int, error) {
	if groupName == "" {
		return 0, &GroupValidationError{"group name can not be empty"}
	}
	result, err := m.findGroup(groupName, []string{"gidNumber"})
	if err != nil {
		return 0, err
	}
	if len(result.Entries) != 1 {
		return 0, &ZeroOrMultipleGroupsError{Group: groupName, Count: len(result.Entries)}
	}
	gidNumbers := result.Entries[0].GetAttributeValues("gidNumber")
	if len(gidNumbers) != 1 {
		return 0, fmt.Errorf("group %q does not have gidNumber or multiple", groupName)
	}
	return strconv.Atoi(gidNumbers[0])
}

// NewGroupRequest ...
type NewGroupRequest struct {
	Name    string   `json:"name" form:"name"`
	Members []string `json:"members" form:"members"`
}

// NewGroup ...
func (m *LDAPManager) NewGroup(req *NewGroupRequest) error {
	if req.Name == "" {
		return &GroupValidationError{"group name can not be empty"}
	}
	result, err := m.findGroup(req.Name, []string{"dn", m.GroupMembershipAttribute})
	if err != nil {
		return err
	}
	if len(result.Entries) > 0 {
		return &GroupAlreadyExistsError{Group: req.Name}
	}
	highestGID, err := m.getHighestID(m.GroupAttribute)
	if err != nil {
		return err
	}
	newGID := highestGID + 1

	var groupAttributes []ldap.Attribute
	if !m.UseRFC2307BISSchema {
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"top", "posixGroup"}},
			{Type: "cn", Vals: []string{req.Name}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(newGID)}},
		}
	} else {
		if len(req.Members) < 1 {
			return &GroupValidationError{"when using RFC2307BIS (not NIS), you must specify at least one group member"}
		}
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"top", "groupOfUniqueNames", "posixGroup"}},
			{Type: "cn", Vals: []string{req.Name}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(newGID)}},
			// TODO: Do not expect full userDN's here
			{Type: m.GroupMembershipAttribute, Vals: req.Members},
		}
	}
	addGroupRequest := &ldap.AddRequest{
		DN:         fmt.Sprintf("cn=%s,%s", req.Name, m.GroupsDN),
		Attributes: groupAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(addGroupRequest)
	if err := m.ldap.Add(addGroupRequest); err != nil {
		return err
	}
	if err := m.updateLastID("lastGID", newGID); err != nil {
		return err
	}
	log.Infof("added new group %q (gid=%d)", req.Name, newGID)
	return nil
}

// DeleteGroup ...
func (m *LDAPManager) DeleteGroup(groupName string) error {
	if groupName == "" {
		return &GroupValidationError{"group name can not be empty"}
	}
	isAdminGroup := strings.ToLower(groupName) == strings.ToLower(m.DefaultAdminGroup)
	isUserGroup := strings.ToLower(groupName) == strings.ToLower(m.DefaultUserGroup)
	if isAdminGroup || isUserGroup {
		return &GroupValidationError{"deleting the default user or admin group is not allowed"}
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("cn=%s,%s", escapeDN(groupName), m.GroupsDN),
		[]ldap.Control{},
	)); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			return &ZeroOrMultipleGroupsError{Group: groupName}
		}
		return err
	}
	log.Infof("removed group %q", groupName)
	return nil
}

/* GetGroup ...
func (m *LDAPManager) GetGroup(groupName string) (string, error) {
	if groupName == "" {
		return "", &GroupValidationError{"group name can not be empty"}
	}
	result, err := m.findGroup(groupName, []string{"gidNumber"})
	if err != nil {
		return "", err
	}
	if len(result.Entries) != 1 {
		return "", &ZeroOrMultipleGroupsError{Group: groupName, Count: len(result.Entries)}
	}
	// TODO: What do we want to return here
	return result.Entries[0].DN, nil
}
*/

// RenameGroup ...
func (m *LDAPManager) RenameGroup(groupName, newName string) error {
	if groupName == "" || newName == "" {
		return &GroupValidationError{"group name can not be empty"}
	}
	modifyRequest := &ldap.ModifyDNRequest{
		DN:           fmt.Sprintf("cn=%s,%s", groupName, m.GroupsDN),
		NewRDN:       fmt.Sprintf("cn=%s", newName),
		DeleteOldRDN: true,
		NewSuperior:  "",
	}
	log.Debug(modifyRequest)
	if err := m.ldap.ModifyDN(modifyRequest); err != nil {
		return err
	}
	log.Infof("renames group from %q to %q", groupName, newName)
	return nil
}

// GetGroupListRequest ...
type GetGroupListRequest struct {
	ListOptions
	Filters string
}

// GetGroupList ...
func (m *LDAPManager) GetGroupList(req *GetGroupListRequest) ([]string, error) {
	filter := fmt.Sprintf("(&(objectClass=*)%s)", req.Filters)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
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
		if req.SortOrder == SortDescending {
			return !asc
		}
		return asc
	})
	// Clip
	if req.Start >= 0 && req.End < len(groups) && req.Start < req.End {
		return groups[req.Start:req.End], nil
	}
	return groups, nil
}
