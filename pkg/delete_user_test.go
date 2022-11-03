package pkg

import (
	"testing"

	"github.com/romnn/deepequal"
	"github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

func getUserListUsernames(users *pb.UserList) []string {
	var usernames []string
	for _, user := range users.GetUsers() {
		usernames = append(usernames, user.GetUsername())
	}
	return usernames
}

// TestDeleteUser tests deleting users
func TestDeleteUser(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	usernames := []string{"user1", "user2"}
	for _, username := range usernames {
		if err := test.Manager.NewUser(&pb.NewUserRequest{
			Username:  username,
			Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
			Email:     "felix@web.de",
			FirstName: "Felix",
			LastName:  "Heisenberg",
		}); err != nil {
			t.Fatalf("failed to add user: %v", err)
		}
	}

	// assert we find those two users
	users, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get user list: %v", err)
	}
	found := getUserListUsernames(users)
	expected := append(usernames, test.Manager.DefaultAdminUsername)
	recursivesort.Sort(&found)
	recursivesort.Sort(&expected)
	if equal, err := deepequal.DeepEqual(found, expected); !equal {
		t.Fatalf("unexpected users: %v", err)
	}

	// delete the first user
	keepGroups := false
	if err := test.Manager.DeleteUser(&pb.DeleteUserRequest{
		Username: usernames[0],
	}, keepGroups); err != nil {
		t.Fatalf(
			"failed to delete user %q: %v",
			usernames[0], err,
		)
	}

	// assert we find only the second user
	users, err = test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get user list: %v", err)
	}
	found = getUserListUsernames(users)
	expected = append(usernames[1:], test.Manager.DefaultAdminUsername)
	recursivesort.Sort(&found)
	recursivesort.Sort(&expected)
	if equal, err := deepequal.DeepEqual(found, expected); !equal {
		t.Fatalf("unexpected users: %v", err)
	}
}
