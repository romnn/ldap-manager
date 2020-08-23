package ldapmanager

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// BindReadOnly ...
func (s *LDAPManager) BindReadOnly() error {
	return s.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", s.ReadonlyBindUsername), s.ReadonlyBindPassword)
}

// BindAdmin ...
func (s *LDAPManager) BindAdmin() error {
	return s.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", s.AdminBindUsername), s.AdminBindPassword)
}

func (s *LDAPManager) setupOU(dn, name string) error {
	addOURequest := &ldap.AddRequest{
		DN: dn,
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"organizationalUnit"}},
			{Type: "ou", Vals: []string{name}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addOURequest)
	return s.ldap.Add(addOURequest)
}

func (s *LDAPManager) setupGroupsOU() error {
	return s.setupOU(s.GroupsDN, s.GroupsOU)
}

func (s *LDAPManager) setupUsersOU() error {
	return s.setupOU(s.UserGroupDN, s.UsersOU)
}

func (s *LDAPManager) setupLastID(attribute, cn string, desc string) error {
	highestID, err := s.GetHighestID(attribute)
	if err != nil {
		return err
	}
	addLastIDRequest := &ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", cn, s.BaseDN),
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"device", "top"}},
			{Type: "serialnumber", Vals: []string{strconv.Itoa(highestID)}},
			{Type: "description", Vals: []string{desc}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addLastIDRequest)
	return s.ldap.Add(addLastIDRequest)
}

func (s *LDAPManager) setupLastGID() error {
	return s.setupLastID(
		s.GroupAttribute, "lastGID",
		"Records the last GID used to create a Posix group. This prevents the re-use of a GID from a deleted group.",
	)
}

func (s *LDAPManager) setupLastUID() error {
	return s.setupLastID(
		s.AccountAttribute, "lastUID",
		"Records the last UID used to create a Posix account. This prevents the re-use of a UID from a deleted account.",
	)
}

func (s *LDAPManager) setupDefaultGroup() error {
	return s.NewGroup(s.DefaultUserGroup, []string{})
}

func (s *LDAPManager) setupAdminsGroup() error {
	if err := s.NewGroup(s.DefaultAdminGroup, []string{}); err != nil {
		return err
	}
	admins, err := s.GetGroupMembers(s.DefaultAdminGroup, 0, 0, "")
	if err != nil {
		return err
	}
	if len(admins) < 1 {
		return errors.New("no admin user created")
	}
	return nil
}

func (s *LDAPManager) setupAuth(adminPassword string) error {
	return s.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", s.AdminBindUsername), adminPassword)
}

// SetupLDAP ...
func (s *LDAPManager) SetupLDAP() error {
	if err := s.setupGroupsOU(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
		return fmt.Errorf("failed to setup groups organizational unit (OU): %v", err)
	}
	if err := s.setupUsersOU(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
		return fmt.Errorf("failed to setup users organizational unit (OU): %v", err)
	}
	if err := s.setupLastGID(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
		return fmt.Errorf("failed to setup the last GID: %v", err)
	}
	if err := s.setupLastUID(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
		return fmt.Errorf("failed to setup the last UID: %v", err)
	}
	/*
		if err := s.setupDefaultGroup(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup the default user group: %v", err)
		}
		if err := s.setupAdminsGroup(); err != nil && !isErr(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup the default admin group: %v", err)
		}
	*/
	return nil
}
