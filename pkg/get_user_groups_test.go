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
				Name: "users",
				Members: []*pb.GroupMember{
					{
						Dn:       test.Manager.UserDN("ldapadmin"),
						Username: "ldapadmin",
						Group:    "users",
					},
					{
						Dn:       test.Manager.UserDN(username),
						Username: username,
						Group:    "users",
					},
				},
				GID: 2000,
			},
			{
				Name: groupName,
				Members: []*pb.GroupMember{
					{
						Dn:       test.Manager.UserDN(username),
						Username: username,
						Group:    groupName,
					},
				},
				GID: 2002,
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

	sort := recursivesort.RecursiveSort{}
	sort.Sort(&groups)
	sort.Sort(&expected)

	if equal, diff := EqualProto(expected, groups); !equal {
		t.Fatalf("unexpected groups: \n%s", diff)
	}
}
