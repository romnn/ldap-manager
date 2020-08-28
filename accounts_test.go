package ldapmanager

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	ldaphash "github.com/romnnn/ldap-manager/hash"
)

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
			Username:  expected[0],
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		},
		{
			Username:  expected[1],
			Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
			Email:     "uwe-h@mobile.com",
			FirstName: "uwe",
			LastName:  "Heisenberg",
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

// TestAuthenticateUser ...
func TestAuthenticateUser(t *testing.T) {
	if skipAccountTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	samplePasswords := []string{"123456", "Hallo@Welt", "@#73sAdf0^E^RC#+++83230*###$&"}
	for _, algorithm := range ldaphash.LDAPPasswordHashingAlgorithms {
		name, _ := pb.HashingAlgorithm_name[int32(algorithm)]
		for _, pw := range samplePasswords {
			// t.Log(name, algorithm, pw)
			var finalErr error
			attemptsLeft := 5
			for {
				// FIXME: this tests is flaky :(
				attemptsLeft--
				newUserReq := &pb.NewAccountRequest{
					Username:  name + pw + strconv.Itoa(attemptsLeft),
					Password:  pw,
					Email:     "a@b.de",
					FirstName: "roman",
					LastName:  "d",
				}
				if err := test.Manager.NewAccount(newUserReq, algorithm); err != nil {
					if attemptsLeft <= 0 {
						finalErr = fmt.Errorf("failed to add user %q: %v", newUserReq.Username, err)
						break
					}
					continue
				}

				// now check if we can authenticate using the clear password
				if err := test.Manager.AuthenticateUser(&pb.AuthenticateUserRequest{Username: newUserReq.Username, Password: pw}); err != nil {
					if attemptsLeft <= 0 {
						finalErr = fmt.Errorf("failed to authenticate user %q with password %q: %v", newUserReq.Username, pw, err)
						break
					}
					continue
				}
				break
			}
			if finalErr != nil {
				t.Error(finalErr)
			}
		}
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
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing password
		{false, &pb.NewAccountRequest{
			Username:  "peter1",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing email
		{false, &pb.NewAccountRequest{
			Username:  "peter2",
			Password:  "Hallo Welt",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing first name
		{false, &pb.NewAccountRequest{
			Username: "peter3",
			Password: "Hallo Welt",
			Email:    "a@b.de",
			LastName: "d",
		}},
		// invalid: missing last name
		{false, &pb.NewAccountRequest{
			Username:  "peter4",
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
		}},
		// valid: all required fields
		{true, &pb.NewAccountRequest{
			Username:  "peter5",
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "test",
		}},
		// invalid: email is not valid
		{false, &pb.NewAccountRequest{
			Username:  "peter5",
			Password:  "Hallo Welt",
			Email:     "test.de",
			FirstName: "roman",
			LastName:  "test",
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

// TestGetAccount ...
func TestGetAccount(t *testing.T) {
	if skipAccountTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	newUserReq := &pb.NewAccountRequest{
		Username:  "felix",
		Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
		Email:     "felix@web.de",
		FirstName: "Felix",
		LastName:  "Heisenberg",
	}
	if err := test.Manager.NewAccount(newUserReq, pb.HashingAlgorithm_DEFAULT); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// Make sure the users group was created
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of all groups: %v", err)
	}
	if groups.GetGroups()[0] != test.Manager.DefaultUserGroup {
		t.Fatalf("expected the default user group %q to have been created", test.Manager.DefaultUserGroup)
	}

	// Make sure that the new account is in the users group
	group, err := test.Manager.GetGroup(&pb.GetGroupRequest{Name: test.Manager.DefaultUserGroup})
	if err != nil {
		t.Fatalf("failed to get members of the group %q: %v", test.Manager.DefaultUserGroup, err)
	}
	if group.Members[0] != newUserReq.Username {
		t.Fatalf("expected the new user %q to be a member of the default user group %q", newUserReq.Username, test.Manager.DefaultUserGroup)
	}

	memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: newUserReq.Username, Group: test.Manager.DefaultUserGroup})
	if err != nil {
		t.Fatalf("failed to check if user %q is in the group %q: %v", newUserReq.Username, test.Manager.DefaultUserGroup, err)
	}
	if !memberStatus.GetIsMember() {
		t.Fatalf("expected user %q to be a member of the group %q: %v", newUserReq.Username, test.Manager.DefaultUserGroup, err)
	}

	account, err := test.Manager.GetAccount(&pb.GetAccountRequest{Username: "felix"})
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
	expected := map[string]string{
		"cn":            "Felix Heisenberg",
		"displayName":   "Felix Heisenberg",
		"gidNumber":     "2001",
		"givenName":     "Felix",
		"homeDirectory": "/home/felix",
		"loginShell":    "/bin/bash",
		"mail":          "felix@web.de",
		"sn":            "Heisenberg",
		"uid":           "felix",
		"uidNumber":     "2001",
	}
	if diff := cmp.Diff(expected, account.GetData()); diff != "" {
		t.Errorf("got unexpected account result: (-want +got):\n%s", diff)
	}
}

// TestDeleteAccount ...
func TestDeleteAccount(t *testing.T) {
	if skipAccountTests {
		t.Skip()
	}
	test := new(Test).Setup(t)
	defer test.Teardown()

	users := []string{"user1", "user2"}
	for _, user := range users {
		if err := test.Manager.NewAccount(&pb.NewAccountRequest{
			Username:  user,
			Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
			Email:     "felix@web.de",
			FirstName: "Felix",
			LastName:  "Heisenberg",
		}, pb.HashingAlgorithm_DEFAULT); err != nil {
			t.Fatalf("failed to add user: %v", err)
		}
	}

	// Assert we find those two users
	userList, err := test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get users list: %v", err)
	}
	if err := containsUsers(userList, users, test.Manager.AccountAttribute); err != nil {
		t.Error(err)
	}

	// Now delete the first user
	keepGroups := false
	if err := test.Manager.DeleteAccount(&pb.DeleteAccountRequest{Username: users[0]}, keepGroups); err != nil {
		t.Fatalf("failed to delete user %q: %v", users[0], err)
	}

	// Assert we find only the second user
	userList, err = test.Manager.GetUserList(&pb.GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get users list: %v", err)
	}
	if err := containsUsers(userList, users[1:2], test.Manager.AccountAttribute); err != nil {
		t.Error(err)
	}
}
