package pkg

import (
	"testing"

	// "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestAddGroupMember tests adding a group member
func TestAddGroupMember(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	username1 := "test-user-1"
	username2 := "test-user-2"
	for _, username := range []string{username1, username2} {
		req := pb.NewUserRequest{
			Username:  username,
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}

		if err := test.Manager.NewUser(&req); err != nil {
			t.Fatalf("failed to add new user: %v", err)
		}
	}

	// create a new group
	strict := true
	groupName := "test-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username1},
	}, strict); err != nil {
		t.Fatalf("failed to add new group: %v", err)
	}

	// get group members
	group, err := test.Manager.GetGroupByName(groupName)
	if err != nil {
		t.Fatalf("failed to get group %q: %v", groupName, err)
	}

	expected := &pb.Group{
		Name: groupName,
		Members: []string{
			test.Manager.UserNamed(username1),
		},
		GID: 2002,
	}

	t.Log(PrettyPrint(group))
	t.Log(PrettyPrint(expected))

	if equal, diff := EqualProto(expected, group); !equal {
		t.Fatalf("unexpected group: \n%s", diff)
	}

	memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: username1,
		Group:    groupName,
	})
	if err != nil {
		t.Fatalf("failed to get membership status of user %q for group %q: %v", username1, groupName, err)
	}
	if !memberStatus.GetIsMember() {
		t.Fatalf("user %q should be a member of group %q", username1, groupName)
	}

	memberStatus, err = test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: username2,
		Group:    groupName,
	})
	if err != nil {
		t.Fatalf("failed to get membership status of user %q for group %q: %v", username2, groupName, err)
	}
	if memberStatus.GetIsMember() {
		t.Fatalf("user %q should not yet be a member of group %q", username2, groupName)
	}

	// add username2 as group member
	allowNonExistent := false
	if err := test.Manager.AddGroupMember(&pb.GroupMember{
		Username: username2,
		Group:    groupName,
	}, allowNonExistent); err != nil {
		t.Fatalf("failed to add user %q to group %q: %v", username2, groupName, err)
	}

	memberStatus, err = test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: username2,
		Group:    groupName,
	})
	if err != nil {
		t.Fatalf("failed to get membership status of user %q for group %q: %v", username2, groupName, err)
	}
	if !memberStatus.GetIsMember() {
		t.Fatalf("user %q should be a member of group %q", username2, groupName)
	}
}

// TestAddGroupMemberMissing tests adding a group member
// when either the user or the group does not exist.
func TestAddGroupMemberMissing(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	strict := false
	groupName := "test-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{"temp-user"},
	}, strict); err != nil {
		t.Fatalf("failed to add new group: %v", err)
	}

	// add a non-existent member to an existing group
	// this will fail because the user is not present in the users group
	// and strict checking is the default
	// however, if we were adding to the users group this would succeed
	username := "i-am-not-there"
	allowNonExistent := false
	if err := test.Manager.AddGroupMember(&pb.GroupMember{
		Group:    groupName,
		Username: username,
	}, allowNonExistent); err == nil {
		t.Errorf("expected error adding user %q to group %q", username, groupName)
	}

	// add an existing user to an non-existing group
	username = "valid-user"
	if err := test.Manager.NewUser(&pb.NewUserRequest{
		Username:  username,
		Password:  "Hallo Welt",
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}); err != nil {
		t.Fatalf("failed to add new user: %v", err)
	}

	groupName = "group-that-is-ficticious"
	if err := test.Manager.AddGroupMember(&pb.GroupMember{
		Group:    groupName,
		Username: username,
	}, allowNonExistent); err == nil {
		t.Fatalf("expected error adding user %q to a group %q", username, groupName)
	}
}
