package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	// GroupsOu configures the LDAP group organizational unit
	GroupsOu = cli.StringFlag{
		Name:    "groups-ou",
		Value:   "groups",
		EnvVars: []string{"GROUPS_OU"},
		Usage:   "group organizational unit",
	}
	// UsersOu configures the LDAP user organizational unit
	UsersOu = cli.StringFlag{
		Name:    "users-ou",
		Value:   "users",
		EnvVars: []string{"USERS_OU"},
		Usage:   "user organizational unit",
	}
	// GroupsDn configures the LDAP groups DN
	GroupsDn = cli.StringFlag{
		Name:    "groups-dn",
		Value:   "", // default is ou=GROUPS_OU,BASE_DN
		EnvVars: []string{"GROUPS_DN"},
		Usage:   "groups DN (default is ou=$GROUPS_OU,$BASE_DN)",
	}
	// UsersDn configures the LDAP users DN
	UsersDn = cli.StringFlag{
		Name:    "users-dn",
		Value:   "", // default is ou=USERS_DN,BASE_DN
		EnvVars: []string{"USERS_DN"},
		Usage:   "users DN (default is ou=$USERS_DN,$BASE_DN)",
	}
	// GroupMembershipAttribute configures the LDAP group membership attribute
	GroupMembershipAttribute = cli.GenericFlag{
		Name: "group-membership-attribute",
		Value: &EnumValue{
			Enum:    []string{"uniqueMember", "memberUID"},
			Default: "uniqueMember",
		},
		EnvVars: []string{"GROUP_MEMBERSHIP_ATTRIBUTE"},
		Usage:   "group membership attribute (e.g. uniqueMember)",
	}
	// GroupMembershipUsesUID configures if LDAP uses UID for group membership
	GroupMembershipUsesUID = cli.BoolFlag{
		Name:    "group-membership-uses-uid",
		Value:   false,
		EnvVars: []string{"GROUP_MEMBERSHIP_USES_UID"},
		Usage:   "group membership uses UID only instead of full DN",
	}
	// AccountAttribute configures the LDAP account attribute
	AccountAttribute = cli.StringFlag{
		Name:    "account-attribute",
		Value:   "uid",
		EnvVars: []string{"ACCOUNT_ATTRIBUTE"},
		Usage:   "account attribute",
	}
	// GroupAttribute configures the LDAP group attribute
	GroupAttribute = cli.StringFlag{
		Name:    "group-attribute",
		Value:   "gid",
		EnvVars: []string{"GROUP_ATTRIBUTE"},
		Usage:   "group attribute",
	}
	// DefaultUserGroup configures the default LDAP user group
	DefaultUserGroup = cli.StringFlag{
		Name:    "default-user-group",
		Value:   "users",
		EnvVars: []string{"DEFAULT_USER_GROUP"},
		Usage:   "default user group",
	}
	// DefaultAdminGroup configures the default LDAP admin group
	DefaultAdminGroup = cli.StringFlag{
		Name:    "default-admin-group",
		Value:   "admins",
		EnvVars: []string{"DEFAULT_ADMIN_GROUP"},
		Usage:   "default admin group",
	}
	// DefaultLoginShell configures the default LDAP login shell
	DefaultLoginShell = cli.StringFlag{
		Name:    "default-login-shell",
		Value:   "/bin/bash",
		EnvVars: []string{"DEFAULT_LOGIN_SHELL"},
		Usage:   "default login shell",
	}
	// DefaultAdminUsername configures the default LDAP admin username
	DefaultAdminUsername = cli.StringFlag{
		Name:    "default-admin-username",
		Value:   "admin",
		EnvVars: []string{"DEFAULT_ADMIN_USERNAME"},
		Usage:   "default admin username",
	}
	// DefaultAdminPassword configures the default LDAP admin password
	DefaultAdminPassword = cli.StringFlag{
		Name:    "default-admin-password",
		Value:   "admin",
		EnvVars: []string{"DEFAULT_ADMIN_PASSWORD"},
		Usage:   "default admin password",
	}
	// ForceCreateAdmin forces creating the default LDAP admin user
	ForceCreateAdmin = cli.BoolFlag{
		Name:    "force-create-admin",
		Value:   false,
		EnvVars: []string{"FORCE_CREATE_ADMIN"},
		Usage:   "force creation of the admin user even if there is a different user in the admin group",
	}
	// LdapFlags is a collection of all LDAP CLI flags
	LdapFlags = []cli.Flag{
		&GroupsOu,
		&UsersOu,
		&GroupsDn,
		&UsersDn,
		&GroupMembershipAttribute,
		&GroupMembershipUsesUID,
		&AccountAttribute,
		&GroupAttribute,
		&DefaultUserGroup,
		&DefaultAdminGroup,
		&DefaultLoginShell,
		&DefaultAdminUsername,
		&DefaultAdminPassword,
		&ForceCreateAdmin,
	}
)
