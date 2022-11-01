package pkg

import (
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	"testing"
)

// TestAuthenticateUser tests authenticating as a user
func TestAuthenticateUser(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	// create new user
	username := "Test User"
	password := "secret password"
	req := &pb.NewUserRequest{
		Username:  username,
		Password:  password,
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}
	if err := test.Manager.NewUser(req); err != nil {
		t.Fatalf("failed to add user %q: %v", username, err)
	}

	// check that authenticating the user using wrong password fails
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: "wrong",
	}); err == nil {
		t.Fatalf("authenticating user %q with wrong password succeeded", username)
	}

	// check if we can authenticate the user using correct password
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: password,
	}); err != nil {
		t.Fatalf("cannot authenticate user %q with password %q: %v", username, password, err)
	}
}
