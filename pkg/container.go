package pkg

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	// tc "github.com/romnn/testcontainers"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	timeout  = 5 * time.Minute
	imageTag = "latest"
)

// ContainerOptions ...
type ContainerOptions struct {
	// tc.ContainerOptions
	// tc.ContainerConfig
	ldapconfig.OpenLDAPConfig
	ImageTag string
}

// Container ...
type Container struct {
	Container testcontainers.Container
	// tc.ContainerConfig
	ldapconfig.OpenLDAPConfig
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

	// defaultOptions := ContainerOptions{
	// 	OpenLDAPConfig: ldapconfig.NewOpenLDAPConfig(),
	// }

	// tc.MergeOptions(&defaultOptions, options)

	var env = make(map[string]string)
	env["LDAP_ORGANISATION"] = options.Organization
	env["LDAP_DOMAIN"] = options.Domain
	env["LDAP_BASE_DN"] = options.BaseDN

	env["LDAP_ADMIN_PASSWORD"] = options.AdminPassword
	env["LDAP_CONFIG_PASSWORD"] = options.ConfigPassword

	env["LDAP_READONLY_USER"] = strconv.FormatBool(options.ReadonlyUser)
	env["LDAP_READONLY_USER_USERNAME"] = options.ReadonlyUserUsername
	env["LDAP_READONLY_USER_PASSWORD"] = options.ReadonlyUserPassword

	env["LDAP_TLS"] = strconv.FormatBool(options.TLS)
	env["LDAP_RFC2307BIS_SCHEMA"] = strconv.FormatBool(options.UseRFC2307BISSchema)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("osixia/openldap:%s", imageTag),
			Env:          env,
			ExposedPorts: []string{string(port)},
			// WaitingFor:   wait.ForLog("slapd starting").WithStartupTimeout(timeout),
			WaitingFor:   wait.ForListeningPort(port).WithStartupTimeout(timeout),
		},
		Started: true,
	}
	openLDAPContainer, err := testcontainers.GenericContainer(ctx, req)
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

	container.OpenLDAPConfig = options.OpenLDAPConfig
	container.OpenLDAPConfig.Host = host
	container.OpenLDAPConfig.Port = realPort.Int()

	return container, nil
}