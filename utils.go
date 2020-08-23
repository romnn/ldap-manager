package ldapmanager

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

func escape(s string) string {
	return s
}

func isErr(err error, code uint16) bool {
	return strings.HasPrefix(err.Error(), fmt.Sprintf("LDAP Result Code %d %q", code, ldap.LDAPResultCodeMap[code]))
}

func (s *LDAPManager) findGroup(groupName string, attributes []string) (*ldap.SearchResult, error) {
	return s.ldap.Search(ldap.NewSearchRequest(
		s.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", escape(groupName)),
		attributes,
		[]ldap.Control{},
	))
}

func (s *LDAPManager) updateLastID(cn string, newID int) error {
	modifyRequest := ldap.NewModifyRequest(
		fmt.Sprintf("cn=%s,%s", cn, s.BaseDN),
		[]ldap.Control{},
	)
	modifyRequest.Replace("serialNumber", []string{strconv.Itoa(newID)})
	log.Debug(modifyRequest)
	if err := s.ldap.Modify(modifyRequest); err != nil {
		return fmt.Errorf("failed to update cn=%s: %v", cn, err)
	}
	log.Debugf("updated cn=%s with %d", cn, newID)
	return nil
}
