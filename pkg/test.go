package pkg

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
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

// Starts starts the container
func (test *Test) Start(t *testing.T) *Test {
	var err error
	t.Parallel()

	// start OpenLDAP container
	options := ContainerOptions{
		OpenLDAPConfig: ldapconfig.NewOpenLDAPConfig(),
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
	if err := test.Manager.Connect(); err != nil {
		t.Fatalf("failed to connect to OpenLDAP: %v", err)
	}
	return test
}

// Setup runs the setup of the manager
func (test *Test) Setup(t *testing.T) *Test {
	if test.Manager == nil {
		t.Fatal("must call test.Start(..) before running setup")
	}
	if err := test.Manager.Setup(); err != nil {
		t.Fatalf("failed to setup manager: %v", err)
	}
	return test
}

// Teardown stops the container
func (test *Test) Teardown() {
	if test.Container != nil {
		test.Container.Terminate(context.Background())
	}
}
