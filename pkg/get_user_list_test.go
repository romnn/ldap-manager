package pkg

import (
	"testing"

	"github.com/romnn/go-recursive-sort"
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
				Username:      req.GetUsername(),
				FirstName:     req.GetFirstName(),
				LastName:      req.GetLastName(),
				DisplayName:   "roman d",
				CN:            "roman d",
				Email:         req.GetEmail(),
				UID:           2001,
				GID:           2001,
				LoginShell:    "/bin/bash",
				HomeDirectory: "/home/test-user",
			},
			{
				Username:      test.Manager.DefaultAdminUsername,
				FirstName:     "changeme",
				LastName:      "changeme",
				DisplayName:   "changeme changeme",
				CN:            "changeme changeme",
				Email:         "changeme@changeme.com",
				UID:           2000,
				GID:           2001,
				LoginShell:    "/bin/bash",
				HomeDirectory: "/home/ldapadmin",
			},
		},
		Total: 2,
	}

	sort := recursivesort.RecursiveSort{StructSortField: "GID"}
	sort.Sort(&users)
	sort.Sort(&expected)

	t.Log(PrettyPrint(users))
	t.Log(PrettyPrint(expected))

	if equal, diff := EqualProto(expected, users); !equal {
		t.Fatalf("unexpected users: \n%s", diff)
	}
}
