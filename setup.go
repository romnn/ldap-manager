package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap"
	log "github.com/sirupsen/logrus"
)

// BindReadOnly ...
func (s *LDAPManagerServer) BindReadOnly() error {
	return s.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", s.ReadonlyBindUsername), s.ReadonlyBindPassword)
}

// BindAdmin ...
func (s *LDAPManagerServer) BindAdmin() error {
	return s.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", s.AdminBindUsername), s.AdminBindPassword)
}

func (s *LDAPManagerServer) setupOU(dn, name string) error {
	addOURequest := &ldap.AddRequest{
		DN: dn,
		Attributes: []ldap.Attribute{
			{"objectClass", []string{"organizationalUnit"}},
			{"ou", []string{name}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addOURequest)
	return s.ldap.Add(addOURequest)
}

func (s *LDAPManagerServer) setupGroupsOU() error {
	return s.setupOU(s.GroupsDN, s.GroupsOU)
}

func (s *LDAPManagerServer) setupUsersOU() error {
	return s.setupOU(s.UserGroupDN, s.UsersOU)
}

func (s *LDAPManagerServer) setupLastID(attribute, cn string, desc string) error {
	highestID, err := s.GetHighestID(attribute)
	if err != nil {
		return err
	}
	addLastIDRequest := &ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", cn, s.BaseDN),
		Attributes: []ldap.Attribute{
			{"objectClass", []string{"device", "top"}},
			{"serialnumber", []string{strconv.Itoa(highestID)}},
			{"description", []string{desc}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(addLastIDRequest)
	return s.ldap.Add(addLastIDRequest)
}

func (s *LDAPManagerServer) setupLastGID() error {
	return s.setupLastID(
		s.GroupAttribute, "lastGID",
		"Records the last GID used to create a Posix group. This prevents the re-use of a GID from a deleted group.",
	)
}

func (s *LDAPManagerServer) setupLastUID() error {
	return s.setupLastID(
		s.AccountAttribute, "lastUID",
		"Records the last UID used to create a Posix account. This prevents the re-use of a UID from a deleted account.",
	)
}

func (s *LDAPManagerServer) setupDefaultGroup() error {
	return s.NewGroup(s.DefaultUserGroup, []string{})
}

func (s *LDAPManagerServer) setupAdminsGroup() error {
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

func (s *LDAPManagerServer) setupAuth(adminPassword string) error {
	return s.ldap.Bind(fmt.Sprintf("cn=%s,dc=example,dc=org", s.AdminBindUsername), adminPassword)
}

// SetupLDAP ...
func (s *LDAPManagerServer) SetupLDAP() error {
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
