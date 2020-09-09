package ldapmanager

import (
	"context"
	"io/ioutil"
	tclog "log"
	"testing"

	ldapconfig "github.com/romnnn/ldap-manager/config"
	ldaptest "github.com/romnnn/ldap-manager/test"
	tc "github.com/romnnn/testcontainers"
	"github.com/romnnn/testcontainers-go"
	log "github.com/sirupsen/logrus"
)

const (
	parallel = false

	skipAccountTests        = false
	skipChangePasswordTests = false
	skipGroupTests          = false
	skipGroupMemberTests    = false
	skipSetupTests          = false
)

// Test ...
type Test struct {
	OpenLDAPC       testcontainers.Container
	OpenLDAPCConfig ldapconfig.OpenLDAPConfig
	Manager         *LDAPManager
}

func (test *Test) setup(t *testing.T, skipSetupLDAP bool) *Test {
	var err error
	if parallel {
		t.Parallel()
	}

	// This will disable the native `log.Printf` calls by testcontainers-go
	tclog.SetFlags(0)
	tclog.SetOutput(ioutil.Discard)

	// This wil disable the application logger
	log.SetOutput(ioutil.Discard)

	// Note: if you want to log in tests, use `t.Log`

	containerOptions := tc.ContainerOptions{
		ContainerRequest: testcontainers.ContainerRequest{},
	}

	// Start mongodb container
	options := ldaptest.ContainerOptions{
		ContainerOptions: containerOptions,
		OpenLDAPConfig:   ldapconfig.OpenLDAPConfig{},
	}
	test.OpenLDAPC, test.OpenLDAPCConfig, err = ldaptest.StartOpenLDAPContainer(context.Background(), options)
	if err != nil {
		t.Fatalf("failed to start the OpenLDAP container: %v", err)
		return test
	}

	// create and setup the LDAP Manager service
	test.Manager = NewLDAPManager(test.OpenLDAPCConfig)
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
	if test.OpenLDAPC != nil {
		_ = test.OpenLDAPC.Terminate(context.Background())
	}
}
