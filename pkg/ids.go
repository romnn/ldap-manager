package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

type highestIDRequest struct {
	attribute      string
	min            int
	entryBaseDN    string
	entryFilter    string
	entryAttribute string
}

// GetHighestUID gets the highest UID
func (m *LDAPManager) GetHighestUID() (int, error) {
	req := highestIDRequest{
		attribute:      m.AccountAttribute,
		min:            MinUID,
		entryBaseDN:    m.UserGroupDN,
		entryFilter:    fmt.Sprintf("(%s=*)", m.AccountAttribute),
		entryAttribute: "uidNumber",
	}
	return m.getHighestID(&req)
}

// GetHighestGID gets the highest GID
func (m *LDAPManager) GetHighestGID() (int, error) {
	req := highestIDRequest{
		attribute:      m.GroupAttribute,
		min:            MinGID,
		entryBaseDN:    m.GroupsDN,
		entryFilter:    "(objectClass=posixGroup)",
		entryAttribute: "gidNumber",
	}
	return m.getHighestID(&req)
}

func (m *LDAPManager) getHighestID(req *highestIDRequest) (int, error) {
	// check for cached lastUID / lastGID value first
	attribute := strings.ToUpper(req.attribute)
	filter := fmt.Sprintf(
		"(&(objectClass=device)(cn=last%s))",
		attribute,
	)

	conn, err := m.Pool.Get()
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	result, err := conn.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"serialNumber"},
		[]ldap.Control{},
	))
	if err != nil {
		return 0, err
	}
	if len(result.Entries) > 0 {
		serial := result.Entries[0].GetAttributeValue("serialNumber")
		fetchedID, err := strconv.Atoi(serial)
		if err == nil && fetchedID >= req.min {
			return fetchedID, nil
		}
	}

	// cache miss requires traversing all entries
	result, err = conn.Search(ldap.NewSearchRequest(
		req.entryBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		req.entryFilter,
		[]string{req.entryAttribute},
		[]ldap.Control{},
	))
	if err != nil {
		return req.min, err
	}
	highestID := req.min
	for _, entry := range result.Entries {
		if id := entry.GetAttributeValue(req.entryAttribute); id != "" {
			if id, err := strconv.Atoi(id); err == nil {
				if id > highestID {
					highestID = id
				}
			}
		}
	}
	return highestID, nil
}

// updateLastID updates the id cache holding the last ID
func (m *LDAPManager) updateLastID(cn string, lastID int) error {
	req := ldap.NewModifyRequest(
		fmt.Sprintf("cn=%s,%s", cn, m.BaseDN),
		[]ldap.Control{},
	)
	req.Replace("serialNumber", []string{
		strconv.Itoa(lastID),
	})
	log.Debug(PrettyPrint(req))

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Modify(req); err != nil {
		return fmt.Errorf(
			"failed to update cn=%s: %v",
			cn, err,
		)
	}
	return nil
}
