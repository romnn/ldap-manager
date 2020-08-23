package ldapmanager

import (
	"testing"

	ldaptest "github.com/romnnn/ldap-manager/test"
)

// TestAddNewUser ...
func TestAddNewUser(t *testing.T) {
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()
	manager := NewLDAPManager(test.OpenLDAPCConfig)
	if err := manager.Setup(); err != nil {
		t.Fatal(err)
	}

	// Add user
	newUserReq := &NewAccountRequest{
		Username:  "romnn",
		Password:  "Hallo Welt",
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}
	if err := manager.NewAccount(newUserReq); err != nil {
		t.Errorf("failed to add user: %v", err)
	}

	// List all users
	users, err := manager.GetUserList(&GetUserListRequest{})
	if err != nil {
		t.Errorf("failed to add user: %v", err)
	}
	found := false
	for _, user := range users {
		if uid, ok := user[manager.AccountAttribute]; ok && uid == newUserReq.Username {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected to find user %q after it was added but only got %v", newUserReq.Username, users)
	}
}
