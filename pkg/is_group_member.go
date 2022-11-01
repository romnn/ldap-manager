package pkg

import (
	"fmt"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// IsGroupMember checks if a user is member of a group
func (m *LDAPManager) IsGroupMember(req *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	// result, err := m.findGroup(req.Group, []string{"dn", m.GroupMembershipAttribute})
	group, err := m.GetGroupByName(req.GetGroup())
	// , []string{"dn", m.GroupMembershipAttribute})
	if err != nil {
		return nil, err
	}
	// if len(result.Entries) != 1 {
	// 	return nil, &ZeroOrMultipleGroupsError{
	// 		Group: req.GetGroup(),
	// 		Count: len(result.Entries),
	// 	}
	// }
	username := req.GetUsername()
	if !m.GroupMembershipUsesUID {
		username = fmt.Sprintf("%s=%s,%s", m.AccountAttribute, username, m.UserGroupDN)
	}
	// group := groups.Entries[0]
	// for _, member := range group.GetAttributeValues(m.GroupMembershipAttribute) {
	for _, member := range group.GetMembers() {
		if member == username {
			return &pb.GroupMemberStatus{IsMember: true}, nil
		}
	}
	return &pb.GroupMemberStatus{IsMember: false}, nil
}
