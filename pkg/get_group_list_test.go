package pkg

import (
	"testing"

	recursivesort "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestGetGroupList tests getting a list of all groups
func TestGetGroupList(t *testing.T) {
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

	// get all groups
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of groups: %v", err)
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
						// Dn: "uid=ldapadmin,ou=users,dc=example,dc=org",
					},
					{
						Username: username,
						Dn:       test.Manager.UserDN(username),
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
				Name: groupName,
				Members: []*pb.GroupMember{
					{
						Username: username,
						Dn:       test.Manager.UserDN(username),
						Group:    groupName,
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
		t.Fatalf("unexpected groups: \n%s", diff)
	}
}
