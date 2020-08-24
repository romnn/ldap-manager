package ldapmanager

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap"
	"github.com/neko-neko/echo-logrus/v2/log"
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

// Group ...
type Group struct {
	Members []string `json:"members" form:"members"`
	Name    string   `json:"name" form:"name"`
	DN      string   `json:"dn" form:"dn"`
}

func (m *LDAPManager) getGroup(groupName string) (*Group, error) {
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
	return &Group{
		Members: members,
		Name:    groupName,
		DN:      group.DN,
	}, nil
}

// IsGroupMemberRequest ...
type IsGroupMemberRequest struct {
	Username string `json:"username" form:"username"`
	Group    string `json:"group" form:"group"`
}

// IsGroupMember ...
func (m *LDAPManager) IsGroupMember(req *IsGroupMemberRequest) (bool, error) {
	result, err := m.findGroup(req.Group, []string{"dn", m.GroupMembershipAttribute})
	if err != nil {
		return false, err
	}
	if len(result.Entries) != 1 {
		return false, &ZeroOrMultipleGroupsError{Group: req.Group, Count: len(result.Entries)}
	}
	if !m.GroupMembershipUsesUID {
		req.Username = fmt.Sprintf("%s=%s,%s", m.AccountAttribute, req.Username, m.UserGroupDN)
	}
	for _, member := range result.Entries[0].GetAttributeValues(m.GroupMembershipAttribute) {
		if member == req.Username {
			return true, nil
		}
	}
	return false, nil
}

// GetGroupRequest ...
type GetGroupRequest struct {
	Options ListOptions `json:"options" form:"options"`
	Group   string      `json:"group" form:"group"`
}

// GetGroup ...
func (m *LDAPManager) GetGroup(req *GetGroupRequest) (*Group, error) {
	group, err := m.getGroup(req.Group)
	if err != nil {
		return nil, err
	}
	normGroup := &Group{Name: group.Name}

	// Convert member DN's to usernames
	for _, memberDN := range group.Members {
		if memberUsername, err := extractAttribute(memberDN, m.AccountAttribute); err == nil && memberUsername != "" {
			normGroup.Members = append(normGroup.Members, memberUsername)
		}
	}

	// Sort
	sort.Slice(normGroup.Members, func(i, j int) bool {
		asc := normGroup.Members[i] < normGroup.Members[j]
		if req.Options.SortOrder == SortDescending {
			return !asc
		}
		return asc
	})
	// Clip
	if req.Options.Start >= 0 && req.Options.End < len(normGroup.Members) && req.Options.Start < req.Options.End {
		normGroup.Members = normGroup.Members[req.Options.Start:req.Options.End]
		return normGroup, nil
	}
	return normGroup, nil
}

// AddGroupMemberRequest ...
type AddGroupMemberRequest struct {
	Username         string `json:"username" form:"username"`
	Group            string `json:"group" form:"group"`
	AllowNonExistent bool
}

// AddGroupMember ...
func (m *LDAPManager) AddGroupMember(req *AddGroupMemberRequest) error {
	if req.Group == "" || req.Username == "" {
		return &GroupValidationError{"group and user name can not be empty"}
	}
	if !req.AllowNonExistent && !m.IsProtectedGroup(req.Group) {
		isMember, err := m.IsGroupMember(&IsGroupMemberRequest{Username: req.Username, Group: m.DefaultUserGroup})
		if err != nil {
			return fmt.Errorf("failed to check if member %q exists: %v", req.Username, err)
		}
		if !isMember {
			return &ZeroOrMultipleAccountsError{
				Username: req.Username,
			}
		}
	}

	username := escapeDN(req.Username)
	if !m.GroupMembershipUsesUID {
		username = m.AccountNamed(req.Username)
	}
	modifyRequest := ldap.NewModifyRequest(
		m.GroupNamed(req.Group),
		[]ldap.Control{},
	)
	modifyRequest.Add(m.GroupMembershipAttribute, []string{username})
	log.Debug(modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		return err
	}
	log.Infof("added user %q to group %q", username, req.Group)
	return nil
}

// DeleteGroupMemberRequest ...
type DeleteGroupMemberRequest struct {
	Username                   string `json:"username" form:"username"`
	Group                      string `json:"group" form:"group"`
	AllowDeleteOfDefaultGroups bool
}

// DeleteGroupMember ...
func (m *LDAPManager) DeleteGroupMember(req *DeleteGroupMemberRequest) error {
	if req.Group == "" || req.Username == "" {
		return &GroupValidationError{"group and user name can not be empty"}
	}
	if !req.AllowDeleteOfDefaultGroups && m.IsProtectedGroup(req.Group) {
		return &GroupValidationError{"deleting members from the default user or admin group is not allowed"}
	}
	username := escapeDN(req.Username)
	if !m.GroupMembershipUsesUID {
		username = m.AccountNamed(req.Username)
	}
	modifyRequest := ldap.NewModifyRequest(
		m.GroupNamed(req.Group),
		[]ldap.Control{},
	)
	modifyRequest.Delete(m.GroupMembershipAttribute, []string{username})
	log.Debug(modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultObjectClassViolation) {
			return &RemoveLastGroupMemberError{req.Group}
		}
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) || ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchAttribute) {
			return &NoSuchMemberError{Group: req.Group, Member: req.Username}
		}
		return err
	}
	log.Infof("removed user %q from group %q", username, req.Group)
	return nil
}
