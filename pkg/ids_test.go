package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// TestGetHighestGID tests getting the highest GID
func TestGetHighestGID(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// setup adds the admin user and admin and user group
	GID, err := test.Manager.GetHighestGID()
	if err != nil {
		t.Fatalf(
			"failed to get highest GID: %v",
			err,
		)
	}
	if GID != MinGID+2 {
		t.Fatalf(
			"expected GID %d but got %d",
			MinGID+2, GID,
		)
	}

	// add a new group that will claim the highest GID
	strict := false
	groupName := "test-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{test.Manager.DefaultAdminUsername},
	}, strict); err != nil {
		t.Fatalf(
			"failed to create new group: %v",
			err,
		)
	}

	group, err := test.Manager.GetGroupByName(groupName)
	if err != nil {
		t.Fatalf(
			"failed to get group %q: %v",
			groupName, err,
		)
	}
	t.Log(PrettyPrint(group))
	if group.GetGID() != MinGID+2 {
		t.Fatalf(
			"expected GID %d but got %d",
			MinGID+2, group.GetGID(),
		)
	}

	// check that the last GID was updated
	GID, err = test.Manager.GetHighestGID()
	if err != nil {
		t.Fatalf(
			"failed to get highest GID: %v",
			err,
		)
	}
	if GID != MinGID+3 {
		t.Fatalf(
			"expected GID %d but got %d",
			MinGID+3, group.GetGID(),
		)
	}
}

// TestGetHighestUID tests getting the highest UID
func TestGetHighestUID(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// setup adds the admin user and admin and user group
	UID, err := test.Manager.GetHighestUID()
	if err != nil {
		t.Fatalf(
			"failed to get highest UID: %v",
			err,
		)
	}
	if UID != MinUID+1 {
		t.Fatalf(
			"expected UID %d but got %d",
			MinUID+1, UID,
		)
	}

	username := "test-user"
	req := pb.NewUserRequest{
		Username:  username,
		Password:  "some password",
		FirstName: "changeme",
		LastName:  "changeme",
		Email:     "changeme@changeme.com",
	}

	if err := test.Manager.NewUser(&req); err != nil {
		t.Fatalf(
			"failed to create user: %v",
			err,
		)
	}

	// add a new user that will claim the highest UID
	user, err := test.Manager.GetUser(username)
	if err != nil {
		t.Fatalf(
			"failed to get user %q: %v",
			username, err,
		)
	}
	t.Log(PrettyPrint(user))
	if user.GetUID() != MinUID+1 {
		t.Fatalf(
			"expected UID %d but got %d",
			MinUID+1, user.GetUID(),
		)
	}

	// check that the last UID was updated
	UID, err = test.Manager.GetHighestUID()
	if err != nil {
		t.Fatalf(
			"failed to get highest UID: %v",
			err,
		)
	}
	if UID != MinUID+2 {
		t.Fatalf(
			"expected UID %d but got %d",
			MinUID+2, UID,
		)
	}
}
