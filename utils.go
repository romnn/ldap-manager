package ldapmanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap"
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

// ListOptions ...
type ListOptions struct {
	Start     int    `json:"start" form:"start"`
	End       int    `json:"end" form:"end"`
	SortOrder string `json:"sort_order" form:"sort_order"`
	SortKey   string `json:"sort_key" form:"sort_key"`
}

func escape(s string) string {
	return ldap.EscapeFilter(s)
}

func isErr(err error, code uint16) bool {
	return strings.HasPrefix(err.Error(), fmt.Sprintf("LDAP Result Code %d %q", code, ldap.LDAPResultCodeMap[code]))
}

func (m *LDAPManager) findGroup(groupName string, attributes []string) (*ldap.SearchResult, error) {
	return m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escape(groupName)),
		attributes,
		[]ldap.Control{},
	))
}

func (m *LDAPManager) updateLastID(cn string, newID int) error {
	modifyRequest := ldap.NewModifyRequest(
		fmt.Sprintf("cn=%s,%s", cn, m.BaseDN),
		[]ldap.Control{},
	)
	modifyRequest.Replace("serialNumber", []string{strconv.Itoa(newID)})
	log.Debug(modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		return fmt.Errorf("failed to update cn=%s: %v", cn, err)
	}
	log.Debugf("updated cn=%s with %d", cn, newID)
	return nil
}

func (m *LDAPManager) getHighestID(attribute string) (int, error) {
	var highestID int
	var entryBaseDN, entryFilter, entryAttribute string
	switch strings.ToUpper(attribute) {
	case strings.ToUpper(m.GroupAttribute):
		highestID = MinGID
		entryBaseDN = m.GroupsDN
		entryFilter = "(objectClass=posixGroup)"
		entryAttribute = "gidNumber"
	case strings.ToUpper(m.AccountAttribute):
		highestID = MinUID
		entryBaseDN = m.UserGroupDN
		entryFilter = fmt.Sprintf("(%s=*)", m.AccountAttribute)
		entryAttribute = "uidNumber"
	default:
		return highestID, fmt.Errorf("unknown id attribute %q", attribute)
	}

	filter := fmt.Sprintf("(&(objectClass=device)(cn=last%s))", strings.ToUpper(attribute))
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.BaseDN,
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
	result, err = m.ldap.Search(ldap.NewSearchRequest(
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
