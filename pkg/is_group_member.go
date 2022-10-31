package pkg

import (
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// IsGroupMember ...
func (m *LDAPManager) IsGroupMember(req *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	var status pb.GroupMemberStatus
	// result, err := m.findGroup(req.Group, []string{"dn", m.GroupMembershipAttribute})
	// if err != nil {
	// 	return &status, err
	// }
	// if len(result.Entries) != 1 {
	// 	return &status, &ZeroOrMultipleGroupsError{Group: req.GetGroup(), Count: len(result.Entries)}
	// }
	// if !m.GroupMembershipUsesUID {
	// 	req.Username = fmt.Sprintf("%s=%s,%s", m.AccountAttribute, req.GetUsername(), m.UserGroupDN)
	// }
	// for _, member := range result.Entries[0].GetAttributeValues(m.GroupMembershipAttribute) {
	// 	if member == req.GetUsername() {
	// 		return &pb.GroupMemberStatus{IsMember: true}, nil
	// 	}
	// }
	return &status, nil
}
