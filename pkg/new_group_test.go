package pkg

import (
	"testing"

	"github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestNewGroupEmpty tests that adding an empty group fails
func TestNewGroupEmpty(t *testing.T) {
	test := new(Test).Setup(t)
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
	test := new(Test).Setup(t)
	defer test.Teardown()
	username := "this-user-does-not-exist"
	groupName := "my-group-with-non-existent-members"

	// adding a group with missing members should fail in strict mode
	strict := true
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err == nil {
		t.Fatalf("expected error adding group %q with missing member %q (strict=%t)", groupName, username, strict)
	}

	// adding a group with missing members is allowed if not strict
	strict = false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		t.Fatalf("failed to add group %q with missing member %q (strict=%t)", groupName, username, strict)
	}
}

// TestNewGroup tests adding a valid new group
func TestNewGroup(t *testing.T) {
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

	// get all groups
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of groups: %v", err)
	}

	expected := []*pb.Group{
		{
			Name: "admins",
			Members: []string{
				"uid=ldapadmin,ou=users,dc=example,dc=org",
			},
			GID: 2000,
		},
		{
			Name: "my-group",
			Members: []string{
				"uid=some-user,ou=users,dc=example,dc=org",
			},
			GID: 2002,
		},
		{
			Name: "users",
			Members: []string{
				"uid=ldapadmin,ou=users,dc=example,dc=org",
				"uid=some-user,ou=users,dc=example,dc=org",
			},
			GID: 2001,
		},
	}

	sort := recursivesort.RecursiveSort{StructSortField: "GID"}
	sort.Sort(&groups)
	sort.Sort(&expected)

	t.Log(PrettyPrint(groups))
	t.Log(PrettyPrint(expected))

	for i := range expected {
		if equal, diff := EqualProto(expected[i], groups.GetGroups()[i]); !equal {
			t.Fatalf("unexpected group at index %d: \n%s", i, diff)
		}
	}
}
