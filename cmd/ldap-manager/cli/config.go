package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	LdapHost = cli.StringFlag{
		Name:    "openldap-host",
		Value:   "localhost",
		EnvVars: []string{"OPENLDAP_HOST"},
		Usage:   "openldap host",
	}
	LdapPort = cli.IntFlag{
		Name:    "openldap-port",
		Value:   389,
		EnvVars: []string{"OPENLDAP_PORT"},
		Usage:   "openldap port",
	}
	LdapProtocol = cli.StringFlag{
		Name:    "openldap-protocol",
		Value:   "ldap",
		EnvVars: []string{"OPENLDAP_PROTOCOL"},
		Usage:   "openldap protocol",
	}
	LdapAdminPassword = cli.StringFlag{
		Name:    "openldap-admin-password",
		Value:   "admin",
		EnvVars: []string{"OPENLDAP_ADMIN_PASSWORD"},
		Usage:   "openldap admin password",
	}
	LdapConfigPassword = cli.StringFlag{
		Name:    "openldap-config-password",
		Value:   "config",
		EnvVars: []string{"OPENLDAP_CONFIG_PASSWORD"},
		Usage:   "openldap config password",
	}
	LdapReadOnlyUser = cli.StringFlag{
		Name:    "openldap-readonly-user",
		Value:   "", // no readonly user
		EnvVars: []string{"OPENLDAP_READONLY_USER"},
		Usage:   "openldap readonly user",
	}
	LdapReadOnlyPassword = cli.StringFlag{
		Name:    "openldap-readonly-password",
		Value:   "", // no readonly user
		EnvVars: []string{"OPENLDAP_READONLY_PASSWORD"},
		Usage:   "openldap readonly password",
	}
	LdapOrganization = cli.StringFlag{
		Name:    "openldap-organization",
		Value:   "Example Inc.",
		EnvVars: []string{"OPENLDAP_ORGANIZATION"},
		Usage:   "openldap organization",
	}
	LdapDomain = cli.StringFlag{
		Name:    "openldap-domain",
		Value:   "example.org",
		EnvVars: []string{"OPENLDAP_DOMAIN"},
		Usage:   "openldap domain",
	}
	LdapBaseDn = cli.StringFlag{
		Name:    "openldap-base-dn",
		Value:   "dc=example,dc=org",
		EnvVars: []string{"OPENLDAP_BASE_DN"},
		Usage:   "openldap base DN",
	}
	LdapTls = cli.BoolFlag{
		Name:    "openldap-tls",
		Value:   false,
		EnvVars: []string{"OPENLDAP_TLS"},
		Usage:   "openldap tls",
	}
	LdapUseRfc2307Bis = cli.BoolFlag{
		Name:    "openldap-use-rfc2307bis",
		Value:   true,
		EnvVars: []string{"OPENLDAP_USE_RFC2307BIS"},
		Usage:   "openldap use RFC2307BIS schema",
	}
	LdapConfigFlags = []cli.Flag{
		&LdapHost,
		&LdapPort,
		&LdapProtocol,
		&LdapAdminPassword,
		&LdapConfigPassword,
		&LdapReadOnlyUser,
		&LdapReadOnlyPassword,
		&LdapOrganization,
		&LdapDomain,
		&LdapBaseDn,
		&LdapTls,
		&LdapUseRfc2307Bis,
	}
)
