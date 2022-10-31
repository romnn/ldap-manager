package pkg

import (
	"testing"
	// "fmt"
	// "strconv"
	// "strings"

	// "github.com/google/go-cmp/cmp"
	// pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// ldaphash "github.com/romnn/ldap-manager/pkg/hash"
)

// TestGetUser ...
func TestGetUser(t *testing.T) {
	// test := new(ldaptest.Test).Setup(t)
	// defer test.Teardown()

	// newUserReq := &pb.NewAccountRequest{
	// 	Account: &pb.Account{
	// 		Username:  "felix",
	// 		Password:  "y&*T R&EGGSAdsnbdjhfb887gfdwe7fFWEFGDSSDEF",
	// 		Email:     "felix@web.de",
	// 		FirstName: "Felix",
	// 		LastName:  "Heisenberg",
	// 	},
	// }
	// if err := test.Manager.NewAccount(newUserReq, pb.HashingAlgorithm_DEFAULT); err != nil {
	// 	t.Fatalf("failed to add user: %v", err)
	// }
	// username := newUserReq.GetAccount().GetUsername()

	// // Make sure the users group was created
	// groups, err := test.Manager.GetGroupList(&pb.GetGroupListRequest{})
	// if err != nil {
	// 	t.Fatalf("failed to get list of all groups: %v", err)
	// }
	// if !contains(groups.GetGroups(), test.Manager.DefaultUserGroup) {
	// 	t.Fatalf("expected the default user group %q to have been created", test.Manager.DefaultUserGroup)
	// }

	// // Make sure that the new account is in the users group
	// group, err := test.Manager.GetGroup(&pb.GetGroupRequest{Name: test.Manager.DefaultUserGroup})
	// if err != nil {
	// 	t.Fatalf("failed to get members of the group %q: %v", test.Manager.DefaultUserGroup, err)
	// }
	// if !contains(group.Members, username) {
	// 	t.Fatalf("expected the new user %q to be a member of the default user group %q", username, test.Manager.DefaultUserGroup)
	// }

	// memberStatus, err := test.Manager.IsGroupMember(&pb.IsGroupMemberRequest{Username: username, Group: test.Manager.DefaultUserGroup})
	// if err != nil {
	// 	t.Fatalf("failed to check if user %q is in the group %q: %v", username, test.Manager.DefaultUserGroup, err)
	// }
	// if !memberStatus.GetIsMember() {
	// 	t.Fatalf("expected user %q to be a member of the group %q: %v", username, test.Manager.DefaultUserGroup, err)
	// }

	// account, err := test.Manager.GetAccount(&pb.GetAccountRequest{Username: "felix"})
	// if err != nil {
	// 	t.Fatalf("failed to get account: %v", err)
	// }
	// expected := map[string]string{
	// 	"cn":            "Felix Heisenberg",
	// 	"displayName":   "Felix Heisenberg",
	// 	"gidNumber":     "2001", // users group should be 2001
	// 	"givenName":     "Felix",
	// 	"homeDirectory": "/home/felix",
	// 	"loginShell":    "/bin/bash",
	// 	"mail":          "felix@web.de",
	// 	"sn":            "Heisenberg",
	// 	"uid":           "felix",
	// 	"uidNumber":     "2002", // admin user should be 2001
	// }
	// if diff := cmp.Diff(expected, account.GetData()); diff != "" {
	// 	t.Errorf("got unexpected account result: (-want +got):\n%s", diff)
	// }
}
