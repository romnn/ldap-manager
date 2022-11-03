package pkg

import (
	"testing"

	"github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestUpdateUser tests updating a user
func TestUpdateUser(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// add a user
	username := "some-user"
	password := "Hallo Welt"
	if err := test.Manager.NewUser(&pb.NewUserRequest{
		Username:  username,
		Password:  password,
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
			"failed to add group %q with member %v: %v",
			groupName, username, err,
		)
	}

	// update the user (including the username and password)
	newUsername := "fancy-new-username"
	newPassword := "another password"
	isAdmin := false
	updatedUsername, err := test.Manager.UpdateUser(&pb.UpdateUserRequest{
		Username: username,
		Update: &pb.NewUserRequest{
			FirstName: "new first name",
			LastName:  "new last name",
			UID:       2300,
			GID:       2300,
			Username:  newUsername,
			Password:  newPassword,
		},
	}, isAdmin)
	if err != nil {
		t.Fatalf(
			"failed to update and rename user from %q to %q: %v",
			username, newUsername, err,
		)
	}
	if updatedUsername != newUsername {
		t.Fatalf(
			"expected updated username to be %q but got %q",
			newUsername, updatedUsername,
		)
	}

	// assert authentication with the new password succeeds
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: newUsername,
		Password: newPassword,
	}); err != nil {
		t.Fatalf(
			"authenticating user %q with new password %q failed: %v",
			newUsername, newPassword, err,
		)
	}

	// assert authentication with the old password fails
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: newUsername,
		Password: password,
	}); err == nil {
		t.Fatalf(
			"authenticating user %q with old password %q succeeded",
			newUsername, password,
		)
	}

	// assert the user was updated in all groups
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of groups: %v", err)
	}

	expected := &pb.GroupList{
		Groups: []*pb.Group{
			{
				Name: "admins",
				Members: []string{
					"uid=ldapadmin,ou=users,dc=example,dc=org",
				},
				GID: 2000,
			},
			{
				Name: groupName,
				Members: []string{
					test.Manager.UserNamed(newUsername),
				},
				GID: 2002,
			},
			{
				Name: "users",
				Members: []string{
					"uid=ldapadmin,ou=users,dc=example,dc=org",
					test.Manager.UserNamed(newUsername),
				},
				GID: 2001,
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
		t.Fatalf("unexpected group members: \n%s", diff)
	}
}
