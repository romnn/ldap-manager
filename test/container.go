package test

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/romnn/ldap-manager/pkg/config"
	tc "github.com/romnn/testcontainers"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ContainerOptions ...
type ContainerOptions struct {
	tc.ContainerOptions
	tc.ContainerConfig
	ldapconfig.OpenLDAPConfig
	ImageTag string
}

// Container ...
type Container struct {
	Container testcontainers.Container
	tc.ContainerConfig
	ldapconfig.OpenLDAPConfig
	// Host     string
	// Port     int64
	// Password string
}

// Terminate ...
func (c *Container) Terminate(ctx context.Context) {
	if c.Container != nil {
		c.Container.Terminate(ctx)
	}
}

// StartOpenLDAP ...
func StartOpenLDAP(ctx context.Context, options ContainerOptions) (Container, error) {
	var container Container
	port, err := nat.NewPort("", "389")
	if err != nil {
		return container, fmt.Errorf("failed to build port: %v", err)
	}

	defaultOptions := ContainerOptions{
		OpenLDAPConfig: ldapconfig.NewOpenLDAPConfig(),
	}

	tc.MergeOptions(&defaultOptions, options)

	var env = make(map[string]string)
	env["LDAP_ORGANISATION"] = defaultOptions.Organization
	env["LDAP_DOMAIN"] = defaultOptions.Domain
	env["LDAP_BASE_DN"] = defaultOptions.BaseDN

	env["LDAP_ADMIN_PASSWORD"] = defaultOptions.AdminPassword
	env["LDAP_CONFIG_PASSWORD"] = defaultOptions.ConfigPassword

	env["LDAP_READONLY_USER"] = strconv.FormatBool(defaultOptions.ReadonlyUser)
	env["LDAP_READONLY_USER_USERNAME"] = defaultOptions.ReadonlyUserUsername
	env["LDAP_READONLY_USER_PASSWORD"] = defaultOptions.ReadonlyUserPassword

	env["LDAP_TLS"] = strconv.FormatBool(defaultOptions.TLS)
	env["LDAP_RFC2307BIS_SCHEMA"] = strconv.FormatBool(defaultOptions.UseRFC2307BISSchema)

	timeout := options.ContainerOptions.StartupTimeout
	if int64(timeout) < 1 {
		timeout = 5 * time.Minute // Default timeout
	}

	tag := "latest"
	if options.ImageTag != "" {
		tag = options.ImageTag
	}

	req := testcontainers.ContainerRequest{
		Image:        fmt.Sprintf("osixia/openldap:%s", tag),
		Env:          env,
		ExposedPorts: []string{string(port)},
		WaitingFor:   wait.ForLog("slapd starting").WithStartupTimeout(timeout),
	}

	tc.MergeRequest(&req, &options.ContainerOptions.ContainerRequest)

	openLDAPContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return container, fmt.Errorf("failed to start container: %v", err)
	}
	container.Container = openLDAPContainer

	host, err := openLDAPContainer.Host(ctx)
	if err != nil {
		return container, fmt.Errorf("failed to get container host: %v", err)
	}

	realPort, err := openLDAPContainer.MappedPort(ctx, port)
	if err != nil {
		return container, fmt.Errorf("failed to get exposed container port: %v", err)
	}

	container.OpenLDAPConfig = defaultOptions.OpenLDAPConfig
	// openldapCConfig = defaultOptions.OpenLDAPConfig
	container.OpenLDAPConfig.Host = host
	container.OpenLDAPConfig.Port = realPort.Int()

	return container, nil
}
