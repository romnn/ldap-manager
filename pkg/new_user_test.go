package pkg

import (
	// "fmt"
	"testing"
	// "strconv"
	// "strings"
	// "github.com/google/go-cmp/cmp"
	"github.com/romnn/deepequal"
	// "github.com/romnn/go-recursive-sort"
	// "github.com/k0kubun/pp"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// ldaptest "github.com/romnn/ldap-manager/test"
	// ldaphash "github.com/romnn/ldap-manager/pkg/hash"
)

// func getAttribute(users *pb.UserList, attr string) ([]string, error) {
// 	results := []string{}
// 	for _, user := range users.GetUsers() {
// 		value, ok := user.GetData()[attr]
// 		if !ok {
// 			return results, fmt.Errorf("user %q has no attribute %q", user, attr)
// 		}
// 		results = append(results, value)
// 	}
// 	return results, nil
// }

// TestNewUser tests adding a new user
func TestNewUser(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	username := "romnn"
	req := pb.NewUserRequest{
		Account: &pb.Account{
			Username:  username,
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		},
	}
	expected := map[string]string{
		"uid":           "romnn",
		"givenName":     "roman",
		"displayName":   "roman d",
		"uidNumber":     "2001",
		"sn":            "d",
		"cn":            "roman d",
		"gidNumber":     "2002",
		"loginShell":    "/bin/bash",
		"homeDirectory": "/home/romnn",
		"mail":          "a@b.de",
	}

	if err := test.Manager.NewUser(&req, pb.HashingAlgorithm_DEFAULT); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}
	user, err := test.Manager.GetUser(username)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	// pls := user.(proto.Message)
	// t.Log(pls.DebugString())
	t.Log(PrettyPrint(user))
	// // sort := recursivesort.RecursiveSort{}
	// // sort.Sort(&received)
	// // sort.Sort(&expected)

	// t.Log(received)
	// t.Log(expected)
	if equal, err := deepequal.DeepEqual(user.GetData(), expected); !equal {
		t.Fatalf("unexpected user data: %v", err)
	}

	// try to bind as the user

	// todo: try to login / bind as that user
	// todo: extend the types in their respective packages by placing files next to them

	// t.Fatalf("todo")
	// // add two valid users
	// expected := []string{"romnn", "uwe12"}
	// requests := []*pb.NewUserRequest{
	// 	{
	// 		Account: &pb.Account{
	// 			Username:  expected[0],
	// 			Password:  "Hallo Welt",
	// 			Email:     "a@b.de",
	// 			FirstName: "roman",
	// 			LastName:  "d",
	// 		},
	// 	},
	// 	{
	// 		Account: &pb.Account{
	// 			Username:  expected[1],
	// 			Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
	// 			Email:     "uwe-h@mobile.com",
	// 			FirstName: "uwe",
	// 			LastName:  "Heisenberg",
	// 		},
	// 	},
	// }
	// for _, req := range requests {
	// 	if err := test.Manager.NewUser(req, pb.HashingAlgorithm_DEFAULT); err != nil {
	// 		t.Errorf("failed to add user: %v", err)
	// 	}
	// }

	// user, err := test.Manager.GetUser(username{})
	// if err != nil {
	// 	t.Fatalf("failed to get user: %v", err)
	// }
	// t.Log(user)
	// t.Fatalf("todo")
	// users, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	// if err != nil {
	// 	t.Fatalf("failed to get users: %v", err)
	// }
	// received, err := getAttribute(users, test.Manager.AccountAttribute)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// // sort := recursivesort.RecursiveSort{}
	// // sort.Sort(&received)
	// // sort.Sort(&expected)

	// t.Log(received)
	// t.Log(expected)
	// if equal, err := deepequal.DeepEqual(received, expected); !equal {
	// 	t.Fatal(err)
	// }
}

// TestNewUserValidation ...
func TestNewUserValidation(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	cases := []struct {
		valid   bool
		request *pb.NewUserRequest
	}{
		// invalid: missing everything
		{false, &pb.NewUserRequest{}},
		// invalid: missing username
		{false, &pb.NewUserRequest{
			Account: &pb.Account{
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "d",
			},
		}},
		// invalid: missing password
		{false, &pb.NewUserRequest{
			Account: &pb.Account{
				Username:  "peter1",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "d",
			},
		}},
		// invalid: missing email
		{false, &pb.NewUserRequest{
			Account: &pb.Account{
				Username:  "peter2",
				Password:  "Hallo Welt",
				FirstName: "roman",
				LastName:  "d",
			},
		}},
		// invalid: missing first name
		{false, &pb.NewUserRequest{
			Account: &pb.Account{
				Username: "peter3",
				Password: "Hallo Welt",
				Email:    "a@b.de",
				LastName: "d",
			},
		}},
		// invalid: missing last name
		{false, &pb.NewUserRequest{
			Account: &pb.Account{
				Username:  "peter4",
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
			},
		}},
		// valid: all required fields
		{true, &pb.NewUserRequest{
			Account: &pb.Account{
				Username:  "peter5",
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "test",
			},
		}},
		// invalid: email is not valid
		{false, &pb.NewUserRequest{
			Account: &pb.Account{
				Username:  "peter5",
				Password:  "Hallo Welt",
				Email:     "test.de",
				FirstName: "roman",
				LastName:  "test",
			},
		}},
	}
	for _, c := range cases {
		err := test.Manager.NewUser(c.request, pb.HashingAlgorithm_DEFAULT)
		if err != nil && c.valid {
			t.Errorf("failed to add valid user: %v", err)
		}
		if err == nil && !c.valid {
			t.Errorf("expected error when adding invalid user %v", c.request)
		}
	}
}

// func containsUsers(observed *pb.UserList, expected []string, attr string) error {
// 	for _, e := range expected {
// 		found := false
// 		for _, o := range observed.GetUsers() {
// 			if uid, ok := o.GetData()[attr]; ok && uid == e {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			return fmt.Errorf("expected user %q after it was added but only got %v", e, observed)
// 		}
// 	}
// 	return nil
// }
