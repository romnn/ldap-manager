package pkg

import (
	"fmt"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

func (test *Test) isGroupMember(
	t *testing.T,
	username string,
	group string,
	expected bool,
) (*pb.GroupMemberStatus, error) {
	// retry a few times (memberOf can take some time to become available)
	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 10 * time.Second,
	}, 10)

	var memberStatus *pb.GroupMemberStatus
	err := backoff.Retry(func() error {
		var err error
		memberStatus, err = test.Manager.IsGroupMember(
			&pb.IsGroupMemberRequest{
				Username: username,
				Group:    group,
			},
		)
		if err != nil {
			return err
		}
		if memberStatus.IsMember != expected {
			return fmt.Errorf(
				"unexpected member status for user %q of group %q (expected %v but got %v)",
				username, group, expected, memberStatus.IsMember,
			)
		}
		return nil
	}, b)
	return memberStatus, err
}

// TestRemoveGroupMember tests removing a group member
func TestRemoveGroupMember(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	// assert removing from the default user group always fails
	allowRemoveFromDefaultGroups := false
	if err := test.Manager.RemoveGroupMember(&pb.GroupMember{
		Group:    test.Manager.DefaultUserGroup,
		Username: "someone",
	}, allowRemoveFromDefaultGroups); err == nil {
		t.Error("expected removing member of the users group to fail")
	}

	groupName := "test-group"
	usernames := []string{"user1", "user2"}
	for _, username := range usernames {
		if err := test.Manager.NewUser(&pb.NewUserRequest{
			Username:  username,
			Password:  "Hallo Welt",
			Email:     "a@b.de",
			FirstName: "roman",
			LastName:  "d",
		}); err != nil {
			t.Fatalf("failed to add new user: %v", err)
		}
	}
	strict := false
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: usernames,
	}, strict); err != nil {
		t.Fatalf("failed to add new group: %v", err)
	}

	memberStatus, _ := test.isGroupMember(
		t,
		usernames[0], groupName, true,
	)
	if !memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to be a member of group %q",
			usernames[0], groupName,
		)
	}

	// remove first member
	if err := test.Manager.RemoveGroupMember(&pb.GroupMember{
		Group:    groupName,
		Username: usernames[0],
	}, allowRemoveFromDefaultGroups); err != nil {
		t.Fatalf(
			"failed to remove member %q of group %q: %v",
			usernames[0], groupName, err,
		)
	}
	memberStatus, _ = test.isGroupMember(t, usernames[0], groupName, false)
	if memberStatus.GetIsMember() {
		t.Fatalf(
			"expected user %q to be no longer a member of group %q",
			usernames[0], groupName,
		)
	}

	// assert username2 can not be removed
	err := test.Manager.RemoveGroupMember(&pb.GroupMember{
		Username: usernames[1],
		Group:    groupName,
	}, allowRemoveFromDefaultGroups)
	_, lastMember := err.(*RemoveLastGroupMemberError)
	if err == nil || !lastMember {
		t.Fatalf(
			"expected removing last member %q of group %q to fail",
			usernames[1], groupName,
		)
	}
}

// TestRemoveGroupMemberMissing tests removing a group member
// when either the user or the group does not exist.
func TestRemoveGroupMemberMissing(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	strict := false
	groupName := "test-group"
	if err := test.Manager.NewGroup(&pb.NewGroupRequest{
		Name:    groupName,
		Members: []string{"temp-user"},
	}, strict); err != nil {
		t.Fatalf("failed to add new group: %v", err)
	}

	// delete a non-existent member of an existing group
	username := "i-am-not-there"
	allowRemoveFromDefaultGroups := false
	if err := test.Manager.RemoveGroupMember(&pb.GroupMember{
		Group:    groupName,
		Username: username,
	}, allowRemoveFromDefaultGroups); err == nil {
		t.Errorf(
			"expected error removing user %q from group %q",
			username, groupName,
		)
	}

	// add an existing user to an non-existing group
	username = "valid-user"
	if err := test.Manager.NewUser(&pb.NewUserRequest{
		Username:  username,
		Password:  "Hallo Welt",
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}); err != nil {
		t.Fatalf("failed to add new user: %v", err)
	}

	// delete an existent user from an non-existing group
	groupName = "group-that-is-ficticious"
	if err := test.Manager.RemoveGroupMember(&pb.GroupMember{
		Group:    groupName,
		Username: username,
	}, allowRemoveFromDefaultGroups); err == nil {
		t.Errorf(
			"expected error removing user %q from group %q",
			username, groupName,
		)
	}
}
