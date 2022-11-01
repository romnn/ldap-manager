package pkg

import (
	"context"
	// "io/ioutil"
	// tclog "log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	// ldapmanager "github.com/romnn/ldap-manager/pkg"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	// tc "github.com/romnn/testcontainers"
	// log "github.com/sirupsen/logrus"
	// "github.com/testcontainers/testcontainers-go"
)

// EqualProto checks if two proto messages are equal
// if they are not equal, a diff describes how they differ
func EqualProto(a proto.Message, b proto.Message) (bool, string) {
	opts := cmp.Options{
		protocmp.Transform(),
	}
	equal := cmp.Equal(a, b, opts)
	diff := cmp.Diff(a, b, opts)
	return equal, diff
}

// Test wraps a pre-configured OpenLDAP container and Manager instance
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
