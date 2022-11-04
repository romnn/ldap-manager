package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	// LdapHost configures the LDAP server host
	LdapHost = cli.StringFlag{
		Name:    "ldap-host",
		Value:   "localhost",
		EnvVars: []string{"LDAP_HOST"},
		Usage:   "LDAP host",
	}
	// LdapPort configures the LDAP server port
	LdapPort = cli.IntFlag{
		Name:    "ldap-port",
		Value:   389,
		EnvVars: []string{"LDAP_PORT"},
		Usage:   "LDAP port",
	}
	// LdapProtocol configures the LDAP server protocol
	LdapProtocol = cli.StringFlag{
		Name:    "ldap-protocol",
		Value:   "ldap",
		EnvVars: []string{"LDAP_PROTOCOL"},
		Usage:   "LDAP protocol",
	}
	// LdapAdminUsername configures the LDAP admin username
	LdapAdminUsername = cli.StringFlag{
		Name:    "ldap-admin-username",
		Value:   "admin",
		EnvVars: []string{"LDAP_ADMIN_USERNAME"},
		Usage:   "LDAP admin username",
	}
	// LdapAdminPassword configures the LDAP admin password
	LdapAdminPassword = cli.StringFlag{
		Name:    "ldap-admin-password",
		Value:   "admin",
		EnvVars: []string{"LDAP_ADMIN_PASSWORD"},
		Usage:   "LDAP admin password",
	}
	// LdapReadOnlyUser configures the LDAP read-only user
	LdapReadOnlyUser = cli.StringFlag{
		Name:    "ldap-readonly-user",
		Value:   "", // no read-only user
		EnvVars: []string{"LDAP_READONLY_USER"},
		Usage:   "LDAP read-only user",
	}
	// LdapReadOnlyPassword configures the LDAP read-only user
	LdapReadOnlyPassword = cli.StringFlag{
		Name:    "openldap-readonly-password",
		Value:   "", // no read-only user
		EnvVars: []string{"OPENLDAP_READONLY_PASSWORD"},
		Usage:   "LDAP read-only password",
	}
	// LdapOrganization configures the LDAP organization
	LdapOrganization = cli.StringFlag{
		Name:    "ldap-organization",
		Value:   "Example Inc.",
		EnvVars: []string{"LDAP_ORGANIZATION"},
		Usage:   "LDAP organization",
	}
	// LdapDomain configures the LDAP domain
	LdapDomain = cli.StringFlag{
		Name:    "ldap-domain",
		Value:   "example.org",
		EnvVars: []string{"LDAP_DOMAIN"},
		Usage:   "LDAP domain",
	}
	// LdapBaseDn configures the LDAP base DN
	LdapBaseDn = cli.StringFlag{
		Name:    "ldap-base-dn",
		Value:   "dc=example,dc=org",
		EnvVars: []string{"LDAP_BASE_DN"},
		Usage:   "LDAP base DN",
	}
	// LdapTLS configures if TLS shoudld be used for LDAP
	LdapTLS = cli.BoolFlag{
		Name:    "ldap-tls",
		Value:   false,
		EnvVars: []string{"LDAP_TLS"},
		Usage:   "LDAP use TLS",
	}
	// LdapUseRfc2307Bis configures if the LDAP server uses the RFC2307BIS schema
	LdapUseRfc2307Bis = cli.BoolFlag{
		Name:    "ldap-use-rfc2307bis",
		Value:   true,
		EnvVars: []string{"LDAP_USE_RFC2307BIS"},
		Usage:   "LDAP use RFC2307BIS schema",
	}
	// LdapConfigFlags is a set of all LDAP CLI flags
	LdapConfigFlags = []cli.Flag{
		&LdapHost,
		&LdapPort,
		&LdapProtocol,
		&LdapAdminUsername,
		&LdapAdminPassword,
		&LdapReadOnlyUser,
		&LdapReadOnlyPassword,
		&LdapOrganization,
		&LdapDomain,
		&LdapBaseDn,
		&LdapTLS,
		&LdapUseRfc2307Bis,
	}
)
