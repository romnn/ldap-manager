package groups

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// GroupAlreadyExistsError ...
type GroupAlreadyExistsError struct {
	ApplicationError
	Group string
}

// GroupAlreadyExistsError ...
func (e *GroupAlreadyExistsError) Error() string {
	return fmt.Sprintf("group %q already exists", e.Group)
}

// Code ...
func (e *GroupAlreadyExistsError) Code() codes.Code {
	return codes.AlreadyExists
}

// ZeroOrMultipleGroupsError ...
type ZeroOrMultipleGroupsError struct {
	ApplicationError
	Group string
	Count int
}

// Error ...
func (e *ZeroOrMultipleGroupsError) Error() string {
	if e.Count > 1 {
		return fmt.Sprintf("multiple (%d) groups with name %q", e.Count, e.Group)
	}
	return fmt.Sprintf("no group with name %q", e.Group)
}

// Code ...
func (e *ZeroOrMultipleGroupsError) Code() codes.Code {
	if e.Count > 1 {
		return codes.Internal
	}
	return codes.NotFound
}

func (m *LDAPManager) getGroupGID(groupName string) (int, error) {
	if groupName == "" {
		return 0, &ValidationError{Message: "group name can not be empty"}
	}
	result, err := m.findGroup(groupName, []string{"gidNumber"})
	if err != nil {
		return 0, err
	}
	if len(result.Entries) != 1 {
		return 0, &ZeroOrMultipleGroupsError{Group: groupName, Count: len(result.Entries)}
	}
	return strconv.Atoi(result.Entries[0].GetAttributeValue("gidNumber"))
}

// IsProtectedGroup ...
func (m *LDAPManager) IsProtectedGroup(groupName string) bool {
	isAdminGroup := strings.ToLower(groupName) == strings.ToLower(m.DefaultAdminGroup)
	isUserGroup := strings.ToLower(groupName) == strings.ToLower(m.DefaultUserGroup)
	return isAdminGroup || isUserGroup
}

// GroupNamed ...
func (m *LDAPManager) GroupNamed(name string) string {
	return fmt.Sprintf("cn=%s,%s", escapeDN(name), m.GroupsDN)
}

// NewGroup ...
func (m *LDAPManager) NewGroup(req *pb.NewGroupRequest, strict bool) error {
	if req.GetName() == "" {
		return &ValidationError{Message: "group name can not be empty"}
	}
	result, err := m.findGroup(req.GetName(), []string{"dn", m.GroupMembershipAttribute})
	if err != nil {
		return err
	}
	if len(result.Entries) > 0 {
		return &GroupAlreadyExistsError{Group: req.GetName()}
	}
	highestGID, err := m.getHighestID(m.GroupAttribute)
	if err != nil {
		return err
	}
	newGID := highestGID + 1

	var memberList []string
	for _, username := range req.GetMembers() {
		if strict {
			memberStatus, err := m.IsGroupMember(&pb.IsGroupMemberRequest{Username: username, Group: m.DefaultUserGroup})
			if err != nil {
				return fmt.Errorf("failed to check if member %q exists: %v", username, err)
			}
			if !memberStatus.GetIsMember() {
				log.Warnf("Skipping user %q to be added to group %q because it is not in the %q group", username, req.GetName(), m.DefaultUserGroup)
				continue
			}
		}
		member := escapeDN(username)
		if !m.GroupMembershipUsesUID {
			member = m.AccountNamed(username)
		}
		memberList = append(memberList, member)
	}

	var groupAttributes []ldap.Attribute
	if !m.UseRFC2307BISSchema {
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"top", "posixGroup"}},
			{Type: "cn", Vals: []string{escapeDN(req.GetName())}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(newGID)}},
		}
	} else {
		if len(memberList) < 1 {
			return &ValidationError{Message: "when using RFC2307BIS (not NIS), you must specify at least one existing group member"}
		}
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"top", "groupOfUniqueNames", "posixGroup"}},
			{Type: "cn", Vals: []string{escapeDN(req.GetName())}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(newGID)}},
		}
	}

	groupAttributes = append(groupAttributes, ldap.Attribute{
		Type: m.GroupMembershipAttribute, Vals: memberList,
	})

	addGroupRequest := &ldap.AddRequest{
		DN:         m.GroupNamed(req.GetName()),
		Attributes: groupAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debugf("addGroupRequest=%v", addGroupRequest)
	if err := m.ldap.Add(addGroupRequest); err != nil {
		return err
	}
	if err := m.updateLastID("lastGID", newGID); err != nil {
		return err
	}
	log.Infof("added new group %q with %d members (gid=%d)", req.GetName(), len(memberList), newGID)
	return nil
}

// DeleteGroup ...
func (m *LDAPManager) DeleteGroup(req *pb.DeleteGroupRequest) error {
	if req.GetName() == "" {
		return &ValidationError{Message: "group name can not be empty"}
	}
	if m.IsProtectedGroup(req.GetName()) {
		return &ValidationError{Message: "deleting the default user or admin group is not allowed"}
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		m.GroupNamed(req.GetName()),
		[]ldap.Control{},
	)); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			return &ZeroOrMultipleGroupsError{Group: req.GetName()}
		}
		return err
	}
	log.Infof("removed group %q", req.GetName())
	return nil
}

// UpdateGroup ...
func (m *LDAPManager) UpdateGroup(req *pb.UpdateGroupRequest) error {
	if req.GetName() == "" {
		return &ValidationError{Message: "group name can not be empty"}
	}

	groupName := req.GetName()
	if req.GetNewName() != "" && req.GetNewName() != groupName {
		modifyRequest := &ldap.ModifyDNRequest{
			DN:           m.GroupNamed(groupName),
			NewRDN:       fmt.Sprintf("cn=%s", req.GetNewName()),
			DeleteOldRDN: true,
			NewSuperior:  "",
		}
		log.Debugf("UpdateGroup modifyRequest=%v", modifyRequest)
		if err := m.ldap.ModifyDN(modifyRequest); err != nil {
			return err
		}
		log.Infof("renamed group from %q to %q", req.GetName(), req.GetNewName())
		groupName = req.GetNewName()
	}

	modifyGroupRequest := ldap.NewModifyRequest(
		m.GroupNamed(groupName),
		[]ldap.Control{},
	)
	if req.GetGid() >= MinGID {
		modifyGroupRequest.Replace("gidNumber", []string{strconv.Itoa(int(req.GetGid()))})
	}
	if err := m.ldap.Modify(modifyGroupRequest); err != nil {
		return fmt.Errorf("failed to modify group %q: %v", groupName, err)
	}
	log.Infof("updated %d attributes of group %q", len(modifyGroupRequest.Changes), groupName)
	return nil
}

func (m *LDAPManager) countGroups() (int, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(objectClass=posixGroup)",
		[]string{"cn"},
		[]ldap.Control{},
	))
	if err != nil {
		return 0, err
	}
	return len(result.Entries), nil
}

// GetGroupList ...
func (m *LDAPManager) GetGroupList(req *pb.GetGroupListRequest) (*pb.GroupList, error) {
	filter := parseFilter(req.Filter)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=posixGroup)%s)", filter),
		[]string{},
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	total, err := m.countGroups()
	if err != nil {
		return nil, err
	}
	groupList := &pb.GroupList{Total: int64(total)}
	for _, group := range result.Entries {
		if cn := group.GetAttributeValue("cn"); cn != "" {
			groupList.Groups = append(groupList.Groups, cn)
		}
	}
	// Sort
	groups := groupList.GetGroups()
	sort.Slice(groups, func(i, j int) bool {
		asc := groups[i] < groups[j]
		if req.GetSortOrder() == pb.SortOrder_DESCENDING {
			return !asc
		}
		return asc
	})
	// Clip
	if req.GetStart() >= 0 && req.GetEnd() < int32(len(groups)) && req.GetStart() < req.GetEnd() {
		groupList.Groups = groups[req.GetStart():req.GetEnd()]
	}
	return groupList, nil
}
