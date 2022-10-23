package main

import (
	"context"
	"io/ioutil"
	tclog "log"
	"net"
	"sync"
	"testing"

	// gogrpcservice "github.com/romnn/go-grpc-service"
	// "github.com/romnn/go-grpc-service/auth"
	ldapmanager "github.com/romnn/ldap-manager"
	ldapbase "github.com/romnn/ldap-manager/cmd/ldap-manager/base"
	ldapgrpc "github.com/romnn/ldap-manager/cmd/ldap-manager/grpc"
	ldaphttp "github.com/romnn/ldap-manager/cmd/ldap-manager/http"
	"github.com/romnn/ldap-manager/pkg/config"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	ldaptest "github.com/romnn/ldap-manager/test"
	tc "github.com/romnn/testcontainers"
	log "github.com/sirupsen/logrus"
	"github.com/testcontainers/testcontainers-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
)

const (
	parallel        = false
	enableDebugLogs = false

	bufSize = 1024 * 1024
)

// Test ...
type Test struct {
	OpenLDAPC         testcontainers.Container
	OpenLDAPCConfig   ldap.OpenLDAPConfig
	ManagerEndpoint   *grpc.ClientConn
	ManagerClient     pb.LDAPManagerClient
	ManagerServer     *ldapbase.LDAPManagerServer
	ManagerGRPCServer *ldapgrpc.LDAPManagerServer
	ManagerHTTPServer *ldaphttp.LDAPManagerServer
}

func (test *Test) setup(t *testing.T, skipSetupLDAP bool) *Test {
	var err error
	if parallel {
		t.Parallel()
	}
	if !enableDebugLogs {
		// disable the native `log.Printf` calls by testcontainers-go
		tclog.SetFlags(0)
		tclog.SetOutput(ioutil.Discard)
		// disable the application logger
		log.SetOutput(ioutil.Discard)
	}

	containerOptions := tc.ContainerOptions{
		ContainerRequest: testcontainers.ContainerRequest{},
	}

	// Start OpenLDAP container
	options := ldaptest.ContainerOptions{
		ContainerOptions: containerOptions,
		OpenLDAPConfig:   config.OpenLDAPConfig{},
	}
	container, err := ldaptest.StartOpenLDAP(context.Background(), options)
	if err != nil {
		t.Fatalf("failed to start the OpenLDAP container: %v", err)
		return test
	}

	authenticator := &auth.Authenticator{
		ExpireSeconds: 60 * 60 * 10,
		Issuer:        "issuer@example.com",
		Audience:      "example.com",
	}
	if err := authenticator.SetupKeys(&auth.AuthenticatorKeyConfig{
		Generate: true,
	}); err != nil {
		t.Fatalf("failed to setup keys: %v", err)
	}

	// create and setup the LDAP Manager service
	manager := ldapmanager.NewLDAPManager(test.OpenLDAPCConfig)
	manager.DefaultAdminUsername = "ldapadmin"
	manager.DefaultAdminPassword = "123456"

	if err := manager.Setup(skipSetupLDAP); err != nil {
		t.Fatal(err)
	}

	test.ManagerServer = &ldapbase.LDAPManagerServer{
		Service: gogrpcservice.Service{
			Name:               "ldap manager service",
			HTTPHealthCheckURL: "/healthz",
		},
		Authenticator: authenticator,
		AuthKeyConfig: &auth.AuthenticatorKeyConfig{
			Generate: true,
		},
		Manager: manager,
		Static:  false,
	}

	// create listeners
	grpcListener := bufconn.Listen(bufSize)
	httpListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	test.ManagerGRPCServer = ldapgrpc.NewGRPCLDAPManagerServer(test.ManagerServer, grpcListener)
	ctx := context.Background()
	var wg sync.WaitGroup
	wg.Add(2)

	go test.ManagerGRPCServer.Serve(ctx, &wg)

	test.ManagerEndpoint, err = grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return grpcListener.Dial()
	}), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		test.Teardown()
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	test.ManagerClient = pb.NewLDAPManagerClient(test.ManagerEndpoint)
	test.ManagerHTTPServer = ldaphttp.NewHTTPLDAPManagerServer(test.ManagerServer, httpListener, test.ManagerEndpoint)
	go test.ManagerHTTPServer.Serve(ctx, &wg)

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
	if test.ManagerHTTPServer != nil {
		test.ManagerHTTPServer.Shutdown()
	}
	if test.ManagerGRPCServer != nil {
		test.ManagerGRPCServer.Shutdown()
	}
	if test.ManagerEndpoint != nil {
		test.ManagerEndpoint.Close()
	}
	if test.OpenLDAPC != nil {
		_ = test.OpenLDAPC.Terminate(context.Background())
	}
}

// TestSetup ...
func TestSetup(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	// check if the default admin and user groups were created
	if _, err := test.ManagerServer.Manager.GetGroup(&pb.GetGroupRequest{Name: test.ManagerServer.Manager.DefaultUserGroup}); err != nil {
		t.Errorf("setup failed: failed to get default users group %q: %v", test.ManagerServer.Manager.DefaultUserGroup, err)
	}
	if _, err := test.ManagerServer.Manager.GetGroup(&pb.GetGroupRequest{Name: test.ManagerServer.Manager.DefaultAdminGroup}); err != nil {
		t.Errorf("setup failed: failed to get default admin group %q: %v", test.ManagerServer.Manager.DefaultAdminGroup, err)
	}

	// Check if the default admin user was created
	if _, err := test.ManagerServer.Manager.AuthenticateUser(&pb.LoginRequest{Username: test.ManagerServer.Manager.DefaultAdminUsername, Password: test.ManagerServer.Manager.DefaultAdminPassword}); err != nil {
		t.Errorf("setup failed: failed to authenticate as admin %q: %v", test.ManagerServer.Manager.DefaultAdminGroup, err)
	}

	// check if the default admin user is in the admins group
	adminsMemberStatus, err := test.ManagerServer.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: test.ManagerServer.Manager.DefaultAdminUsername,
		Group:    test.ManagerServer.Manager.DefaultAdminGroup,
	})
	if err != nil {
		t.Errorf("setup failed: failed to check if admin user %q is in group %q: %v", test.ManagerServer.Manager.DefaultAdminUsername, test.ManagerServer.Manager.DefaultAdminGroup, err)
	}
	if isAdmin := adminsMemberStatus.GetIsMember(); !isAdmin {
		t.Errorf("setup failed: default admin user %q is not an admin (in group %q)", test.ManagerServer.Manager.DefaultAdminUsername, test.ManagerServer.Manager.DefaultAdminGroup)
	}

	// check if the default admin user is in the users group as well
	usersMemberStatus, err := test.ManagerServer.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: test.ManagerServer.Manager.DefaultAdminUsername,
		Group:    test.ManagerServer.Manager.DefaultUserGroup,
	})
	if err != nil {
		t.Errorf("setup failed: failed to check if admin user %q is in group %q: %v", test.ManagerServer.Manager.DefaultAdminUsername, test.ManagerServer.Manager.DefaultUserGroup, err)
	}
	if isUser := usersMemberStatus.GetIsMember(); !isUser {
		t.Errorf("setup failed: default admin user %q is not a user (in group %q)", test.ManagerServer.Manager.DefaultAdminUsername, test.ManagerServer.Manager.DefaultUserGroup)
	}

	cases := []struct {
		username, password string
		success, admin     bool
	}{
		{
			username: test.ManagerServer.Manager.DefaultAdminUsername,
			password: test.ManagerServer.Manager.DefaultAdminPassword,
			success:  true,
			admin:    true,
		},
		{
			username: test.ManagerServer.Manager.DefaultAdminUsername,
			password: "wrong-password",
			success:  false,
		},
	}

	for _, c := range cases {
		response, err := test.ManagerClient.Login(context.Background(), &pb.LoginRequest{
			Username: c.username,
			Password: c.password,
		})
		if c.success && err != nil {
			t.Errorf("login failed unexpectedly for %s:%s: %v", c.username, c.password, err)
			continue
		}
		if !c.success && err == nil {
			t.Errorf("expected login error for %s:%s", c.username, c.password)
			continue
		}
		if !c.success && err != nil {
			continue
		}

		if response.GetIsAdmin() != c.admin {
			t.Errorf("expected admin=%t but got %t for %s:%s", c.admin, response.GetIsAdmin(), c.username, c.password)
			continue
		}

		// check that the admin user is can not access routes without passing the token
		if _, err := test.ManagerClient.GetUserList(context.Background(), &pb.GetUserListRequest{}); err == nil {
			t.Error("expected error when calling GetUserList without auth metadata")
		}

		// check that the admin user can access protected routes
		ctx := metadata.NewOutgoingContext(context.Background(), metadata.New(map[string]string{
			"x-user-token": response.GetToken(),
		}))
		if _, err := test.ManagerClient.GetUserList(ctx, &pb.GetUserListRequest{}); err != nil {
			t.Errorf("unexpected error when calling GetUserList: %v", err)
		}

		// TODO: also check the gatway via REST
	}
}
