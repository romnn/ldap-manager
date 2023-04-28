package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestGetGroup tests getting a group
func TestGetGroup(t *testing.T) {
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
		t.Fatalf("failed to add new user: %v", err)
	}

	// add the user to a group
	strict := true
	groupName := "my-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{username},
	}, strict); err != nil {
		t.Fatalf(
			"failed to create group %q with member %v: %v",
			groupName, username, err,
		)
	}

	expected := &pb.Group{
		Name: groupName,
		Members: []*pb.GroupMember{
			{
				Username: username,
				Dn:       test.Manager.UserDN(username),
				Group:    groupName,
			},
		},
		GID: 2002,
	}
	group, err := test.Manager.GetGroupByName(groupName)

	t.Log(PrettyPrint(expected))
	t.Log(PrettyPrint(group))

	if err != nil {
		t.Fatalf(
			"failed to get group %q: %v",
			groupName, err,
		)
	}
	if equal, diff := EqualProto(expected, group); !equal {
		t.Fatalf(
			"unexpected group %q: \n%s",
			groupName, diff,
		)
	}
}

// TestGetDefaultGroup tests getting the default groups
func TestGetDefaultGroup(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	for _, c := range []struct {
		name     string
		expected *pb.Group
	}{
		{
			name: test.Manager.DefaultUserGroup,
			expected: &pb.Group{
				Name: "users",
				Members: []*pb.GroupMember{
					{
						Username: "ldapadmin",
						Dn:       test.Manager.UserDN("ldapadmin"),
						Group:    "users",
					},
				},
				GID: 2000,
			},
		},

		{
			name: test.Manager.DefaultAdminGroup,
			expected: &pb.Group{
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
		},
	} {
		group, err := test.Manager.GetGroupByName(c.name)

		t.Log(PrettyPrint(c.expected))
		t.Log(PrettyPrint(group))

		if err != nil {
			t.Fatalf(
				"failed to get group %q: %v",
				c.name, err,
			)
		}
		if equal, diff := EqualProto(c.expected, group); !equal {
			t.Fatalf(
				"unexpected group %q: \n%s",
				c.name, diff,
			)
		}
	}
}
