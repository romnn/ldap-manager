package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestDeleteGroupMissing tests deleting a missing group will fail
func TestDeleteGroupMissing(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// assert deleting a non-existent group failed
	if err := test.Manager.DeleteGroup(&pb.DeleteGroupRequest{
		Name: "group-that-does-not-exist",
	}); err == nil {
		t.Fatalf("expected error deleting group that does not exist")
	}
}

// TestDeleteDefaultGroup tests deleting default groups will fail
func TestDeleteDefaultGroup(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// make sure deleting the users group is not allowed
	for _, group := range []string{
		test.Manager.DefaultUserGroup,
		test.Manager.DefaultAdminGroup,
	} {
		if err := test.Manager.DeleteGroup(&pb.DeleteGroupRequest{
			Name: group,
		}); err == nil {
			t.Errorf("expected error deleting default group %q", group)
		}
	}
}

// TestDeleteGroup tests deleting a valid group
func TestDeleteGroup(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// add a user
	username := "some-user"
	if err := test.Manager.NewUser(&pb.NewUserRequest{
		Username:  username,
		Password:  "Hallo Welt",
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}); err != nil {
		t.Fatalf("failed to add new user: %v", err)
	}

	// now add a valid group with the user
	strict := true
	groupName := "my-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		t.Fatalf(
			"failed to add group %q with member %v: %v",
			groupName, username, err,
		)
	}

	// assert the group is found
	if _, err := test.Manager.GetGroupByName(groupName); err != nil {
		t.Fatalf(
			"failed to get group %q: %v",
			groupName, err,
		)
	}

	// now delete the group
	if err := test.Manager.DeleteGroup(&pb.DeleteGroupRequest{
		Name: groupName,
	}); err != nil {
		t.Fatalf(
			"failed to delete group %q: %v",
			groupName, err,
		)
	}

	// assert that the group is no longer found
	if _, err := test.Manager.GetGroupByName(groupName); err == nil {
		t.Fatalf(
			"unexpected success getting deleted group %q",
			groupName,
		)
	}
}
