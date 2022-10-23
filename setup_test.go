package ldapmanager

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestSetup ...
func TestSetup(t *testing.T) {
	if skipSetupTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// check if the default admin and user groups were created
	if _, err := test.Manager.GetGroup(&pb.GetGroupRequest{Name: test.Manager.DefaultUserGroup}); err != nil {
		t.Errorf("setup failed: failed to get default users group %q: %v", test.Manager.DefaultUserGroup, err)
	}
	if _, err := test.Manager.GetGroup(&pb.GetGroupRequest{Name: test.Manager.DefaultAdminGroup}); err != nil {
		t.Errorf("setup failed: failed to get default admin group %q: %v", test.Manager.DefaultAdminGroup, err)
	}

	// Check if the default admin user was created
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: test.Manager.DefaultAdminUsername, Password: test.Manager.DefaultAdminPassword}); err != nil {
		t.Errorf("setup failed: failed to authenticate as admin %q: %v", test.Manager.DefaultAdminGroup, err)
	}

	// check if the default admin user is in the admins group
	adminsMemberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: test.Manager.DefaultAdminUsername,
		Group:    test.Manager.DefaultAdminGroup,
	})
	if err != nil {
		t.Errorf("setup failed: failed to check if admin user %q is in group %q: %v", test.Manager.DefaultAdminUsername, test.Manager.DefaultAdminGroup, err)
	}
	if isAdmin := adminsMemberStatus.GetIsMember(); !isAdmin {
		t.Errorf("setup failed: default admin user %q is not an admin (in group %q)", test.Manager.DefaultAdminUsername, test.Manager.DefaultAdminGroup)
	}

	// check if the default admin user is in the users group as well
	usersMemberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: test.Manager.DefaultAdminUsername,
		Group:    test.Manager.DefaultUserGroup,
	})
	if err != nil {
		t.Errorf("setup failed: failed to check if admin user %q is in group %q: %v", test.Manager.DefaultAdminUsername, test.Manager.DefaultUserGroup, err)
	}
	if isUser := usersMemberStatus.GetIsMember(); !isUser {
		t.Errorf("setup failed: default admin user %q is not a user (in group %q)", test.Manager.DefaultAdminUsername, test.Manager.DefaultUserGroup)
	}
}

// TestForceSetup ...
func TestForceSetup(t *testing.T) {
	if skipSetupTests {
		t.Skip()
	}
	test := new(Test).SkipSetup(t)
	defer test.Teardown()

	differentAdminUser := &pb.NewAccountRequest{
		Account: &pb.Account{
			Username:  "differentAdmin",
			Password:  "differentAdmin",
			FirstName: "changeme",
			LastName:  "changeme",
			Email:     "changeme@changeme.com",
		},
	}

	_ = test.Manager.setupGroupsOU()
	_ = test.Manager.setupUsersOU()
	_ = test.Manager.setupLastGID()
	_ = test.Manager.setupLastUID()

	// create a different admin user
	if err := test.Manager.NewAccount(differentAdminUser, pb.HashingAlgorithm_DEFAULT); err != nil {
		t.Fatalf("failed to create different admin account: %v", err)
	}
	// create the group
	strict := false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{Name: test.Manager.DefaultAdminGroup, Members: []string{
		differentAdminUser.GetAccount().GetUsername(),
	}}, strict); err != nil {
		t.Fatalf("failed to create admins group: %v", err)
	}

	if err := test.Manager.SetupLDAP(); err != nil {
		t.Fatalf("failed to setup ldap manager service: %v", err)
	}

	// make sure we cannot authenticate with the default admin user because an admin already existed
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: test.Manager.DefaultAdminUsername, Password: test.Manager.DefaultAdminPassword}); err == nil {
		t.Errorf("expected error when authenticating as the default admin %q when another admin account already existed", test.Manager.DefaultAdminUsername)
	}

	// make sure the admin is created when forced
	test.Manager.ForceCreateAdmin = true
	if err := test.Manager.SetupLDAP(); err != nil {
		t.Fatalf("failed to setup ldap manager service: %v", err)
	}
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: test.Manager.DefaultAdminUsername, Password: test.Manager.DefaultAdminPassword}); err != nil {
		t.Errorf("failed to authenticate as the default admin %q after forced creation: %v", test.Manager.DefaultAdminUsername, err)
	}
}
