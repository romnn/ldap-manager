package ldapmanager

import (
	"testing"

	pb "github.com/romnn/ldap-manager/grpc/ldap-manager"
)

// TestChangePassword ...
func TestChangePassword(t *testing.T) {
	if skipChangePasswordTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// Add user with initial password "123"
	username := "testuser"
	initialPassword := "123"
	if err := test.Manager.NewAccount(&pb.NewAccountRequest{
		Account: &pb.Account{
			Username:  username,
			Password:  initialPassword,
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		},
	}, pb.HashingAlgorithm_DEFAULT); err != nil {
		t.Fatalf("failed to add user %q: %v", username, err)
	}

	// Test we can authenticate with password "123"
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: username, Password: initialPassword}); err != nil {
		t.Fatalf("failed to authenticate user %q with password %q: %v", username, initialPassword, err)
	}

	userListBefore, _ := test.Manager.GetUserList(&pb.GetUserListRequest{})

	// Invalid change password request
	if err := test.Manager.ChangePassword(&pb.ChangePasswordRequest{
		Username: username,
		Password: "", // invalid
	}); err == nil {
		t.Fatalf("expected error changing the password for user %q to be empty", username)
	}

	// Valid change password request
	newPassword := "456"
	if err := test.Manager.ChangePassword(&pb.ChangePasswordRequest{
		Username: username,
		Password: newPassword, // valid
	}); err != nil {
		t.Fatalf("failed to change password of user %q to %q: %v", username, newPassword, err)
	}

	// Assert we can no longer authenticate with the old password
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: username, Password: initialPassword}); err == nil {
		t.Fatalf("expected error authenticating user %q with the initial password %q: %v", username, initialPassword, err)
	}

	// Assert we can authenticate with the new password
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: username, Password: newPassword}); err != nil {
		t.Fatalf("failed to authenticate user %q with the new password %q: %v", username, newPassword, err)
	}

	// Assert the number of users did not change in the process
	userListAfter, _ := test.Manager.GetUserList(&pb.GetUserListRequest{})
	before := len(userListBefore.GetUsers())
	after := len(userListAfter.GetUsers())
	if before != after {
		t.Fatalf("expected the number of users to not change, but got %d after and %d before", after, before)
	}
}
