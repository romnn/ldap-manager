package main

import (
	"testing"

	ldaptest "github.com/romnnn/ldap-manager/test"
)

// TestAddNewUser ...
func TestAddNewUser(t *testing.T) {
	// t.Skip()
	test := new(ldaptest.Test).Setup(t)
	defer test.Teardown()

	database, err := Connect(&test.MongoConfig)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(database)
	t.Fatalf("not implemented")
}
