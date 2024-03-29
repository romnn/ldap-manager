package pkg

import (
	"testing"

	recursivesort "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestNewGroupEmpty tests that adding an empty group fails
func TestNewGroupEmpty(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// adding a group with no members should fail
	strict := false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    "my-group",
		Members: []string{},
	}, strict); err == nil {
		t.Fatalf("expected error when adding a group with no members")
	}
}

// TestNewGroupMissingMember tests that adding an group with missing members fails
func TestNewGroupMissingMember(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()
	username := "this-user-does-not-exist"
	groupName := "my-group-with-non-existent-members"

	// adding a group with missing members should fail in strict mode
	strict := true
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err == nil {
		t.Fatalf(
			"expected error adding group %q with missing member %q (strict=%t)",
			groupName, username, strict,
		)
	}

	// adding a group with missing members is allowed if not strict
	strict = false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		t.Fatalf(
			"failed to add group %q with missing member %q (strict=%t)",
			groupName, username, strict,
		)
	}
}

// TestNewGroup tests adding a valid new group
func TestNewGroup(t *testing.T) {
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
		t.Fatalf(
			"failed to add new user: %v",
			err,
		)
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

	// get all groups
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf(
			"failed to get list of groups: %v",
			err,
		)
	}

	expected := &pb.GroupList{
		Groups: []*pb.Group{
			{
				Name: "users",
				Members: []*pb.GroupMember{
					{
						Username: "ldapadmin",
						Dn:       test.Manager.UserDN("ldapadmin"),
						Group:    "users",
					},
					{
						Username: "some-user",
						Dn:       test.Manager.UserDN("some-user"),
						Group:    "users",
					},
				},
				GID: 2000,
			},
			{
				Name: "admins",
				Members: []*pb.GroupMember{
					{
						Username: "ldapadmin",
						Dn:       test.Manager.UserDN("ldapadmin"),
						Group:    "admins",
					},
				},
				GID: 2001,
			},
			{
				Name: "my-group",
				Members: []*pb.GroupMember{
					{
						Username: "some-user",
						Dn:       test.Manager.UserDN("some-user"),
						Group:    "my-group",
					},
				},
				GID: 2002,
			},
		},
		Total: 3,
	}

	sort := recursivesort.RecursiveSort{StructSortField: "GID"}
	sort.Sort(&groups)
	sort.Sort(&expected)

	t.Log(PrettyPrint(groups))
	t.Log(PrettyPrint(expected))

	if equal, diff := EqualProto(expected, groups); !equal {
		t.Fatalf("unexpected group: \n%s", diff)
	}
}
