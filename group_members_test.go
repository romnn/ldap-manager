package ldapmanager

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// TestIsGroupMember ...
func TestIsGroupMember(t *testing.T) {
	if skipGroupMemberTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// add sample users
	users, err := addSampleUsers(test.Manager, 3)
	if err != nil {
		t.Fatalf("failed to add sample users: %v", err)
	}

	// make sure every user is member of the users group by default
	for _, user := range users {
		if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: user, Group: test.Manager.DefaultUserGroup}); err != nil || !memberStatus.GetIsMember() {
			t.Errorf("expected user %q to be member of group %q: %v", user, test.Manager.DefaultUserGroup, err)
		}
	}

	// create a new group with the first user as the initial member
	_, groupName, err := addSampleGroup(test.Manager, "my-group", []string{users[0]}, 0)
	if err != nil {
		t.Fatalf("failed to add sample group: %v", err)
	}
	assertHasGroups(t, test.Manager, []string{groupName})

	// add the second user as well
	allowNonExistent := false
	if err := test.Manager.AddGroupMember(&pb.GroupMember{Group: groupName, Username: users[1]}, allowNonExistent); err != nil {
		t.Fatalf("failed to add user %q to group %q: %v", users[1], groupName, err)
	}

	// check that user 0 and 1 are members
	for _, user := range []string{users[0], users[1]} {
		if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: user, Group: groupName}); err != nil || !memberStatus.GetIsMember() {
			t.Errorf("expected user %q to be member of group %q: %v", user, groupName, err)
		}
	}

	// make sure a non-existent user is not a member of an existent group
	nonExistantUser := "i-dont-exist"
	if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: nonExistantUser, Group: groupName}); memberStatus.GetIsMember() {
		t.Errorf("expected non-existent user %q to be no member of group %q: %v", nonExistantUser, groupName, err)
	}

	// make sure an existent user is not a member of a non-existent group
	nonExistantGroup := "group-that-is-ficticious"
	if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: users[1], Group: nonExistantGroup}); memberStatus.GetIsMember() {
		t.Errorf("expected user %q to be no member of non-existent group %q: %v", users[1], nonExistantGroup, err)
	}

	// make sure after a user is deleted it is not longer a member of any group
	allGroups, _ := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	leaveGroups := false
	if err := test.Manager.DeleteAccount(&pb.DeleteAccountRequest{Username: users[0]}, leaveGroups); err != nil {
		t.Fatalf("failed to delete user %q: %v", users[0], err)
	}
	for _, group := range allGroups.GetGroups() {
		if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: users[0], Group: group}); memberStatus.GetIsMember() {
			t.Errorf("expected deleted user %q to be no member of group %q: %v", users[0], group, err)
		}
	}
}

// TestGetGroup ...
func TestGetGroup(t *testing.T) {
	if skipGroupMemberTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	users, groupName, err := addSampleGroup(test.Manager, "my-group", []string{}, 5) // user-0 to user-4
	if err != nil {
		t.Fatalf("failed to add sample group: %v", err)
	}
	assertHasGroups(t, test.Manager, []string{groupName})

	group, err := test.Manager.GetGroup(&pb.GetGroupRequest{Group: groupName})
	if err != nil {
		t.Errorf("failed to add get group: %v", err)
	}
	if diff := cmp.Diff(users, group.GetMembers()); diff != "" {
		t.Errorf("got unexpected members for group %q: (-want +got):\n%s", groupName, diff)
	}

	// get a non-existent group
	nonExistentGroup := "i-dont-exist"
	if _, err := test.Manager.GetGroup(&pb.GetGroupRequest{Group: nonExistentGroup}); err == nil {
		t.Errorf("expected error getting non-existant group %q", nonExistentGroup)
	}
}

// TestAddGroupMember ...
func TestAddGroupMember(t *testing.T) {
	if skipGroupMemberTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// add two sample users and a new group with the first one as the initial member
	users, err := addSampleUsers(test.Manager, 2)
	if err != nil {
		t.Fatalf("failed to add sample users: %v", err)
	}
	_, groupName, err := addSampleGroup(test.Manager, "my-group", []string{users[0]}, 0)
	if err != nil {
		t.Fatalf("failed to add sample group: %v", err)
	}
	assertHasGroups(t, test.Manager, []string{groupName})

	// add user 1 to "my-group" as well
	if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: users[1], Group: groupName}); memberStatus.GetIsMember() {
		t.Errorf("expected user %q to be not yet a member of group %q: %v", users[1], groupName, err)
	}
	allowNonExistent := false
	if err := test.Manager.AddGroupMember(&pb.GroupMember{Group: groupName, Username: users[1]}, allowNonExistent); err != nil {
		t.Errorf("failed to add user %q to group %q: %v", users[1], groupName, err)
	}
	if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: users[1], Group: groupName}); !memberStatus.GetIsMember() {
		t.Errorf("expected user %q to have become a member of group %q: %v", users[1], groupName, err)
	}

	// add a non-existent member to an existing group
	// this will fail because the user is not present in the users group and strict checking is the default
	// however, if we were adding to the users group this would succeed
	nonExistantUser := "i-am-not-there"
	if err := test.Manager.AddGroupMember(&pb.GroupMember{Group: groupName, Username: nonExistantUser}, allowNonExistent); err == nil {
		t.Errorf("expected error adding a non-existent user %q to group %q", nonExistantUser, groupName)
	}

	// add an existent user to a non-existing group
	nonExistantGroup := "group-that-is-ficticious"
	if err := test.Manager.AddGroupMember(&pb.GroupMember{Group: nonExistantGroup, Username: users[0]}, allowNonExistent); err == nil {
		t.Errorf("expected error when adding user %q to a non-existent group %q", users[0], nonExistantGroup)
	}
}

// TestDeleteGroupMember ...
func TestDeleteGroupMember(t *testing.T) {
	if skipGroupMemberTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// make sure deleting from the default user group always fails
	nonExistantUser := "i-dont-exist"
	allowDeleteOfDefaultGroups := false
	if err := test.Manager.DeleteGroupMember(&pb.GroupMember{Group: test.Manager.DefaultUserGroup, Username: "someone"}, allowDeleteOfDefaultGroups); err == nil {
		t.Error("expected removing member of the users group to fail")
	}

	// create a new group with the first user as the initial member
	users, groupName, err := addSampleGroup(test.Manager, "my-group", []string{}, 2) // user-0 and user-1
	if err != nil {
		t.Fatalf("failed to add sample group: %v", err)
	}
	assertHasGroups(t, test.Manager, []string{groupName})

	// delete a non-existent member of an existing group
	if err := test.Manager.DeleteGroupMember(&pb.GroupMember{Group: groupName, Username: nonExistantUser}, allowDeleteOfDefaultGroups); err == nil {
		t.Errorf("expected error when deleting a non-existent user %q from group %q", nonExistantUser, groupName)
	}

	// delete an existent user from an non-existing group
	nonExistantGroup := "group-that-is-ficticious"
	if err := test.Manager.DeleteGroupMember(&pb.GroupMember{Group: nonExistantGroup, Username: users[0]}, allowDeleteOfDefaultGroups); err == nil {
		t.Errorf("expected error when deleting user %q from a non-existent group %q", users[0], nonExistantGroup)
	}

	// delete user 1
	if err := test.Manager.DeleteGroupMember(&pb.GroupMember{Group: groupName, Username: users[1]}, allowDeleteOfDefaultGroups); err != nil {
		t.Fatalf("failed to delete member %q of group %q: %v", users[1], groupName, err)
	}
	if memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: users[1], Group: groupName}); memberStatus.GetIsMember() {
		t.Errorf("expected user %q to be no longer a member of group %q: %v", users[1], groupName, err)
	}

	// make sure user 0 can not be deleted from "my-group" because it is the only remaining member
	if err := test.Manager.DeleteGroupMember(&pb.GroupMember{Group: groupName, Username: users[0]}, allowDeleteOfDefaultGroups); err == nil {
		t.Fatalf("failed to delete member %q of group %q: %v", users[0], groupName, err)
	}
}
