package pkg

import (
	// "github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// DeleteAccount ...
func (m *LDAPManager) DeleteUser(req *pb.DeleteUserRequest, keepGroups bool) error {
	// if req.GetUsername() == "" {
	// 	return errors.New("username must not be empty")
	// }
	// if !keepGroups {
	// 	// delete the account from all its groups
	// 	groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{Username: req.GetUsername()})
	// 	if err != nil {
	// 		return fmt.Errorf("failed to get list of groups: %v", err)
	// 	}
	// 	for _, group := range groups.GetGroups() {
	// 		allowDeleteOfDefaultGroups := true
	// 		if err := m.DeleteGroupMember(&pb.GroupMember{Group: group, Username: req.GetUsername()}, allowDeleteOfDefaultGroups); err != nil {
	// 			if _, ok := err.(*RemoveLastGroupMemberError); ok {
	// 				return err
	// 			}
	// 			if _, ok := err.(*NoSuchMemberError); !ok {
	// 				return fmt.Errorf("failed to remove deleted user %q from group %q: %v", req.GetUsername(), group, err)
	// 			}
	// 		}
	// 	}
	// }
	// if err := m.ldap.Del(ldap.NewDelRequest(
	// 	fmt.Sprintf("%s=%s,%s", m.AccountAttribute, escapeDN(req.GetUsername()), m.UserGroupDN),
	// 	[]ldap.Control{},
	// )); err != nil {
	// 	return err
	// }
	// log.Infof("removed account %q", req.GetUsername())
	return nil
}
