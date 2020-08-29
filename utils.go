package ldapmanager

import (
	"fmt"
	"regexp"
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

	hex = "0123456789abcdef"
)

// ListOptions ...
type ListOptions struct {
	Start     int    `json:"start" form:"start"`
	End       int    `json:"end" form:"end"`
	SortOrder string `json:"sort_order" form:"sort_order"`
	SortKey   string `json:"sort_key" form:"sort_key"`
}

func isValidAttribute(attr string) bool {
	switch attr {
	case "uid":
		return true
	case "cn":
		return true
	case "uidNumber":
		return true
	case "gidNumber":
		return true
	case "mail":
		return true
	case "sn":
		return true
	case "givenName":
		return true
	case "displayName":
		return true
	case "loginShell":
		return true
	case "homeDirectory":
		return true
	}
	return false
}

func escapeFilter(s string) string {
	return ldap.EscapeFilter(s)
}

func mustescapeFilter(c byte) bool {
	if c > 0x7f {
		return true
	}
	switch c {
	case ',', ';', '(', ')', '\\', '*', '"', '#', '=', '+', '<', '>', 0:
		return true
	}
	return false
}

func escapeDN(dn string) string {
	// escapes https://ldapwiki.com/wiki/DN%20Escape%20Values
	escape := 0
	for i := 0; i < len(dn); i++ {
		if mustescapeFilter(dn[i]) {
			escape++
		}
	}
	if escape == 0 {
		return dn
	}
	buf := make([]byte, len(dn)+escape*2)
	for i, j := 0, 0; i < len(dn); i++ {
		c := dn[i]
		if mustescapeFilter(c) {
			buf[j+0] = '\\'
			buf[j+1] = hex[c>>4]
			buf[j+2] = hex[c&0xf]
			j += 3
		} else {
			buf[j] = c
			j++
		}
	}
	return string(buf)
}

func parseFilter(filters []string) string {
	var filter string
	for _, f := range filters {
		if pair := strings.Split(f, "="); len(pair) == 2 {
			if attr := strings.ToLower(pair[0]); isValidAttribute(attr) {
				filter += fmt.Sprintf("(%s=*%s*)", attr, escapeFilter(pair[1]))
			}
		}
	}
	return filter
}

func extractAttribute(dn string, attribute string) (string, error) {
	reg, err := regexp.Compile(fmt.Sprintf("%s=(?P<Attribute>.*?),", attribute))
	if err != nil {
		return "", err
	}
	if matches := reg.FindStringSubmatch(dn); len(matches) > 1 {
		if match := matches[1]; match != "" {
			return match, nil
		}
	}
	return "", fmt.Errorf("could not find attribute %q in %q", attribute, dn)
}

func (m *LDAPManager) findGroup(groupName string, attributes []string) (*ldap.SearchResult, error) {
	return m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escapeFilter(groupName)),
		attributes,
		[]ldap.Control{},
	))
}

func (m *LDAPManager) getGroupByGID(gid int) (string, int, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(gid=%d)", gid),
		[]string{"cn"},
		[]ldap.Control{},
	))
	if err != nil {
		return "", 0, err
	}
	if len(result.Entries) != 1 {
		return "", 0, fmt.Errorf("zero or multiple groups with gid=%d", gid)
	}
	group := result.Entries[0]
	cn := group.GetAttributeValue("cn")
	if cn == "" {
		return "", 0, fmt.Errorf("group with gid=%d has no valid cn attribute", gid)
	}
	return cn, gid, nil
}

func (m *LDAPManager) updateLastID(cn string, newID int) error {
	modifyRequest := ldap.NewModifyRequest(
		fmt.Sprintf("cn=%s,%s", cn, m.BaseDN),
		[]ldap.Control{},
	)
	modifyRequest.Replace("serialNumber", []string{strconv.Itoa(newID)})
	log.Debugf("modifyRequest=%v", modifyRequest)
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
