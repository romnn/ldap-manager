package pkg

import (
	// "fmt"
	"testing"
	// "strconv"
	// "strings"
	// "github.com/google/go-cmp/cmp"
	// "github.com/romnn/deepequal"
	// "github.com/romnn/go-recursive-sort"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// ldaptest "github.com/romnn/ldap-manager/test"
	// ldaphash "github.com/romnn/ldap-manager/pkg/hash"
)

// TestGetGroupList tests getting a list of all groups
func TestGetGroupList(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	// username := "romnn"
	// req := pb.NewUserRequest{
	// 	Account: &pb.Account{
	// 		Username:  username,
	// 		Password:  "Hallo Welt",
	// 		Email:     "a@b.de",
	// 		FirstName: "roman",
	// 		LastName:  "d",
	// 	},
	// }
	// expected := map[string]string{
	// 	"uid":           "romnn",
	// 	"givenName":     "roman",
	// 	"displayName":   "roman d",
	// 	"uidNumber":     "2001",
	// 	"sn":            "d",
	// 	"cn":            "roman d",
	// 	"gidNumber":     "2002",
	// 	"loginShell":    "/bin/bash",
	// 	"homeDirectory": "/home/romnn",
	// 	"mail":          "a@b.de",
	// }

	// if err := test.Manager.NewUser(&req, pb.HashingAlgorithm_DEFAULT); err != nil {
	// 	t.Fatalf("failed to add user: %v", err)
	// }
	// user, err := test.Manager.GetUser(username)
	// if err != nil {
	// 	t.Fatalf("failed to get user: %v", err)
	// }
	username := "test-user"
	req := pb.NewUserRequest{
		Account: &pb.Account{
			Username:  username,
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		},
	}

	// create a new user (will create the user group)
	if err := test.Manager.NewUser(&req, pb.HashingAlgorithm_DEFAULT); err != nil {
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

	t.Log(PrettyPrint(groups))
}
