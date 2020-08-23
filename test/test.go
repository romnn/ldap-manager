package test

import (
	"context"
	"io/ioutil"
	tclog "log"
	"testing"

	ldapconfig "github.com/romnnn/ldap-manager/config"
	tc "github.com/romnnn/testcontainers"
	"github.com/romnnn/testcontainers-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	// This wil disable the native `log.Printf` calls by testcontainers-go
	tclog.SetFlags(0)
	tclog.SetOutput(ioutil.Discard)

	// This wil disable the application logger
	log.SetOutput(ioutil.Discard)

	// Note: if you want to log in tests, use `t.Log`
}

// Test ...
type Test struct {
	OpenLDAPC       testcontainers.Container
	OpenLDAPCConfig ldapconfig.OpenLDAPConfig
}

const (
	parallel = false
)

// Setup ...
func (test *Test) Setup(t *testing.T) *Test {
	var err error
	if parallel {
		t.Parallel()
	}

	containerOptions := tc.ContainerOptions{
		ContainerRequest: testcontainers.ContainerRequest{},
	}

	// Start mongodb container
	options := ContainerOptions{
		ContainerOptions: containerOptions,
		OpenLDAPConfig:   ldapconfig.OpenLDAPConfig{},
	}
	test.OpenLDAPC, test.OpenLDAPCConfig, err = StartOpenLDAPContainer(context.Background(), options)
	if err != nil {
		t.Fatalf("failed to start the OpenLDAP container: %v", err)
		return test
	}
	return test
}

// Teardown ...
func (test *Test) Teardown() {
	if test.OpenLDAPC != nil {
		_ = test.OpenLDAPC.Terminate(context.Background())
	}
}
