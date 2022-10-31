package pkg

import (
	// "errors"
	"fmt"
	"strconv"
	"strings"

	"crypto/tls"
	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// BindAdmin binds as the admin user
func (m *LDAPManager) BindAdmin() error {
	adminUser := fmt.Sprintf("cn=%s,%s", "admin", m.OpenLDAPConfig.BaseDN)
	adminPassword := m.OpenLDAPConfig.AdminPassword
	return m.ldap.Bind(adminUser, adminPassword)
}

func (m *LDAPManager) setupOU(dn, ou string) error {
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

func (m *LDAPManager) setupGroupOU() error {
	return m.setupOU(m.GroupsDN, m.GroupsOU)
}

func (m *LDAPManager) setupUserOU() error {
	return m.setupOU(m.UserGroupDN, m.UsersOU)
}

func (m *LDAPManager) setupLastID(id int, cn string, desc string) error {
	req := ldap.AddRequest{
		DN: fmt.Sprintf("cn=%s,%s", cn, m.BaseDN),
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{"device", "top"}},
			{Type: "serialNumber", Vals: []string{strconv.Itoa(id)}},
			{Type: "description", Vals: []string{desc}},
		},
		Controls: []ldap.Control{},
	}
	log.Debugf("addLastIDRequest=%v", req)
	return m.ldap.Add(&req)
}

func (m *LDAPManager) setupLastGID() error {
	highestGID, err := m.GetHighestGID()
	if err != nil {
		return err
	}
	return m.setupLastID(
		highestGID,
		// m.GroupAttribute,
		"lastGID",
		`the last GID used to create a posix group,
prevents re-use of a GID from a deleted group.`,
	)
}

func (m *LDAPManager) setupLastUID() error {
	highestUID, err := m.GetHighestUID()
	if err != nil {
		return err
	}
	return m.setupLastID(
		highestUID,
		// m.AccountAttribute,
		"lastUID",
		`last UID used to create a posix user,
prevents the re-use of a UID from a deleted user.`,
	)
}

// func (m *LDAPManager) setupUserGroup() error {
// 	strict := false
// 	return m.NewGroup(&pb.NewGroupRequest{Name: m.DefaultUserGroup}, strict)
// }

func (m *LDAPManager) setupAdmin() error {
	admin := &pb.Account{
		Username:  m.DefaultAdminUsername,
		Password:  m.DefaultAdminPassword,
		FirstName: "changeme",
		LastName:  "changeme",
		Email:     "changeme@changeme.com",
	}

	// create admin group (admin user may not exist yet)
	strict := false
	if err := m.NewGroup(&pb.NewGroupRequest{
		Name:    m.DefaultAdminGroup,
		Members: []string{admin.GetUsername()},
	}, strict); err != nil {
		if _, exists := err.(*GroupAlreadyExistsError); !exists {
			return fmt.Errorf("failed to create admins group: %v", err)
		}
	}

	// get the admin group
	adminGroup, err := m.GetGroupByName(m.DefaultAdminGroup)
	if err != nil {
		return fmt.Errorf("failed to get admin group: %v", err)
	}

	if len(adminGroup.Members) < 1 || m.ForceCreateAdmin {
		// add the initial admin
		if err := m.NewUser(&pb.NewUserRequest{
			Account: admin,
		}, pb.HashingAlgorithm_DEFAULT); err != nil {
			if _, exists := err.(*UserAlreadyExistsError); !exists {
				return fmt.Errorf("failed to create initial admin account: %v", err)
			}
		}

		// add the initial admin to group
		allowNonExistent := false
		if err := m.AddGroupMember(&pb.GroupMember{
			Username: admin.GetUsername(),
			Group:    m.DefaultAdminGroup,
		}, allowNonExistent); err != nil {
			if _, exists := err.(*MemberAlreadyExistsError); !exists {
				return fmt.Errorf("failed to add initial admin user to admins group: %v", err)
			}
		}
	}
	return nil
}

// func (m *LDAPManager) setupAdminGroup() error {
// 	initialAdmin := &pb.NewUserRequest{
// 		Account: &pb.Account{
// 			Username:  m.DefaultAdminUsername,
// 			Password:  m.DefaultAdminPassword,
// 			FirstName: "changeme",
// 			LastName:  "changeme",
// 			Email:     "changeme@changeme.com",
// 		},
// 	}

// 	// Check if the group already exists
// 	adminGroup, err := m.GetGroupByName(m.DefaultAdminGroup)
// 	if err != nil {
// 		if _, ok := err.(*ZeroOrMultipleGroupsError); ok {
// 			// Create the initial admin user in the group
// 			if err := m.NewUser(initialAdmin, pb.HashingAlgorithm_DEFAULT); err != nil {
// 				if _, exists := err.(*UserAlreadyExistsError); !exists {
// 					return fmt.Errorf("failed to create initial admin account: %v", err)
// 				}
// 			}
// 			// Create the group
// 			strict := false
// 			if err := m.NewGroup(&pb.NewGroupRequest{
// 				Name:    m.DefaultAdminGroup,
// 				Members: []string{initialAdmin.GetAccount().GetUsername()},
// 			}, strict); err != nil {
// 				return fmt.Errorf("failed to create admins group: %v", err)
// 			}
// 			return nil
// 		}
// 		return fmt.Errorf("failed to check if the admins group already exists: %v", err)
// 	}

// 	// Group already exists
// 	if len(adminGroup.Members) < 1 || m.ForceCreateAdmin {
// 		// Add the default admin user
// 		if err := m.NewUser(initialAdmin, pb.HashingAlgorithm_DEFAULT); err != nil {
// 			if _, exists := err.(*UserAlreadyExistsError); !exists {
// 				return fmt.Errorf("failed to create initial admin account: %v", err)
// 			}
// 		}
// 		allowNonExistent := false
// 		if err := m.AddGroupMember(&pb.GroupMember{
// 			Username: initialAdmin.GetAccount().GetUsername(),
// 			Group:    m.DefaultAdminGroup,
// 		}, allowNonExistent); err != nil {
// 			if _, ok := err.(*MemberAlreadyExistsError); !ok {
// 				return fmt.Errorf("failed to add the default admin user to the admins group: %v", err)
// 			}
// 		}
// 	}
// 	return nil
// }

// func (m *LDAPManager) setupDefaultAdmin() error {
// 	// Check if there are already admins
// 	adminGroup, err := m.GetGroupByName(m.DefaultAdminGroup)
// 	if err != nil {
// 		return err
// 	}
// 	if len(adminGroup.Members) < 1 {
// 		return errors.New("no admin user created")
// 	}
// 	return nil
// }

// SetupLDAP ...
func (m *LDAPManager) SetupLDAP() error {
	if err := m.setupGroupOU(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if !exists {
			return fmt.Errorf("failed to setup group organizational unit (OU): %v", err)
		}
	} else {
		log.Debug("completed setup of group organizational unit")
	}

	if err := m.setupUserOU(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if !exists {
			return fmt.Errorf("failed to setup user organizational unit (OU): %v", err)
		}
	} else {
		log.Debug("completed setup of user organizational unit")
	}

	if err := m.setupLastGID(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		if !exists && !notFound {
			return fmt.Errorf("failed to setup the last GID: %v", err)
		}
	} else {
		log.Info("completed setup of the last GID")
	}

	if err := m.setupLastUID(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		if !exists && !notFound {
			return fmt.Errorf("failed to setup the last UID: %v", err)
		}
	} else {
		log.Info("completed setup of the last UID")
	}

	if err := m.setupAdmin(); err != nil {
		return err
	}
	return nil
}

// Setup ...
// func (m *LDAPManager) Setup(skipSetupLDAP bool) error {
func (m *LDAPManager) Setup() error {
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
	// if !skipSetupLDAP {
	if err := m.SetupLDAP(); err != nil {
		return err
	}
	// }
	return nil
}
