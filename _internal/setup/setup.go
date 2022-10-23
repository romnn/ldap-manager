package setup

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// BindAdmin ...
func (m *ldapmanager.LDAPManager) BindAdmin() error {
	return m.ldap.Bind(fmt.Sprintf("cn=%s,%s", "admin", m.OpenLDAPConfig.BaseDN), m.OpenLDAPConfig.AdminPassword)
}

func (m *ldapmanager.LDAPManager) setupOU(dn, ou string) error {
	addOURequest := &ldap.AddRequest{
		DN: dn,
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"organizationalUnit"}},
			{Type: "ou", Vals: []string{ou}},
		},
		Controls: []ldap.Control{},
	}
	log.Debugf("addOURequest=%v", addOURequest)
	return m.ldap.Add(addOURequest)
}

func (m *ldapmanager.LDAPManager) setupGroupsOU() error {
	return m.setupOU(m.GroupsDN, m.GroupsOU)
}

func (m *ldapmanager.LDAPManager) setupUsersOU() error {
	return m.setupOU(m.UserGroupDN, m.UsersOU)
}

func (m *ldapmanager.LDAPManager) setupLastID(attribute, cn string, desc string) error {
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
	log.Debugf("addLastIDRequest=%v", addLastIDRequest)
	return m.ldap.Add(addLastIDRequest)
}

func (m *ldapmanager.LDAPManager) setupLastGID() error {
	return m.setupLastID(
		m.GroupAttribute, "lastGID",
		"Records the last GID used to create a Posix group. This prevents the re-use of a GID from a deleted group.",
	)
}

func (m *ldapmanager.LDAPManager) setupLastUID() error {
	return m.setupLastID(
		m.AccountAttribute, "lastUID",
		"Records the last UID used to create a Posix account. This prevents the re-use of a UID from a deleted account.",
	)
}

func (m *ldapmanager.LDAPManager) setupDefaultGroup() error {
	strict := false
	return m.NewGroup(&pb.NewGroupRequest{Name: m.DefaultUserGroup}, strict)
}

func (m *ldapmanager.LDAPManager) setupAdminsGroup() error {
	initialAdmin := &pb.NewAccountRequest{
		Account: &pb.Account{
			Username:  m.DefaultAdminUsername,
			Password:  m.DefaultAdminPassword,
			FirstName: "changeme",
			LastName:  "changeme",
			Email:     "changeme@changeme.com",
		},
	}

	// Check if the group already exists
	adminGroup, err := m.GetGroup(&pb.GetGroupRequest{Name: m.DefaultAdminGroup})
	if err != nil {
		if _, ok := err.(*ZeroOrMultipleGroupsError); ok {
			// Create the initial admin user in the group
			if err := m.NewAccount(initialAdmin, pb.HashingAlgorithm_DEFAULT); err != nil {
				if _, ok := err.(*AccountAlreadyExistsError); !ok {
					return fmt.Errorf("failed to create initial admin account: %v", err)
				}
			}
			// Create the group
			strict := false
			if err := m.NewGroup(&pb.NewGroupRequest{Name: m.DefaultAdminGroup, Members: []string{initialAdmin.GetAccount().GetUsername()}}, strict); err != nil {
				return fmt.Errorf("failed to create admins group: %v", err)
			}
			return nil
		}
		return fmt.Errorf("failed to check if the admins group already exists: %v", err)
	}

	// Group already exists
	if len(adminGroup.Members) < 1 || m.ForceCreateAdmin {
		// Add the default admin user
		if err := m.NewAccount(initialAdmin, pb.HashingAlgorithm_DEFAULT); err != nil {
			if _, ok := err.(*AccountAlreadyExistsError); !ok {
				return fmt.Errorf("failed to create initial admin account: %v", err)
			}
		}
		allowNonExistent := false
		if err := m.AddGroupMember(&pb.GroupMember{Username: initialAdmin.GetAccount().GetUsername(), Group: m.DefaultAdminGroup}, allowNonExistent); err != nil {
			if _, ok := err.(*MemberAlreadyExistsError); !ok {
				return fmt.Errorf("failed to add the default admin user to the admins group: %v", err)
			}
		}
	}
	return nil
}

func (m *ldapmanager.LDAPManager) setupDefaultAdmin() error {
	// Check if there are already admins
	adminGroup, err := m.GetGroup(&pb.GetGroupRequest{Name: m.DefaultAdminGroup})
	if err != nil {
		return err
	}
	if len(adminGroup.Members) < 1 {
		return errors.New("no admin user created")
	}
	return nil
}

// SetupLDAP ...
func (m *ldapmanager.LDAPManager) SetupLDAP() error {
	if err := m.setupGroupsOU(); err != nil {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup groups organizational unit (OU): %v", err)
		}
	} else {
		log.Debug("completed setup of groups organizational unit")
	}

	if err := m.setupUsersOU(); err != nil {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) {
			return fmt.Errorf("failed to setup users organizational unit (OU): %v", err)
		}
	} else {
		log.Debug("completed setup of users organizational unit")
	}

	if err := m.setupLastGID(); err != nil {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) && !ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			return fmt.Errorf("failed to setup the last GID: %v", err)
		}
	} else {
		log.Debug("completed setup of the last GID")
	}

	if err := m.setupLastUID(); err != nil {
		if !ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists) && !ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			return fmt.Errorf("failed to setup the last UID: %v", err)
		}
	} else {
		log.Debug("completed setup of the last UID")
	}

	if err := m.setupAdminsGroup(); err != nil {
		return err
	}
	return nil
}

// Setup ...
func (m *LDAPManager) Setup(skipSetupLDAP bool) error {
	var err error
	URI := m.OpenLDAPConfig.URI()
	log.Debugf("connecting to OpenLDAP at %s", URI)
	m.ldap, err = ldap.DialURL(URI)
	if err != nil {
		return err
	}

	// Check for TLS
	if strings.HasPrefix(URI, "ldaps:") || m.OpenLDAPConfig.TLS {
		if err := m.ldap.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			log.Warnf("failed to connect via TLS: %v", err)
			if m.OpenLDAPConfig.TLS {
				return err
			}
		}
	}

	// Bind as the admin user
	if err := m.BindAdmin(); err != nil {
		return err
	}
	if !skipSetupLDAP {
		if err := m.SetupLDAP(); err != nil {
			return err
		}
	}
	return nil
}
