package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestNewUser tests adding a new user
func TestNewUser(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	username := "romnn"
	password := "hallo welt"
	req := pb.NewUserRequest{
		Username:  username,
		Password:  password,
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}
	expected := &pb.User{
		Username:      "romnn",
		FirstName:     "roman",
		LastName:      "d",
		DisplayName:   "roman d",
		UID:           2001,
		CN:            "roman d",
		GID:           2001,
		LoginShell:    "/bin/bash",
		HomeDirectory: "/home/romnn",
		Email:         "a@b.de",
	}

	// add the user
	if err := test.Manager.NewUser(&req); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// assert the new user is found
	user, err := test.Manager.GetUser(username)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	t.Log(PrettyPrint(user))

	// check if the user data matches
	if equal, diff := EqualProto(expected, user); !equal {
		t.Fatalf("unexpected user: \n%s", diff)
	}

	// check if we can authenticate the user
	user, err = test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		t.Fatalf(
			"cannot authenticate user %q with password %q: %v",
			username, password, err,
		)
	}

	// check if the user data matches
	if equal, diff := EqualProto(expected, user); !equal {
		t.Fatalf("unexpected user: \n%s", diff)
	}
}

// TestNewUserValidation tests validation of user data
func TestNewUserValidation(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	cases := []struct {
		valid   bool
		request *pb.NewUserRequest
	}{
		// invalid: missing everything
		{false, &pb.NewUserRequest{}},
		// invalid: missing username
		{false, &pb.NewUserRequest{
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing password
		{false, &pb.NewUserRequest{
			Username:  "peter1",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing email
		{false, &pb.NewUserRequest{
			Username:  "peter2",
			Password:  "Hallo Welt",
			FirstName: "roman",
			LastName:  "d",
		}},
		// invalid: missing first name
		{false, &pb.NewUserRequest{
			Username: "peter3",
			Password: "Hallo Welt",
			Email:    "a@b.de",
			LastName: "d",
		}},
		// invalid: missing last name
		{false, &pb.NewUserRequest{
			Username:  "peter4",
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
		}},
		// valid: all required fields
		{true, &pb.NewUserRequest{
			Username:  "peter5",
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "test",
		}},
		// invalid: email is not valid
		{false, &pb.NewUserRequest{
			Username:  "peter5",
			Password:  "Hallo Welt",
			Email:     "test.de",
			FirstName: "roman",
			LastName:  "test",
		}},
	}
	for _, c := range cases {
		err := test.Manager.NewUser(c.request)
		if err != nil && c.valid {
			t.Errorf(
				"failed to add valid user: %v",
				err,
			)
		}
		if err == nil && !c.valid {
			t.Errorf(
				"expected error when adding invalid user %v",
				c.request,
			)
		}
	}
}
