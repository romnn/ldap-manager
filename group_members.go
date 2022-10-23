package ldapmanager

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
)

// RemoveLastGroupMemberError ...
type RemoveLastGroupMemberError struct {
	ApplicationError
	Group string
}

// RemoveLastGroupMemberError ...
func (e *RemoveLastGroupMemberError) Error() string {
	return fmt.Sprintf("cannot remove the only remaining group member from group %q. consider deleting the group.", e.Group)
}

// Code ...
func (e *RemoveLastGroupMemberError) Code() codes.Code {
	return codes.FailedPrecondition
}

// NoSuchMemberError ...
type NoSuchMemberError struct {
	ApplicationError
	Group, Member string
}

// NoSuchMemberError ...
func (e *NoSuchMemberError) Error() string {
	return fmt.Sprintf("no such member %q in group %q", e.Member, e.Group)
}

// Code ...
func (e *NoSuchMemberError) Code() codes.Code {
	return codes.NotFound
}

// MemberAlreadyExistsError ...
type MemberAlreadyExistsError struct {
	ApplicationError
	Group, Member string
}

// MemberAlreadyExistsError ...
func (e *MemberAlreadyExistsError) Error() string {
	return fmt.Sprintf("member %q is already a member of group %q", e.Member, e.Group)
}

// Code ...
func (e *MemberAlreadyExistsError) Code() codes.Code {
	return codes.AlreadyExists
}

func (m *LDAPManager) getGroup(groupName string) (*pb.Group, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escapeFilter(groupName)),
		[]string{m.GroupMembershipAttribute, "gidNumber"},
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleGroupsError{Group: groupName, Count: len(result.Entries)}
	}
	var members []string
	group := result.Entries[0]
	for _, member := range group.GetAttributeValues(m.GroupMembershipAttribute) {
		members = append(members, member)
	}
	gid, _ := strconv.Atoi(group.GetAttributeValue("gidNumber"))
	return &pb.Group{
		Members: members,
		Name:    groupName,
		Gid:     int32(gid),
	}, nil
}

// IsGroupMember ...
func (m *LDAPManager) IsGroupMember(req *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	var status pb.GroupMemberStatus
	result, err := m.findGroup(req.Group, []string{"dn", m.GroupMembershipAttribute})
	if err != nil {
		return &status, err
	}
	if len(result.Entries) != 1 {
		return &status, &ZeroOrMultipleGroupsError{Group: req.GetGroup(), Count: len(result.Entries)}
	}
	if !m.GroupMembershipUsesUID {
		req.Username = fmt.Sprintf("%s=%s,%s", m.AccountAttribute, req.GetUsername(), m.UserGroupDN)
	}
	for _, member := range result.Entries[0].GetAttributeValues(m.GroupMembershipAttribute) {
		if member == req.GetUsername() {
			return &pb.GroupMemberStatus{IsMember: true}, nil
		}
	}
	return &status, nil
}

// GetUserGroups ...
func (m *LDAPManager) GetUserGroups(req *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	username := escapeDN(req.GetUsername())
	if !m.GroupMembershipUsesUID {
		username = m.AccountNamed(req.GetUsername())
	}
	filter := fmt.Sprintf("(&(objectClass=posixGroup)(%s=%s))", m.GroupMembershipAttribute, username)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"cn"},
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	groupList := &pb.GroupList{Total: int64(len(result.Entries))}
	for _, group := range result.Entries {
		if cn := group.GetAttributeValue("cn"); cn != "" {
			groupList.Groups = append(groupList.Groups, cn)
		}
	}
	// No sorting and clipping here
	return groupList, nil
}

// GetGroup ...
func (m *LDAPManager) GetGroup(req *pb.GetGroupRequest) (*pb.Group, error) {
	group, err := m.getGroup(req.GetName())
	if err != nil {
		return nil, err
	}
	normGroup := &pb.Group{Name: group.GetName(), Gid: group.GetGid(), Total: int64(len(group.GetMembers()))}

	// Convert member DN's to usernames
	for _, memberDN := range group.GetMembers() {
		if memberUsername, err := extractAttribute(memberDN, m.AccountAttribute); err == nil && memberUsername != "" {
			normGroup.Members = append(normGroup.GetMembers(), memberUsername)
		}
	}

	// Sort
	sort.Slice(normGroup.Members, func(i, j int) bool {
		asc := normGroup.Members[i] < normGroup.Members[j]
		if req.GetSortOrder() == pb.SortOrder_DESCENDING {
			return !asc
		}
		return asc
	})
	// Clip
	if req.GetStart() >= 0 && req.GetEnd() < int32(len(normGroup.GetMembers())) && req.GetStart() < req.GetEnd() {
		normGroup.Members = normGroup.Members[req.GetStart():req.GetEnd()]
		return normGroup, nil
	}
	return normGroup, nil
}

// AddGroupMember ...
func (m *LDAPManager) AddGroupMember(req *pb.GroupMember, allowNonExistent bool) error {
	if req.GetGroup() == "" {
		return &ValidationError{Message: "group name must not be empty"}
	}
	if req.GetUsername() == "" {
		return &ValidationError{Message: "username must not be empty"}
	}
	if !allowNonExistent && !m.IsProtectedGroup(req.GetGroup()) {
		memberStatus, err := m.IsGroupMember(&pb.IsGroupMemberRequest{Username: req.GetUsername(), Group: m.DefaultUserGroup})
		if err != nil {
			return fmt.Errorf("failed to check if member %q exists: %v", req.GetUsername(), err)
		}
		if !memberStatus.GetIsMember() {
			return &ZeroOrMultipleAccountsError{
				Username: req.GetUsername(),
			}
		}
	}

	username := escapeDN(req.GetUsername())
	if !m.GroupMembershipUsesUID {
		username = m.AccountNamed(req.GetUsername())
	}
	modifyRequest := ldap.NewModifyRequest(
		m.GroupNamed(req.GetGroup()),
		[]ldap.Control{},
	)
	modifyRequest.Add(m.GroupMembershipAttribute, []string{username})
	log.Debugf("AddGroupMember: modifyRequest=%v", modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultAttributeOrValueExists) {
			return &MemberAlreadyExistsError{Member: req.GetUsername(), Group: req.GetGroup()}
		}
		return err
	}
	log.Infof("added user %q to group %q", username, req.GetGroup())
	return nil
}

// DeleteGroupMember ...
func (m *LDAPManager) DeleteGroupMember(req *pb.GroupMember, allowDeleteOfDefaultGroups bool) error {
	if req.GetGroup() == "" || req.GetUsername() == "" {
		return &ValidationError{Message: "group and user name can not be empty"}
	}
	if !allowDeleteOfDefaultGroups && m.IsProtectedGroup(req.GetGroup()) {
		return &ValidationError{Message: "deleting members from the default user or admin group is not allowed"}
	}
	username := escapeDN(req.GetUsername())
	if !m.GroupMembershipUsesUID {
		username = m.AccountNamed(req.GetUsername())
	}
	modifyRequest := ldap.NewModifyRequest(
		m.GroupNamed(req.GetGroup()),
		[]ldap.Control{},
	)
	modifyRequest.Delete(m.GroupMembershipAttribute, []string{username})
	log.Debugf("DeleteGroupMember: modifyRequest=%v", modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultObjectClassViolation) {
			return &RemoveLastGroupMemberError{Group: req.GetGroup()}
		}
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) || ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchAttribute) {
			return &NoSuchMemberError{Group: req.GetGroup(), Member: req.GetUsername()}
		}
		return err
	}
	log.Infof("removed user %q from group %q", username, req.GetGroup())
	return nil
}
