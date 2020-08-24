package ldapmanager

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	ldaphash "github.com/romnnn/ldap-manager/hash"
	ldaptest "github.com/romnnn/ldap-manager/test"
)

func containsUsers(observed []map[string]string, expected []string, attr string) error {
	for _, e := range expected {
		found := false
		for _, o := range observed {
			if uid, ok := o[attr]; ok && uid == e {
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
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}

	// Add two valid users
	expected := []string{"romnn", "uwe12"}
	users := []*NewAccountRequest{
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
		if err := manager.NewAccount(newUserReq); err != nil {
			t.Errorf("failed to add user: %v", err)
		}
	}

	// List all users
	userList, err := manager.GetUserList(&GetUserListRequest{})
	if err != nil {
		t.Errorf("failed to get users list: %v", err)
	}
	if err := containsUsers(userList, expected, manager.AccountAttribute); err != nil {
		t.Error(err)
	}
}

// TestAuthenticateUser ...
func TestAuthenticateUser(t *testing.T) {
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}
	samplePasswords := []string{"123456", "Hallo@Welt", "@#73sAdf0^E^RC#+++83230*###$&"}
	for name, algorithm := range ldaphash.LDAPPasswordHashingAlgorithms {
		for _, pw := range samplePasswords {
			// t.Log(name, algorithm, pw)
			newUserReq := &NewAccountRequest{
				Username:         name + pw,
				Password:         pw,
				Email:            "a@b.de",
				FirstName:        "roman",
				LastName:         "d",
				HashingAlgorithm: algorithm,
			}
			if err := manager.NewAccount(newUserReq); err != nil {
				t.Errorf("failed to add user %q: %v", newUserReq.Username, err)
				continue
			}
			// wait some time to process the password hash
			time.Sleep(1 * time.Second)

			// now check if we can authenticate using the clear password
			if _, err := manager.AuthenticateUser(newUserReq.Username, pw); err != nil {
				t.Errorf("failed to authenticate user %q with password %q: %v", newUserReq.Username, pw, err)
			}
		}
	}
}

// TestNewAccountValidation ...
func TestNewAccountValidation(t *testing.T) {
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		valid   bool
		request *NewAccountRequest
	}{
		// invalid: missing everything
		{false, &NewAccountRequest{}},
		// invalid: missing username
		{false, &NewAccountRequest{
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing password
		{false, &NewAccountRequest{
			Username:  "peter1",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing email
		{false, &NewAccountRequest{
			Username:  "peter2",
			Password:  "Hallo Welt",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing first name
		{false, &NewAccountRequest{
			Username: "peter3",
			Password: "Hallo Welt",
			Email:    "a@b.de",
			LastName: "d",
		}},
		// invalid: missing last name
		{false, &NewAccountRequest{
			Username:  "peter4",
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
		}},
		// valid: all required fields
		{true, &NewAccountRequest{
			Username:  "peter5",
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "test",
		}},
		// invalid: email is not valid
		{false, &NewAccountRequest{
			Username:  "peter5",
			Password:  "Hallo Welt",
			Email:     "test.de",
			FirstName: "roman",
			LastName:  "test",
		}},
	}
	for _, c := range cases {
		err := manager.NewAccount(c.request)
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
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}
	newUserReq := &NewAccountRequest{
		Username:  "felix",
		Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
		Email:     "felix@web.de",
		FirstName: "Felix",
		LastName:  "Heisenberg",
	}
	if err := manager.NewAccount(newUserReq); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// Make sure the users group was created
	groups, err := manager.GetGroupList(&GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of all groups: %v", err)
	}
	if groups[0] != manager.DefaultUserGroup {
		t.Fatalf("expected the default user group %q to have been created", manager.DefaultUserGroup)
	}

	// Make sure that the new account is in the users group
	group, err := manager.GetGroup(manager.DefaultUserGroup, &ListOptions{})
	if err != nil {
		t.Fatalf("failed to get members of the group %q: %v", manager.DefaultUserGroup, err)
	}
	if group.Members[0] != newUserReq.Username {
		t.Fatalf("expected the new user %q to be a member of the default user group %q", newUserReq.Username, manager.DefaultUserGroup)
	}

	isMember, err := manager.IsGroupMember(newUserReq.Username, manager.DefaultUserGroup)
	if err != nil {
		t.Fatalf("failed to check if user %q is in the group %q: %v", newUserReq.Username, manager.DefaultUserGroup, err)
	}
	if !isMember {
		t.Fatalf("expected user %q to be a member of the group %q: %v", newUserReq.Username, manager.DefaultUserGroup, err)
	}

	account, err := manager.GetAccount("felix")
	if err != nil {
		t.Fatalf("failed to get account: %v", err)
	}
	if len(account) != 4 {
		t.Errorf("expected GetAccount to return account with 4 propertiesm, but got: %v", account)
	}
	expected := map[string]string{"givenName": "Felix", "mail": "felix@web.de", "sn": "Heisenberg", "uid": "felix"}
	if diff := cmp.Diff(expected, account); diff != "" {
		t.Errorf("got unexpected account result: (-want +got):\n%s", diff)
	}
}

// TestDeleteAccount ...
func TestDeleteAccount(t *testing.T) {
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}
	users := []string{"user1", "user2"}
	for _, user := range users {
		if err := manager.NewAccount(&NewAccountRequest{
			Username:  user,
			Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
			Email:     "felix@web.de",
			FirstName: "Felix",
			LastName:  "Heisenberg",
		}); err != nil {
			t.Fatalf("failed to add user: %v", err)
		}
	}

	// Assert we find those two users
	userList, err := manager.GetUserList(&GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get users list: %v", err)
	}
	if err := containsUsers(userList, users, manager.AccountAttribute); err != nil {
		t.Error(err)
	}

	// Now delete the first user
	if err := manager.DeleteAccount(users[0]); err != nil {
		t.Fatalf("failed to delete user %q: %v", users[0], err)
	}

	// Assert we find only the second user
	userList, err = manager.GetUserList(&GetUserListRequest{})
	if err != nil {
		t.Fatalf("failed to get users list: %v", err)
	}
	if err := containsUsers(userList, users[1:2], manager.AccountAttribute); err != nil {
		t.Error(err)
	}
}
