package ldapmanager

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	log "github.com/sirupsen/logrus"
)

// RemoveLastGroupMemberError ...
type RemoveLastGroupMemberError struct {
	Group string
}

// RemoveLastGroupMemberError ...
func (e *RemoveLastGroupMemberError) Error() string {
	return fmt.Sprintf("cannot remove the only remaining group member from group %q. consider deleting the group.", e.Group)
}

// NoSuchMemberError ...
type NoSuchMemberError struct {
	Group, Member string
}

// NoSuchMemberError ...
func (e *NoSuchMemberError) Error() string {
	return fmt.Sprintf("no such member %q in group %q", e.Member, e.Group)
}

func (m *LDAPManager) getGroup(groupName string) (*pb.Group, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escapeFilter(groupName)),
		[]string{m.GroupMembershipAttribute},
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
		log.Info(member)
		members = append(members, member)
	}
	return &pb.Group{
		Members: members,
		Name:    groupName,
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

// GetGroup ...
func (m *LDAPManager) GetGroup(req *pb.GetGroupRequest) (*pb.Group, error) {
	group, err := m.getGroup(req.GetName())
	if err != nil {
		return nil, err
	}
	normGroup := &pb.Group{Name: group.GetName()}

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
	if req.GetGroup() == "" || req.GetUsername() == "" {
		return &GroupValidationError{"group and user name can not be empty"}
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
	log.Debug(modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		return err
	}
	log.Infof("added user %q to group %q", username, req.GetGroup())
	return nil
}

// DeleteGroupMember ...
func (m *LDAPManager) DeleteGroupMember(req *pb.GroupMember, allowDeleteOfDefaultGroups bool) error {
	if req.GetGroup() == "" || req.GetUsername() == "" {
		return &GroupValidationError{"group and user name can not be empty"}
	}
	if !allowDeleteOfDefaultGroups && m.IsProtectedGroup(req.GetGroup()) {
		return &GroupValidationError{"deleting members from the default user or admin group is not allowed"}
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
	log.Debug(modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultObjectClassViolation) {
			return &RemoveLastGroupMemberError{req.GetGroup()}
		}
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) || ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchAttribute) {
			return &NoSuchMemberError{Group: req.GetGroup(), Member: req.GetUsername()}
		}
		return err
	}
	log.Infof("removed user %q from group %q", username, req.GetGroup())
	return nil
}
