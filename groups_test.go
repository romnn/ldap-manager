package ldapmanager

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func addSampleUsers(manager *LDAPManager, num int) ([]string, error) {
	var added []string
	for n := 0; n < num; n++ {
		username := fmt.Sprintf("user-%d", n)
		if err := manager.NewAccount(&NewAccountRequest{
			Username:  username,
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}); err != nil {
			return added, err
		}
		added = append(added, username)
	}
	return added, nil
}

func addSampleGroup(manager *LDAPManager, name string, members []string, num int) ([]string, string, error) {
	var err error
	if len(members) < 1 {
		members, err = addSampleUsers(manager, num)
		if err != nil {
			return members, "", err
		}
		if name == "" {
			name = "my-group"
		}
	}
	if err := manager.NewGroup(&NewGroupRequest{
		Name:    name,
		Members: members,
	}); err != nil {
		return members, "", err
	}
	return members, name, nil
}

func assertHasGroups(t *testing.T, manager *LDAPManager, expected []string) {
	groups, err := manager.GetGroupList(&GetGroupListRequest{})
	if err != nil {
		t.Errorf("failed to get groups: %v", err)
	}
	if diff := cmp.Diff(append(expected, manager.DefaultUserGroup), groups); diff != "" {
		t.Errorf("got unexpected groups: (-want +got):\n%s", diff)
	}
}

// TestNewGroup ...
func TestNewGroup(t *testing.T) {
	if skipGroupTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// add sample users (in the users group by default)
	validUsers, err := addSampleUsers(test.Manager, 2)
	if err != nil {
		t.Fatalf("failed to add sample users: %v", err)
	}

	// adding a group with no members should fail
	if err := test.Manager.NewGroup(&NewGroupRequest{
		Name:    "my-group",
		Members: []string{},
	}); err == nil {
		t.Error("expected error when adding a group with no members")
	}

	// adding a group with no members should fail when strict checking is enabled
	nonexistentUser := "this-user-does-not-exist"
	nonexistentMembersGroupName := "my-group-with-non-existent-members"
	if err := test.Manager.NewGroup(&NewGroupRequest{
		Name:    nonexistentMembersGroupName,
		Members: []string{nonexistentUser},
		Strict:  true,
	}); err == nil {
		t.Error("expected error when adding a group with a member that does not exist")
	}

	// adding a group with no members is possible when strict checking is disabled
	if err := test.Manager.NewGroup(&NewGroupRequest{
		Name:    nonexistentMembersGroupName,
		Members: []string{nonexistentUser},
		Strict:  false,
	}); err != nil {
		t.Error("failed to add group with nonexistent members and strict=false")
	}

	// now add a valid group with existing members
	validGroupName := "my-group"
	if err := test.Manager.NewGroup(&NewGroupRequest{
		Name:    validGroupName,
		Members: validUsers,
	}); err != nil {
		t.Errorf("failed to add group %q with existing members %v: %v", validGroupName, validUsers, err)
	}

	// make sure that we find the new groups and the default users group
	assertHasGroups(t, test.Manager, []string{validGroupName, nonexistentMembersGroupName})
}

// TestDeleteGroup ...
func TestDeleteGroup(t *testing.T) {
	if skipGroupTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// make sure deleting a non-existent group failed
	if err := test.Manager.DeleteGroup("group-that-does-not-exist"); err == nil {
		t.Error("expected error deleting group that does not exist")
	}

	_, groupName, err := addSampleGroup(test.Manager, "my-group", []string{}, 3)
	if err != nil {
		t.Fatalf("failed to add sample group: %v", err)
	}
	assertHasGroups(t, test.Manager, []string{groupName})

	// now delete the group
	if err := test.Manager.DeleteGroup(groupName); err != nil {
		t.Errorf("failed to delete group %q: %v", groupName, err)
	}
	assertHasGroups(t, test.Manager, []string{})

	// make sure deleting the users group is not allowed
	if err := test.Manager.DeleteGroup(test.Manager.DefaultUserGroup); err == nil {
		t.Error("expected error deleting the default users group")
	}
}

// TestRenameGroup ...
func TestRenameGroup(t *testing.T) {
	if skipGroupTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	_, groupName, err := addSampleGroup(test.Manager, "my-group", []string{}, 3)
	if err != nil {
		t.Fatalf("failed to add sample group: %v", err)
	}
	assertHasGroups(t, test.Manager, []string{groupName})
	groupBefore, err := test.Manager.GetGroup(&GetGroupRequest{Group: groupName})
	if err != nil {
		t.Fatalf("failed to get the group %q before rename: %v", groupName, err)
	}

	// Rename
	renamedGroupName := "my-renamed-group"
	if err := test.Manager.RenameGroup(&RenameGroupRequest{Group: groupName, NewName: renamedGroupName}); err != nil {
		t.Fatalf("failed to rename group from %q to %q: %v", groupName, renamedGroupName, err)
	}
	assertHasGroups(t, test.Manager, []string{renamedGroupName})

	// make sure members are left untouched
	groupAfter, err := test.Manager.GetGroup(&GetGroupRequest{Group: renamedGroupName})
	if err != nil {
		t.Fatalf("failed to get the renamed group %q before rename: %v", renamedGroupName, err)
	}
	if diff := cmp.Diff(groupBefore.Members, groupAfter.Members); diff != "" {
		t.Errorf("got different group members after rename: (-want +got):\n%s", diff)
	}

	// make sure the old name is really gone
	if _, err := test.Manager.GetGroup(&GetGroupRequest{Group: groupName}); err == nil {
		t.Errorf("expected error getting the renamed group by the old name %q", groupName)
	}
}
