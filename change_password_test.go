package ldapmanager

import (
	"testing"

	ldaptest "github.com/romnnn/ldap-manager/test"
)

// TestChangePassword ...
func TestChangePassword(t *testing.T) {
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}

	// Add user with initial password "123"
	username := "testuser"
	initialPassword := "123"
	if err := manager.NewAccount(&NewAccountRequest{
		Username:  username,
		Password:  initialPassword,
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}); err != nil {
		t.Fatalf("failed to add user %q: %v", username, err)
	}

	// Test we can authenticate with password "123"
	if _, err := manager.AuthenticateUser(username, initialPassword); err != nil {
		t.Fatalf("failed to authenticate user %q with password %q: %v", username, initialPassword, err)
	}

	// Invalid change password request
	if err := manager.ChangePassword(&ChangePasswordRequest{
		Username: username,
		Password: "", // invalid
	}); err == nil {
		t.Fatalf("expected error changing the password for user %q to be empty", username)
	}

	// Valid change password request
	newPassword := "456"
	if err := manager.ChangePassword(&ChangePasswordRequest{
		Username: username,
		Password: newPassword, // valid
	}); err != nil {
		t.Fatalf("failed to change password of user %q to %q: %v", username, newPassword, err)
	}

	// Assert we can no longer authenticate with the old password
	if _, err := manager.AuthenticateUser(username, initialPassword); err == nil {
		t.Fatalf("expected error authenticating user %q with the initial password %q: %v", username, initialPassword, err)
	}

	// Assert we can authenticate with the new password
	if _, err := manager.AuthenticateUser(username, newPassword); err != nil {
		t.Fatalf("failed to authenticate user %q with the new password %q: %v", username, newPassword, err)
	}

	// Assert the number of users did not change in the process
	userList, _ := manager.GetUserList(&GetUserListRequest{})
	if len(userList) != 1 {
		t.Fatalf("expected exactly one user, but got %d", len(userList))
	}
}
