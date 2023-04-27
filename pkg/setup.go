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
	ldappool "github.com/romnn/ldap-manager/pkg/pool"
	log "github.com/sirupsen/logrus"
)

// AdminUserDN gets the DN of the admin user
func (m *LDAPManager) AdminUserDN() string {
	return fmt.Sprintf(
		"cn=%s,%s", m.Config.AdminUsername,
		m.Config.BaseDN,
	)
}

// ReadOnlyUserDN gets the DN of the read-only user
func (m *LDAPManager) ReadOnlyUserDN() string {
	return fmt.Sprintf(
		"cn=%s,%s",
		m.Config.ReadOnlyUsername,
		m.Config.BaseDN,
	)
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
	log.Debug(PrettyPrint(addOURequest))
	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Add(addOURequest)
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
	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Add(&req)
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

func (m *LDAPManager) setupReadOnlyUser() error {
	// see https://github.com/osixia/docker-openldap/tree/master/image/service/slapd/assets/config/bootstrap/ldif/readonly-user

	username := m.Config.ReadOnlyUsername
	addUserReq := &ldap.AddRequest{
		DN: m.ReadOnlyUserDN(),
		Attributes: []ldap.Attribute{
			{Type: "objectClass", Vals: []string{
				"simpleSecurityObject",
				"organizationalRole",
			}},
			{Type: "cn", Vals: []string{username}},
			{Type: "userPassword", Vals: []string{"placeholder"}},
			{Type: "description", Vals: []string{"LDAP read only user"}},
		},
		Controls: []ldap.Control{},
	}
	log.Debug(PrettyPrint(addUserReq))

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Add(addUserReq); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if exists {
			return &UserAlreadyExistsError{
				Username: username,
			}
		}
		return fmt.Errorf(
			"failed to add user %q: %v",
			username, err,
		)
	}

	// bind for the config CN to apply ACL rules
	configDN := fmt.Sprintf(
		"cn=%s,cn=config",
		m.Config.AdminUsername,
	)
	configPassword := "config"
	if err := conn.Bind(configDN, configPassword); err != nil {
		return fmt.Errorf(
			"unable to bind as %q with password %q: %v",
			configDN, configPassword, err,
		)
	}

	ldapBackend := "mdb"
	aclReq := ldap.NewModifyRequest(
		fmt.Sprintf(
			// "olcDatabase={0}config,cn=config",
			"olcDatabase={1}%s,cn=config",
			ldapBackend,
		),
		[]ldap.Control{},
	)
	aclReq.Add("olcAccess", []string{
		`to * by dn.exact=gidNumber=0+uidNumber=0,cn=peercred,cn=external,cn=auth manage by * break`,
		fmt.Sprintf(
			`to attrs=userPassword,shadowLastChange by self write by dn="%s" write by anonymous auth by * none`,
			m.AdminUserDN(),
		),
		fmt.Sprintf(
			`to * by self read by dn="%s" write by dn="%s" read by * none`,
			m.AdminUserDN(),
			m.ReadOnlyUserDN(),
		),
	})

	log.Debug(PrettyPrint(aclReq))
	if err := conn.Modify(aclReq); err != nil {
		exists := ldap.IsErrorWithCode(err, ldap.LDAPResultEntryAlreadyExists)
		if exists {
			return &UserAlreadyExistsError{
				Username: username,
			}
		}
		return fmt.Errorf(
			"failed to add ACL rules for %q: %v",
			username, err,
		)
	}

	return nil
}

func (m *LDAPManager) setupAdmin() error {
	// get the admin group (if already exists)
	adminGroup, err := m.GetGroupByName(m.DefaultAdminGroup)

	var presentAdmins []string
	if err == nil {
		presentAdmins = adminGroup.Members
	}

	// if 1 or more admins in the admin group exist, we cannot
	// assume their credentials are still the same unless forced..
	if !m.ForceCreateAdmin && len(presentAdmins) > 0 {
		log.Infof(
			"found existing admins %v: skip create default admin",
			presentAdmins,
		)
		return nil
	}

	// IMPORTANT: create the admin user before the groups
	// otherwise, the memberOf overlay will never pick up that
	// the admin user belongs to the users and admins groups
	// see: https://github.com/osixia/docker-openldap/issues/635
	admin := pb.NewUserRequest{
		Username:  m.DefaultAdminUsername,
		Password:  m.DefaultAdminPassword,
		FirstName: "changeme",
		LastName:  "changeme",
		Email:     "changeme@changeme.com",
	}
	log.Infof("creating default admin %q", admin.GetUsername())
	if err := m.NewUser(&admin); err != nil {
		if _, exists := err.(*UserAlreadyExistsError); !exists {
			return fmt.Errorf(
				"failed to create initial admin user: %v",
				err,
			)
		}
	}

	// create initial groups and add admin user to them
	for _, groupName := range []string{
		m.DefaultAdminGroup,
		m.DefaultUserGroup,
	} {
		strict := false
		if err := m.NewGroup(&pb.NewGroupRequest{
			Name:    groupName,
			Members: []string{admin.GetUsername()},
		}, strict); err != nil {
			if _, exists := err.(*GroupAlreadyExistsError); !exists {
				return fmt.Errorf(
					"failed to create %q group: %v",
					groupName, err,
				)
			}
		}

		allowNonExistent := false
		if err := m.AddGroupMember(&pb.GroupMember{
			Username: admin.GetUsername(),
			Group:    groupName,
		}, allowNonExistent); err != nil {
			if _, exists := err.(*MemberAlreadyExistsError); !exists {
				return fmt.Errorf(
					"failed to add admin user %q to group %q: %v",
					admin.GetUsername(), groupName, err,
				)
			}
		}
	}

	// make sure default admin has admin status
	memberStatus, err := m.IsGroupMember(
		&pb.IsGroupMemberRequest{
			Username: admin.GetUsername(),
			Group:    m.DefaultAdminGroup,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"failed to check admin status for default admin %q: %v",
			admin.GetUsername(), err,
		)
	}
	if !memberStatus.GetIsMember() {
		return fmt.Errorf(
			"default admin %q does not have admin privileges",
			admin.GetUsername(),
		)
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
	// if m.Config.ReadOnlyUser {
	// 	if err := m.SetupReadOnlyUser(); err != nil {
	// 		return err
	// 	}
	// }
	return nil
}

// Connect sets up the connection pool to the LDAP server
func (m *LDAPManager) Connect() error {
	var err error

	// factory for creating new connections
	factory := func() (ldap.Client, error) {
		URI := m.Config.URI()
		b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
			Interval: 2 * time.Second,
		}, 10)

		var conn *ldap.Conn
		err := backoff.Retry(func() error {
			var err error
			conn, err = ldap.DialURL(URI)
			if err != nil {
				log.Warnf("timeout dialing %s: %v", URI, err)
			}
			return err
		}, b)
		if err != nil {
			return nil, err
		}

		// check for TLS
		if strings.HasPrefix(URI, "ldaps:") || m.Config.TLS {
			if err := conn.StartTLS(&tls.Config{
				InsecureSkipVerify: true,
			}); err != nil {
				log.Warnf("failed to connect via TLS: %v", err)
				return nil, err
			}
		}
		return conn, nil
	}

	reset := func(conn ldap.Client) error {
		// re-bind as the admin user
		return conn.Bind(
			m.AdminUserDN(),
			m.Config.AdminPassword,
		)
	}

	m.Pool, err = ldappool.NewChannelPool(10, 20, factory, reset)
	return err
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
