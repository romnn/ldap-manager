package pkg

import (
	"context"
	// "io/ioutil"
	// tclog "log"
	"testing"

	// ldapmanager "github.com/romnn/ldap-manager/pkg"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	// tc "github.com/romnn/testcontainers"
	// log "github.com/sirupsen/logrus"
	// "github.com/testcontainers/testcontainers-go"
)

// Test ...
type Test struct {
	Container *Container
	Manager   *LDAPManager
}

func (test *Test) setup(t *testing.T, skipSetupLDAP bool) *Test {
	var err error
	t.Parallel()

	// containerOptions := tc.ContainerOptions{
	// 	ContainerRequest: testcontainers.ContainerRequest{},
	// }

	// start OpenLDAP container
	options := ContainerOptions{
		// ContainerOptions: containerOptions,
		OpenLDAPConfig: ldapconfig.NewOpenLDAPConfig(),
		// OpenLDAPConfig:   ldapconfig.OpenLDAPConfig{},
	}
	container, err := StartOpenLDAP(context.Background(), options)
	if err != nil {
		t.Fatalf("failed to start OpenLDAP container: %v", err)
	}
	test.Container = &container

	// create and setup the LDAP Manager service
	test.Manager = NewLDAPManager(test.Container.OpenLDAPConfig)
	test.Manager.DefaultAdminUsername = "ldapadmin"
	test.Manager.DefaultAdminPassword = "123456"
	// if err := test.Manager.Setup(skipSetupLDAP); err != nil {
	if err := test.Manager.Setup(); err != nil {
		t.Fatal(err)
	}
	return test
}

// Setup ...
func (test *Test) Setup(t *testing.T) *Test {
	return test.setup(t, false)
}

// SkipSetup ...
func (test *Test) SkipSetup(t *testing.T) *Test {
	return test.setup(t, true)
}

// Teardown ...
func (test *Test) Teardown() {
	if test.Container != nil {
		test.Container.Terminate(context.Background())
	}
}
