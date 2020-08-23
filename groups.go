package ldapmanager

import (
	"errors"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// GroupExistsError ...
type GroupExistsError struct {
	Group string
}

// GroupExistsError ...
func (e *GroupExistsError) Error() string {
	return fmt.Sprintf("group %q already exists", e.Group)
}

// NewGroup ...
func (m *LDAPManager) NewGroup(name string, members []string) error {
	if name == "" {
		return errors.New("group name can not be empty")
	}
	result, err := m.findGroup(name, []string{"dn", m.GroupMembershipAttribute})
	if err != nil {
		return err
	}
	if len(result.Entries) > 0 {
		return &GroupExistsError{Group: name}
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
			{Type: "cn", Vals: []string{name}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(newGID)}},
		}
	} else {
		if len(members) < 1 {
			return errors.New("when using RFC2307BIS (not NIS), you must specify at least one group member")
		}
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"top", "groupOfUniqueNames", "posixGroup"}},
			{Type: "cn", Vals: []string{name}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(newGID)}},
			{Type: m.GroupMembershipAttribute, Vals: members},
		}
	}
	addGroupRequest := &ldap.AddRequest{
		DN:         fmt.Sprintf("cn=%s,%s", name, m.GroupsDN),
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
	log.Infof("added new group %q (gid=%d)", name, newGID)
	return nil
}

// DeleteGroup ...
func (m *LDAPManager) DeleteGroup(groupName string) error {
	if err := m.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("cn=%s,%s", escape(groupName), m.GroupsDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	log.Infof("removed group %q", groupName)
	return nil
}

// GetGroupList ...
func (m *LDAPManager) GetGroupList(start, end int, sortOrder string, filters []string) ([]string, error) {
	filter := fmt.Sprintf("(&(objectClass=*)%s)", filters)
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

// GetGroupGID ...
func (m *LDAPManager) GetGroupGID(groupName string) (int, error) {
	result, err := m.findGroup(groupName, []string{"gidNumber"})
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
