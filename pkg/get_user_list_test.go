package pkg

import (
	"fmt"
	"testing"

	recursivesort "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestGetUserList tests getting a list of all users
func TestGetUserList(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	req := pb.NewUserRequest{
		Username:  "test-user",
		Password:  "Hallo Welt",
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}

	// create a new user (will create the user group)
	if err := test.Manager.NewUser(&req); err != nil {
		t.Fatalf("failed to add new user: %v", err)
	}

	// get all users
	users, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of users: %v", err)
	}

	expected := &pb.UserList{
		Users: []*pb.User{
			{
				Username:    test.Manager.DefaultAdminUsername,
				FirstName:   "changeme",
				LastName:    "changeme",
				DisplayName: "changeme changeme",
				CN:          "changeme changeme",
				DN: fmt.Sprintf(
					"uid=%s,ou=users,dc=example,dc=org",
					test.Manager.DefaultAdminUsername,
				),
				Email:         "changeme@changeme.com",
				UID:           2000,
				GID:           2000,
				LoginShell:    "/bin/bash",
				HomeDirectory: "/home/ldapadmin",
			},
			{
				Username:    req.GetUsername(),
				FirstName:   req.GetFirstName(),
				LastName:    req.GetLastName(),
				DisplayName: "roman d",
				CN:          "roman d",
				DN: fmt.Sprintf(
					"uid=%s,ou=users,dc=example,dc=org",
					req.GetUsername(),
				),
				Email:         req.GetEmail(),
				UID:           2001,
				GID:           2000,
				LoginShell:    "/bin/bash",
				HomeDirectory: "/home/test-user",
			},
		},
		Total: 2,
	}

	sort := recursivesort.RecursiveSort{StructSortField: "UID"}
	sort.Sort(&users)
	sort.Sort(&expected)

	t.Log(PrettyPrint(users))
	t.Log(PrettyPrint(expected))

	if equal, diff := EqualProto(expected, users); !equal {
		t.Fatalf("unexpected users: \n%s", diff)
	}
}
