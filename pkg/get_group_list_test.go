package pkg

import (
	"testing"

	"github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestGetGroupList tests getting a list of all groups
func TestGetGroupList(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	username := "test-user"
	req := pb.NewUserRequest{
		Username:  username,
		Password:  "Hallo Welt",
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}

	// create a new user (will create the user group)
	if err := test.Manager.NewUser(&req); err != nil {
		t.Fatalf("failed to add new user: %v", err)
	}

	// create a new group
	strict := false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    "test-group",
		Members: []string{"test-user"},
	}, strict); err != nil {
		t.Fatalf("failed to add new group: %v", err)
	}

	// get all groups
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of groups: %v", err)
	}

	expected := []*pb.Group{
		&pb.Group{
			Name: "admins",
			Members: []string{
				"uid=ldapadmin,ou=users,dc=example,dc=org",
			},
			GID: 2000,
		},
		&pb.Group{
			Name: "test-group",
			Members: []string{
				"uid=test-user,ou=users,dc=example,dc=org",
			},
			GID: 2002,
		},
		&pb.Group{
			Name: "users",
			Members: []string{
				"uid=ldapadmin,ou=users,dc=example,dc=org",
				"uid=test-user,ou=users,dc=example,dc=org",
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
