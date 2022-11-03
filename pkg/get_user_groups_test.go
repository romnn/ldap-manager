package pkg

import (
	"testing"

	recursivesort "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestGetUserGroups tests getting the groups a user is member of
func TestGetUserGroups(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
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
	groupName := "test-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		t.Fatalf("failed to add new group: %v", err)
	}

	expected := &pb.GroupList{
		Groups: []*pb.Group{
			{
				Name: groupName,
				Members: []string{
					test.Manager.UserDN(username),
				},
				GID: 2002,
			},
			{
				Name: "users",
				Members: []string{
					"uid=ldapadmin,ou=users,dc=example,dc=org",
					test.Manager.UserDN(username),
				},
				GID: 2001,
			},
		},
		Total: 2,
	}

	groups, err := test.Manager.GetUserGroups(&pb.GetUserGroupsRequest{
		Username: username,
	})
	if err != nil {
		t.Fatalf(
			"failed to get the groups of user %q: %v",
			username, err,
		)
	}

	sort := recursivesort.RecursiveSort{StructSortField: "GID"}
	sort.Sort(&groups)
	sort.Sort(&expected)

	if equal, diff := EqualProto(expected, groups); !equal {
		t.Fatalf("unexpected groups: \n%s", diff)
	}
}
