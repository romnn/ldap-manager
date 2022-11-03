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

// StartOpenLDAP starts the OpenLDAP container
func StartOpenLDAP(ctx context.Context, options ContainerOptions) (Container, error) {
	var container Container
	port, err := nat.NewPort("", "389")
	if err != nil {
		return container, fmt.Errorf("failed to build port: %v", err)
	}

	var env = make(map[string]string)
	env["LDAP_ORGANISATION"] = options.Organization
	env["LDAP_DOMAIN"] = options.Domain
	env["LDAP_BASE_DN"] = options.BaseDN

	env["LDAP_ADMIN_PASSWORD"] = options.AdminPassword

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
			WaitingFor:   wait.ForListeningPort(port).WithStartupTimeout(timeout),
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
