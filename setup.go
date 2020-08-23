package ldapmanager

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// BindReadOnly ...
func (m *LDAPManager) BindReadOnly() error {
	return m.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", m.OpenLDAPConfig.ReadonlyUserUsername), m.OpenLDAPConfig.ReadonlyUserPassword)
}

// BindAdmin ...
func (m *LDAPManager) BindAdmin() error {
	return m.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", "admin"), m.OpenLDAPConfig.AdminPassword)
}

func (m *LDAPManager) setupOU(dn, name string) error {
	addOURequest := &ldap.AddRequest{
		DN: dn,
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"organizationalUnit"}},
			{Type: "ou", Vals: []string{name}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addOURequest)
	return m.ldap.Add(addOURequest)
}

func (m *LDAPManager) setupGroupsOU() error {
	return m.setupOU(m.GroupsDN, m.GroupsOU)
}

func (m *LDAPManager) setupUsersOU() error {
	return m.setupOU(m.UserGroupDN, m.UsersOU)
}

func (m *LDAPManager) setupLastID(attribute, cn string, desc string) error {
	highestID, err := m.getHighestID(attribute)
	if err != nil {
		return err
	}
	addLastIDRequest := &ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", cn, m.BaseDN),
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"device", "top"}},
			{Type: "serialNumber", Vals: []string{strconv.Itoa(highestID)}},
			{Type: "description", Vals: []string{desc}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addLastIDRequest)
	return m.ldap.Add(addLastIDRequest)
}

func (m *LDAPManager) setupLastGID() error {
	return m.setupLastID(
		m.GroupAttribute, "lastGID",
		"Records the last GID used to create a Posix group. This prevents the re-use of a GID from a deleted group.",
	)
}

func (m *LDAPManager) setupLastUID() error {
	return m.setupLastID(
		m.AccountAttribute, "lastUID",
		"Records the last UID used to create a Posix account. This prevents the re-use of a UID from a deleted account.",
	)
}

func (m *LDAPManager) setupDefaultGroup() error {
	return m.NewGroup(m.DefaultUserGroup, []string{})
}

func (m *LDAPManager) setupAdminsGroup() error {
	if err := m.NewGroup(m.DefaultAdminGroup, []string{}); err != nil {
		return err
	}
	admins, err := m.GetGroupMembers(m.DefaultAdminGroup, 0, 0, "")
	if err != nil {
		return err
	}
	if len(admins) < 1 {
		return errors.New("no admin user created")
	}
	return nil
}

func (m *LDAPManager) setupAuth(adminPassword string) error {
	return m.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", "admin"), adminPassword)
}

// SetupLDAP ...
func (m *LDAPManager) SetupLDAP() error {
	if err := m.setupGroupsOU(); err != nil {
		if !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup groups organizational unit (OU): %v", err)
		}
	} else {
		log.Debug("completed setup of groups organizational unit")
	}

	if err := m.setupUsersOU(); err != nil {
		if !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup users organizational unit (OU): %v", err)
		}
	} else {
		log.Debug("completed setup of users organizational unit")
	}

	if err := m.setupLastGID(); err != nil {
		if !isErr(err, ldap.LDAPResultEntryAlreadyExists) && !isErr(err, ldap.LDAPResultNoSuchObject) {
			return fmt.Errorf("failed to setup the last GID: %v", err)
		}
	} else {
		log.Debug("completed setup of the last GID")
	}

	if err := m.setupLastUID(); err != nil {
		if !isErr(err, ldap.LDAPResultEntryAlreadyExists) && !isErr(err, ldap.LDAPResultNoSuchObject) {
			return fmt.Errorf("failed to setup the last UID: %v", err)
		}
	} else {
		log.Debug("completed setup of the last UID")
	}

	// Unfortunately, we cannot setup groups here without initial members
	/*
		if err := m.setupDefaultGroup(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup the default user group: %v", err)
		}
		if err := m.setupAdminsGroup(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup the default admin group: %v", err)
		}
	*/
	return nil
}
