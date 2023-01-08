package pkg

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/docker/go-connections/nat"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	timeout  = 5 * time.Minute
	imageTag = "latest"
)

// ContainerOptions describes options for the container
type ContainerOptions struct {
	ldapconfig.Config
	ImageTag string
	Networks []string
}

// Container holds the LDAP container
type Container struct {
	ldapconfig.Config
	Container testcontainers.Container
}

// Terminate terminates the container
func (c *Container) Terminate(ctx context.Context) {
	if c.Container != nil {
		c.Container.Terminate(ctx)
	}
}

const (
	// OpenLDAPPort is the OpenLDAP protocol port
	OpenLDAPPort = 389
)

// StartOpenLDAP starts the OpenLDAP container
func StartOpenLDAP(ctx context.Context, options ContainerOptions) (Container, error) {
	var container Container
	port, err := nat.NewPort("", strconv.Itoa(OpenLDAPPort))
	if err != nil {
		return container, fmt.Errorf("failed to build port: %v", err)
	}

	var env = make(map[string]string)
	env["LDAP_ORGANISATION"] = options.Organization
	env["LDAP_DOMAIN"] = options.Domain
	env["LDAP_BASE_DN"] = options.BaseDN

	env["LDAP_ADMIN_PASSWORD"] = options.AdminPassword

	env["LDAP_READONLY_USER"] = "true" // strconv.FormatBool(options.ReadOnlyUser)
	env["LDAP_READONLY_USER_USERNAME"] = options.ReadOnlyUsername
	env["LDAP_READONLY_USER_PASSWORD"] = options.ReadOnlyPassword
	// env["LDAP_CONFIG_PASSWORD"] = "blabla123"

	env["CONTAINER_LOG_LEVEL"] = strconv.Itoa(16)

	// https://www.openldap.org/doc/admin24/slapdconfig.html
	env["LDAP_LOG_LEVEL"] = strconv.Itoa(-1) // 16
	env["LDAP_TLS"] = strconv.FormatBool(options.TLS)
	env["LDAP_RFC2307BIS_SCHEMA"] = strconv.FormatBool(options.UseRFC2307BISSchema)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("osixia/openldap:%s", imageTag),
			Networks:     options.Networks,
			Env:          env,
			ExposedPorts: []string{string(port)},
			Cmd: []string{
				"--loglevel",
				"debug",
			},
			WaitingFor: wait.ForListeningPort(port).WithStartupTimeout(timeout),
		},
		Started: true,
	}
	openLDAPContainer, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, fmt.Errorf(
			"failed to start container: %v",
			err,
		)
	}
	container.Container = openLDAPContainer

	host, err := openLDAPContainer.Host(ctx)
	if err != nil {
		return container, fmt.Errorf(
			"failed to get container host: %v",
			err,
		)
	}

	realPort, err := openLDAPContainer.MappedPort(ctx, port)
	if err != nil {
		return container, fmt.Errorf(
			"failed to get exposed container port: %v",
			err,
		)
	}

	container.Config = options.Config
	container.Config.Host = host
	container.Config.Port = realPort.Int()

	return container, nil
}
