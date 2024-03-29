package pkg

import (
	"testing"

	recursivesort "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestChangeMissingUserPassword tests changing the password of a missing user
func TestChangeMissingUserPassword(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	err := test.Manager.ChangePassword(&pb.ChangePasswordRequest{
		Username: "non existing user",
		Password: "new password",
	})
	t.Log(err)
	if err == nil {
		t.Fatal("changing password for missing user succeeded unexpectedly")
	}
}

// TestChangeExistingUserPassword tests changing the password of an existing user
func TestChangeExistingUserPassword(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	username := "romnn"
	oldPassword := "Hallo Welt"
	req := pb.NewUserRequest{
		Username:  username,
		Password:  oldPassword,
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}

	// add the user
	if err := test.Manager.NewUser(&req); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	before, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of users: %v", err)
	}

	// check if we can authenticate the user using password
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: oldPassword,
	}); err != nil {
		t.Fatalf(
			"cannot authenticate user %q with password %q: %v",
			username, oldPassword, err,
		)
	}

	// change the password
	newPassword := "new password"
	if err := test.Manager.ChangePassword(&pb.ChangePasswordRequest{
		Username: username,
		Password: newPassword,
	}); err != nil {
		t.Fatalf(
			"failed to change password for user %q to %q: %v",
			username, newPassword, err,
		)
	}

	// check if we can authenticate the user using the new password
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: newPassword,
	}); err != nil {
		t.Fatalf(
			"cannot authenticate user %q with password %q: %v",
			username, newPassword, err,
		)
	}

	// check if authenticating the user with the old password fails
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: oldPassword,
	}); err == nil {
		t.Fatalf(
			"expected authenticating user %q with password %q to fail",
			username, oldPassword,
		)
	}

	// assert users did not change during the process
	after, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of users: %v", err)
	}

	sort := recursivesort.RecursiveSort{StructSortField: "GID"}
	sort.Sort(&before)
	sort.Sort(&after)

	if equal, diff := EqualProto(before, after); !equal {
		t.Fatalf("unexpected change in users: \n%s", diff)
	}
}
