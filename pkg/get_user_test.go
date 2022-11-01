package pkg

import (
	"testing"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

func getGroupListGroupNames(groups *pb.GroupList) []string {
	var names []string
	for _, group := range groups.GetGroups() {
		names = append(names, group.GetName())
	}
	return names
}

// TestGetUser tests getting a user
func TestGetUser(t *testing.T) {
	test := new(Test).Setup(t)
	defer test.Teardown()

	username := "felix"
	req := pb.NewUserRequest{
		Username:  username,
		Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
		Email:     "felix@web.de",
		FirstName: "Felix",
		LastName:  "Heisenberg",
	}
	if err := test.Manager.NewUser(&req); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// assert the users group was created
	groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	if err != nil {
		t.Fatalf("failed to get list of all groups: %v", err)
	}
	t.Log(PrettyPrint(groups))
	userGroupName := test.Manager.DefaultUserGroup
	groupNames := getGroupListGroupNames(groups)
	if !Contains(groupNames, userGroupName) {
		t.Fatalf("expected the default user group %q to exist", userGroupName)
	}

	// assert that the new user is in the users group
	group, err := test.Manager.GetGroupByName(userGroupName)
	if err != nil {
		t.Fatalf("failed to get members of group %q: %v", userGroupName, err)
	}
	t.Log(PrettyPrint(group))
	if !Contains(group.Members, test.Manager.UserNamed(username)) {
		t.Fatalf("expected new user %q to be a member of the default user group %q", username, userGroupName)
	}

	// assert that the new user is member of the user group
	memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: username,
		Group:    userGroupName,
	})
	if err != nil {
		t.Fatalf("failed to check if user %q is member of group %q: %v", username, userGroupName, err)
	}
	t.Log(PrettyPrint(memberStatus))
	if !memberStatus.GetIsMember() {
		t.Fatalf("expected user %q to be member of group %q", username, userGroupName)
	}

	// assert the user data matches
	user, err := test.Manager.GetUser(username)
	if err != nil {
		t.Fatalf("failed to get user: %v", err)
	}
	expected := &pb.User{
		Username:      username,
		FirstName:     "Felix",
		LastName:      "Heisenberg",
		DisplayName:   "Felix Heisenberg",
		UID:           2001,
		CN:            "Felix Heisenberg",
		GID:           2001,
		LoginShell:    "/bin/bash",
		HomeDirectory: "/home/felix",
		Email:         "felix@web.de",
	}
	if equal, diff := EqualProto(expected, user); !equal {
		t.Fatalf("unexpected user: \n%s", diff)
	}
}
