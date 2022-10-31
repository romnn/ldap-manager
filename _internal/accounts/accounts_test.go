package accounts

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	ldaphash "github.com/romnn/ldap-manager/pkg/hash"
)

func contains(list []string, a string) bool {
	for _, b := range list {
		if strings.ToLower(b) == strings.ToLower(a) {
			return true
		}
	}
	return false
}

func containsUsers(observed *pb.UserList, expected []string, attr string) error {
	for _, e := range expected {
		found := false
		for _, o := range observed.GetUsers() {
			if uid, ok := o.GetData()[attr]; ok && uid == e {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected to find user %q after it was added but only got %v", e, observed)
		}
	}
	return nil
}

// TestAddNewUserAndGetUserList ...
func TestAddNewUserAndGetUserList(t *testing.T) {
	if skipAccountTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	// Add two valid users
	expected := []string{"romnn", "uwe12"}
	users := []*pb.NewAccountRequest{
		{
			Account: &pb.Account{
				Username:  expected[0],
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "d",
			},
		},
		{
			Account: &pb.Account{
				Username:  expected[1],
				Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
				Email:     "uwe-h@mobile.com",
				FirstName: "uwe",
				LastName:  "Heisenberg",
			},
		},
	}
	for _, newUserReq := range users {
		if err := test.Manager.NewAccount(newUserReq, pb.HashingAlgorithm_DEFAULT); err != nil {
			t.Errorf("failed to add user: %v", err)
		}
	}

	// List all users
	userList, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Errorf("failed to get users list: %v", err)
	}
	if err := containsUsers(userList, expected, test.Manager.AccountAttribute); err != nil {
		t.Error(err)
	}
}

// TestNewAccountValidation ...
func TestNewAccountValidation(t *testing.T) {
	if skipAccountTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	cases := []struct {
		valid   bool
		request *pb.NewAccountRequest
	}{
		// invalid: missing everything
		{false, &pb.NewAccountRequest{}},
		// invalid: missing username
		{false, &pb.NewAccountRequest{
			Account: &pb.Account{
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "d",
			},
		}},
		// invalid: missing password
		{false, &pb.NewAccountRequest{
			Account: &pb.Account{
				Username:  "peter1",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "d",
			},
		}},
		// invalid: missing email
		{false, &pb.NewAccountRequest{
			Account: &pb.Account{
				Username:  "peter2",
				Password:  "Hallo Welt",
				FirstName: "roman",
				LastName:  "d",
			},
		}},
		// invalid: missing first name
		{false, &pb.NewAccountRequest{
			Account: &pb.Account{
				Username: "peter3",
				Password: "Hallo Welt",
				Email:    "a@b.de",
				LastName: "d",
			},
		}},
		// invalid: missing last name
		{false, &pb.NewAccountRequest{
			Account: &pb.Account{
				Username:  "peter4",
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
			},
		}},
		// valid: all required fields
		{true, &pb.NewAccountRequest{
			Account: &pb.Account{
				Username:  "peter5",
				Password:  "Hallo Welt",
				Email:     "a@b.de",
				FirstName: "roman",
				LastName:  "test",
			},
		}},
		// invalid: email is not valid
		{false, &pb.NewAccountRequest{
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
		err := test.Manager.NewAccount(c.request, pb.HashingAlgorithm_DEFAULT)
		if err != nil && c.valid {
			t.Errorf("failed to add valid user: %v", err)
		}
		if err == nil && !c.valid {
			t.Errorf("expected error when adding invalid user %v", c.request)
		}
	}
}
