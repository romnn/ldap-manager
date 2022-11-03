package pkg

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cenkalti/backoff/v4"
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
	log.Debug(PrettyPrint(req))
	return m.ldap.Add(&req)
}

func (m *LDAPManager) setupLastGID() error {
	highestGID, err := m.GetHighestGID()
	if err != nil {
		return err
	}
	return m.setupLastID(
		highestGID,
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
		"lastUID",
		`last UID used to create a posix user,
prevents the re-use of a UID from a deleted user.`,
	)
}

func (m *LDAPManager) setupAdmin() error {
	admin := pb.NewUserRequest{
		Username:  m.DefaultAdminUsername,
		Password:  m.DefaultAdminPassword,
		FirstName: "changeme",
		LastName:  "changeme",
		Email:     "changeme@changeme.com",
	}

	// get the admin group
	adminGroup, err := m.GetGroupByName(m.DefaultAdminGroup)
	_, missing := err.(*ZeroOrMultipleGroupsError)

	if missing {
		// create admin group (admin user may not exist yet)
		strict := false
		if err := m.NewGroup(&pb.NewGroupRequest{
			Name:    m.DefaultAdminGroup,
			Members: []string{admin.GetUsername()},
		}, strict); err != nil {
			if _, exists := err.(*GroupAlreadyExistsError); !exists {
				return fmt.Errorf(
					"failed to create admin group: %v",
					err,
				)
			}
		}
	}

	if missing || len(adminGroup.Members) < 1 || m.ForceCreateAdmin {
		// add the initial admin
		if err := m.NewUser(&admin); err != nil {
			if _, exists := err.(*UserAlreadyExistsError); !exists {
				return fmt.Errorf(
					"failed to create initial admin user: %v",
					err,
				)
			}
		}

		// add the initial admin to group
		allowNonExistent := false
		if err := m.AddGroupMember(&pb.GroupMember{
			Username: admin.GetUsername(),
			Group:    m.DefaultAdminGroup,
		}, allowNonExistent); err != nil {
			if _, exists := err.(*MemberAlreadyExistsError); !exists {
				return fmt.Errorf(
					"failed to add admin user to admins group: %v",
					err,
				)
			}
		}
	}
	return nil
}

// SetupLDAP sets up the LDAP server
func (m *LDAPManager) SetupLDAP() error {
	if err := m.setupGroupOU(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if !exists {
			return fmt.Errorf(
				"failed to setup group organizational unit (OU): %v",
				err,
			)
		}
	} else {
		log.Debug("completed group organizational unit (OU) setup")
	}

	if err := m.setupUserOU(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if !exists {
			return fmt.Errorf(
				"failed to setup user organizational unit (OU): %v",
				err,
			)
		}
	} else {
		log.Debug("completed user organizational unit (OU) setup")
	}

	if err := m.setupLastGID(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		if !exists && !notFound {
			return fmt.Errorf(
				"failed to setup GID: %v",
				err,
			)
		}
	} else {
		log.Info("completed GID setup")
	}

	if err := m.setupLastUID(); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		if !exists && !notFound {
			return fmt.Errorf(
				"failed to setup UID: %v",
				err,
			)
		}
	} else {
		log.Info("completed UID setup")
	}

	if err := m.setupAdmin(); err != nil {
		return err
	}
	return nil
}

// Connect connects to the LDAP server
func (m *LDAPManager) Connect() error {
	URI := m.OpenLDAPConfig.URI()
	log.Debugf("connecting to OpenLDAP at %s", URI)

	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 30 * time.Second,
	}, 5)
	// exp := backoff.NewExponentialBackOff()
	// exp.MaxElapsedTime = 3 * time.Minute

	err := backoff.Retry(func() error {
		var err error
		m.ldap, err = ldap.DialURL(URI)
		if err != nil {
			log.Warnf("timeout dialing %s: %v", URI, err)
		}
		return err
	}, b)
	if err != nil {
		return err
	}

	// Check for TLS
	if strings.HasPrefix(URI, "ldaps:") || m.OpenLDAPConfig.TLS {
		if err := m.ldap.StartTLS(&tls.Config{
			InsecureSkipVerify: true,
		}); err != nil {
			log.Warnf("failed to connect via TLS: %v", err)
			return err
		}
	}

	// Bind as the admin user
	if err := m.BindAdmin(); err != nil {
		return err
	}
	return nil
}

// Setup sets up the LDAP server
func (m *LDAPManager) Setup() error {
	if err := m.Connect(); err != nil {
		return err
	}
	if err := m.SetupLDAP(); err != nil {
		return err
	}
	return nil
}
