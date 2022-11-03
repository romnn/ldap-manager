package pkg

import (
	"testing"

	// "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestUpdateGroup tests that adding an empty group fails
func TestUpdateGroup(t *testing.T) {
	test := new(Test).Setup(t)
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
		t.Fatalf("failed to add group %q with member %v: %v", groupName, username, err)
	}

	before, err := test.Manager.GetGroupByName(groupName)
	if err != nil {
		t.Fatalf("failed to get group %q before rename: %v", groupName, err)
	}
	t.Log(PrettyPrint(before))

	// update the group (including the name)
	newGroupName := "my-renamed-group"
	newGID := int64(2055)
	if err := test.Manager.UpdateGroup(&pb.UpdateGroupRequest{
		Name:    groupName,
		NewName: newGroupName,
		GID:     newGID,
	}); err != nil {
		t.Fatalf("failed to update and rename group from %q to %q: %v", groupName, newGroupName, err)
	}

	// assert members are left untouched
	after, err := test.Manager.GetGroupByName(newGroupName)
	if err != nil {
		t.Fatalf("failed to get the renamed group %q: %v", newGroupName, err)
	}
	t.Log(PrettyPrint(after))

	expected := &pb.Group{
		Name:    newGroupName,
		Members: before.GetMembers(),
		GID:     newGID,
	}
	if equal, diff := EqualProto(expected, after); !equal {
		t.Fatalf("unexpected group: \n%s", diff)
	}

	// assert old name is really gone
	if _, err := test.Manager.GetGroupByName(groupName); err == nil {
		t.Fatalf("expected error getting group by its old name %q", groupName)
	}
}
