package ldapmanager

import (
	"context"
	// "io/ioutil"
	// tclog "log"
	"testing"

	ldapconfig "github.com/romnn/ldap-manager/config"
	ldaptest "github.com/romnn/ldap-manager/testing"
	tc "github.com/romnn/testcontainers"
	// log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
)

const (
	parallel        = false
	enableDebugLogs = false

	skipAccountTests        = false
	skipChangePasswordTests = false
	skipGroupTests          = false
	skipGroupMemberTests    = false
	skipSetupTests          = false
)

// Test ...
type Test struct {
	LDAPContainer *ldaptest.Container
	Manager       *LDAPManager
}

func (test *Test) setup(t *testing.T, skipSetupLDAP bool) *Test {
	var err error
	if parallel {
		t.Parallel()
	}
	// if !enableDebugLogs {
	// 	// disable the native `log.Printf` calls by testcontainers-go
	// 	tclog.SetFlags(0)
	// 	tclog.SetOutput(ioutil.Discard)
	// 	// disable the application logger
	// 	log.SetOutput(ioutil.Discard)
	// }

	containerOptions := tc.ContainerOptions{
		ContainerRequest: testcontainers.ContainerRequest{},
	}

	// Start OpenLDAP container
	options := ldaptest.ContainerOptions{
		ContainerOptions: containerOptions,
		OpenLDAPConfig:   ldapconfig.OpenLDAPConfig{},
	}

	container, err := ldaptest.StartOpenLDAP(context.Background(), options)
	if err != nil {
		t.Fatalf("failed to start the OpenLDAP container: %v", err)
	}
	test.LDAPContainer = &container

	// create and setup the manager
	test.Manager = NewLDAPManager(container.Config)
	test.Manager.DefaultAdminUsername = "ldapadmin"
	test.Manager.DefaultAdminPassword = "123456"
	if err := test.Manager.Setup(skipSetupLDAP); err != nil {
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
	if test.LDAPContainer != nil {
		test.LDAPContainer.Terminate(context.Background())
	}
}
