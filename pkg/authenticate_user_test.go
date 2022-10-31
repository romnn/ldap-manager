package pkg

import (
	"testing"
	// "fmt"
	// "strconv"
	// "strings"
	// "github.com/google/go-cmp/cmp"
	// ldaphash "github.com/romnn/ldap-manager/pkg/hash"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestAuthenticateUser tests authenticating as a user
func TestAuthenticateUser(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	// create new user
	username := "Test User"
	password := "secret password"
	req := &pb.NewUserRequest{
		Account: &pb.Account{
			// Username:  name + pw + strconv.Itoa(attemptsLeft),
			Username:  username,
			Password:  password,
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		},
	}
	// username := newUserReq.GetAccount().GetUsername()
	if err := test.Manager.NewUser(req, pb.HashingAlgorithm_DEFAULT); err != nil {
		t.Fatalf("failed to add user %q: %v", username, err)
		// if attemptsLeft <= 0 {
		// finalErr = fmt.Errorf("failed to add user %q: %v", newUserReq.GetAccount().GetUsername(), err)
		// break
		// }
		// continue
	}

	// check that authenticating the user using wrong password fails
	if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: "wrong",
	}); err == nil {
		t.Fatalf("authenticating user %q with wrong password succeeded", username)
	}

	// check if we can authenticate the user using correct password
	if user, err := test.Manager.AuthenticateUser(&pb.LoginRequest{
		Username: username,
		Password: password,
	}); err != nil {
		t.Fatalf("failed to authenticate user %q with password %q: %v", username, password, err)
		t.Logf("user: %v", user)

		// if attemptsLeft <= 0 {
		// finalErr = fmt.Errorf("failed to authenticate user %q with password %q: %v", newUserReq.GetAccount().GetUsername(), pw, err)
		// break
		// }
		// continue
	}

	// samplePasswords := []string{"123456", "Hallo@Welt", "@#73sAdf0^E^RC#+++83230*###$&"}
	// for _, algorithm := range ldaphash.LDAPPasswordHashingAlgorithms {
	// 	name, _ := pb.HashingAlgorithm_name[int32(algorithm)]
	// 	for _, pw := range samplePasswords {
	// 		// t.Log(name, algorithm, pw)
	// 		var finalErr error
	// 		attemptsLeft := 5
	// 		for {
	// 			// FIXME: this tests is flaky :(
	// 			attemptsLeft--
	// 			newUserReq := &pb.NewAccountRequest{
	// 				Account: &pb.Account{
	// 					Username:  name + pw + strconv.Itoa(attemptsLeft),
	// 					Password:  pw,
	// 					Email:     "a@b.de",
	// 					FirstName: "roman",
	// 					LastName:  "d",
	// 				},
	// 			}
	// 			if err := test.Manager.NewAccount(newUserReq, algorithm); err != nil {
	// 				if attemptsLeft <= 0 {
	// 					finalErr = fmt.Errorf("failed to add user %q: %v", newUserReq.GetAccount().GetUsername(), err)
	// 					break
	// 				}
	// 				continue
	// 			}

	// 			// now check if we can authenticate using the clear password
	// 			if _, err := test.Manager.AuthenticateUser(&pb.LoginRequest{Username: newUserReq.GetAccount().GetUsername(), Password: pw}); err != nil {
	// 				if attemptsLeft <= 0 {
	// 					finalErr = fmt.Errorf("failed to authenticate user %q with password %q: %v", newUserReq.GetAccount().GetUsername(), pw, err)
	// 					break
	// 				}
	// 				continue
	// 			}
	// 			break
	// 		}
	// 		if finalErr != nil {
	// 			t.Error(finalErr)
	// 		}
	// 	}
	// }
}
