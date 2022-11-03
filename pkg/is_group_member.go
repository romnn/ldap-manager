package pkg

import (
	"fmt"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// IsGroupMember checks if a user is member of a group
func (m *LDAPManager) IsGroupMember(req *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	group, err := m.GetGroupByName(req.GetGroup())
	if err != nil {
		return nil, err
	}
	username := req.GetUsername()
	if !m.GroupMembershipUsesUID {
		username = fmt.Sprintf(
			"%s=%s,%s",
			m.AccountAttribute,
			username,
			m.UserGroupDN,
		)
	}
	for _, member := range group.GetMembers() {
		if member == username {
			return &pb.GroupMemberStatus{
				IsMember: true,
			}, nil
		}
	}
	return &pb.GroupMemberStatus{
		IsMember: false,
	}, nil
}
