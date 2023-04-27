package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestIsGroupMember tests checking if a user is member of a group
func TestIsGroupMember(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	strict := false
	groupName := "test-group"
	usernames := []string{"user1", "user2"}

	// users must be created first
	for _, username := range usernames {
		if err := test.Manager.NewUser(&pb.NewUserRequest{
			Username:  username,
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}); err != nil {
			t.Fatalf(
				"failed to add new user: %v",
				err,
			)
		}
	}

	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: usernames,
	}, strict); err != nil {
		t.Fatalf(
			"failed to add new group: %v",
			err,
		)
	}

	// assert every user is member of the users group by default
	for _, username := range usernames {
		userGroup := test.Manager.DefaultUserGroup
		memberStatus, _ := test.isGroupMember(t, username, userGroup, true)
		if !memberStatus.GetIsMember() {
			t.Fatalf(
				"expected user %q to be a member of group %q",
				username, userGroup,
			)
		}
	}

	// assert every user is member of the new group
	for _, username := range usernames {
		memberStatus, _ := test.isGroupMember(t, username, groupName, true)
		if !memberStatus.GetIsMember() {
			t.Fatalf(
				"expected user %q to be a member of group %q",
				username, groupName,
			)
		}
	}

	// create group with only the first user
	username := usernames[0]
	groupName2 := "another-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName2,
		Members: []string{username},
	}, strict); err != nil {
		t.Fatalf(
			"failed to add new group: %v",
			err,
		)
	}
	memberStatus, _ := test.isGroupMember(t, username, groupName2, true)
	if !memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to be a member of group %q",
			username, groupName2,
		)
	}
	username = usernames[1]
	memberStatus, _ = test.isGroupMember(t, username, groupName2, false)
	if memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to not be a member of group %q",
			username, groupName2,
		)
	}

	// assert a user is not longer member of any group after it is deleted
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf(
			"failed to get list of all groups: %v",
			err,
		)
	}
	t.Log(PrettyPrint(groups))

	username = usernames[1]
	keepGroups := false
	if err := test.Manager.DeleteUser(&pb.DeleteUserRequest{
		Username: username,
	}, keepGroups); err != nil {
		t.Fatalf(
			"failed to delete user %q: %v",
			username, err,
		)
	}
	for _, group := range groups.GetGroups() {
		memberStatus, _ = test.isGroupMember(t, username, group.GetName(), false)
		if memberStatus.GetIsMember() {
			t.Errorf(
				"expected user %q to not be a member of group %q",
				username, group.GetName(),
			)
		}
	}
}

// TestIsGroupMemberMissing tests checking group membership where either the user or group does not exist
func TestIsGroupMemberMissing(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// assert sure a non-existent user is not a member of an existent group
	username := "i-dont-exist"
	groupName := test.Manager.DefaultUserGroup
	memberStatus, _ := test.isGroupMember(t, username, groupName, false)
	if memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to not be a member of group %q",
			username, groupName,
		)
	}

	// make sure an existent user is not a member of a non-existent group
	username = test.Manager.DefaultAdminUsername
	groupName = "group-that-is-ficticious"
	memberStatus, err := test.isGroupMember(t, username, groupName, false)
	_, missingGroup := err.(*ZeroOrMultipleGroupsError)
	if err == nil || !missingGroup {
		t.Fatalf(
			"expected error due to missing group %q",
			groupName,
		)
	}
}
