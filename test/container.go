package test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	ldapconfig "github.com/romnnn/ldap-manager/config"
	tc "github.com/romnnn/testcontainers"
	"github.com/romnnn/testcontainers-go"
	"github.com/romnnn/testcontainers-go/wait"
)

// ContainerOptions ...
type ContainerOptions struct {
	tc.ContainerOptions
	tc.ContainerConfig
	ldapconfig.OpenLDAPConfig
}

const (
	defaultOpenLDAPPort = 389
)

// StartOpenLDAPContainer ...
func StartOpenLDAPContainer(ctx context.Context, options ContainerOptions) (openldapC testcontainers.Container, openldapCConfig ldapconfig.OpenLDAPConfig, err error) {
	openLDAPPort, _ := nat.NewPort("", strconv.Itoa(defaultOpenLDAPPort))

	defaultOptions := ContainerOptions{
		OpenLDAPConfig: ldapconfig.NewOpenLDAPConfig(),
	}

	tc.MergeOptions(&defaultOptions, options)

	var env = make(map[string]string)
	env["LDAP_ORGANISATION"] = defaultOptions.LDAPOrganization
	env["LDAP_DOMAIN"] = defaultOptions.LDAPDomain
	env["LDAP_BASE_DN"] = defaultOptions.LDAPBaseDN

	env["LDAP_ADMIN_PASSWORD"] = defaultOptions.LDAPAdminPassword
	env["LDAP_CONFIG_PASSWORD"] = defaultOptions.LDAPConfigPassword

	env["LDAP_READONLY_USER"] = strconv.FormatBool(defaultOptions.LDAPReadonlyUser)
	env["LDAP_READONLY_USER_USERNAME"] = defaultOptions.LDAPReadonlyUserUsername
	env["LDAP_READONLY_USER_PASSWORD"] = defaultOptions.LDAPReadonlyUserPassword

	env["LDAP_TLS"] = strconv.FormatBool(defaultOptions.LDAPTLS)
	env["LDAP_RFC2307BIS_SCHEMA"] = strconv.FormatBool(defaultOptions.LDAPRFC2307BISSchema)

	timeout := options.ContainerOptions.StartupTimeout
	if int64(timeout) < 1 {
		timeout = 5 * time.Minute // Default timeout
	}

	req := testcontainers.ContainerRequest{
		Env:          env,
		Image:        "osixia/openldap",
		ExposedPorts: []string{string(openLDAPPort)},
		WaitingFor:   wait.ForLog("slapd starting").WithStartupTimeout(timeout),
	}

	tc.MergeRequest(&req, &options.ContainerOptions.ContainerRequest)

	tc.ClientMux.Lock()
	openldapC, err = testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	tc.ClientMux.Unlock()
	if err != nil {
		err = fmt.Errorf("Failed to start OpenLDAP container: %v", err)
		return
	}

	host, err := openldapC.Host(ctx)
	if err != nil {
		err = fmt.Errorf("Failed to get OpenLDAP container host: %v", err)
		return
	}

	port, err := openldapC.MappedPort(ctx, openLDAPPort)
	if err != nil {
		err = fmt.Errorf("Failed to get exposed OpenLDAP container port: %v", err)
		return
	}

	openldapCConfig = defaultOptions.OpenLDAPConfig
	openldapCConfig.Host = host
	openldapCConfig.Port = port.Int()

	if options.CollectLogs {
		options.ContainerConfig.Log = new(tc.LogCollector)
		go tc.EnableLogger(openldapC, options.ContainerConfig.Log)
	}
	return
}
