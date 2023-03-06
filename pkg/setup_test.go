package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestSetup tests the default LDAP setup
func TestSetup(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	userGroup := test.Manager.DefaultUserGroup
	adminGroup := test.Manager.DefaultAdminGroup

	// check if the default admin and user groups were created
	if _, err := test.Manager.GetGroupByName(userGroup); err != nil {
		t.Fatalf(
			"failed to get default user group %q: %v",
			userGroup, err,
		)
	}
	if _, err := test.Manager.GetGroupByName(adminGroup); err != nil {
		t.Fatalf(
			"failed to get default admin group %q: %v",
			adminGroup, err,
		)
	}

	// assert the default admin user was created
	adminUsername := test.Manager.DefaultAdminUsername
	adminPassword := test.Manager.DefaultAdminPassword
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: adminUsername,
		Password: adminPassword,
	}); err != nil {
		t.Errorf(
			"failed to authenticate admin user %q with password %q: %v",
			adminUsername, adminPassword, err,
		)
	}

	// assert the default admin user is in the admins group
	memberStatus, _ := test.isGroupMember(t, adminUsername, adminGroup, true)
	if !memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to be a member of group %q",
			adminUsername, adminGroup,
		)
	}

	// assert the default admin user is in the users group as well
	memberStatus, _ = test.isGroupMember(t, adminUsername, userGroup, true)
	if !memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to be a member of group %q",
			adminUsername, userGroup,
		)
	}
}

// TestForceSetup tests the default LDAP setup
func TestForceSetup(t *testing.T) {
	test := new(Test).Start(t)
	defer test.Teardown()

	defaultAdminGroup := test.Manager.DefaultAdminGroup
	defaultAdminUsername := test.Manager.DefaultAdminUsername
	defaultAdminPassword := test.Manager.DefaultAdminPassword

	differentAdminUsername := "differentAdmin"
	_ = test.Manager.setupGroupOU()
	_ = test.Manager.setupUserOU()
	_ = test.Manager.setupLastGID()
	_ = test.Manager.setupLastUID()

	// create a different admin user
	if err := test.Manager.NewUser(&pb.NewUserRequest{
		Username:  differentAdminUsername,
		Password:  "differentAdmin",
		FirstName: "changeme",
		LastName:  "changeme",
		Email:     "changeme@changeme.com",
	}); err != nil {
		t.Fatalf("failed to create different admin user: %v", err)
	}

	// create the admin group manually
	strict := false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    defaultAdminGroup,
		Members: []string{differentAdminUsername},
	}, strict); err != nil {
		_, exists := err.(*GroupAlreadyExistsError)
		if !exists {
			t.Fatalf("failed to create admin group: %v", err)
		}
	}

	if err := test.Manager.SetupLDAP(); err != nil {
		t.Fatalf("failed to setup LDAP service: %v", err)
	}

	// assert we cannot authenticate with the default admin user,
	// because an admin already existed
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: defaultAdminUsername,
		Password: defaultAdminPassword,
	}); err == nil {
		t.Errorf(
			"expected error authenticating as the default admin %q, when another admin account already existed",
			defaultAdminUsername,
		)
	}

	// assert the default admin is created when forced
	test.Manager.ForceCreateAdmin = true
	if err := test.Manager.SetupLDAP(); err != nil {
		t.Fatalf("failed to setup LDAP service: %v", err)
	}
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: defaultAdminUsername,
		Password: defaultAdminPassword,
	}); err != nil {
		t.Errorf(`failed to authenticate as the default admin %q, 
after forced creation: %v`,
			defaultAdminUsername, err,
		)
	}
}
