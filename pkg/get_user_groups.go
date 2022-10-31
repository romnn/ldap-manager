package pkg

import (
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// GetUserGroups ...
func (m *LDAPManager) GetUserGroups(req *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	// username := escapeDN(req.GetUsername())
	// if !m.GroupMembershipUsesUID {
	// 	username = m.AccountNamed(req.GetUsername())
	// }
	// filter := fmt.Sprintf("(&(objectClass=posixGroup)(%s=%s))", m.GroupMembershipAttribute, username)
	// result, err := m.ldap.Search(ldap.NewSearchRequest(
	// 	m.GroupsDN,
	// 	ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	// 	filter,
	// 	[]string{"cn"},
	// 	[]ldap.Control{},
	// ))
	// if err != nil {
	// 	return nil, err
	// }
	// groupList := &pb.GroupList{Total: int64(len(result.Entries))}
	// for _, group := range result.Entries {
	// 	if cn := group.GetAttributeValue("cn"); cn != "" {
	// 		groupList.Groups = append(groupList.Groups, cn)
	// 	}
	// }
	// // No sorting and clipping here
	// return groupList, nil
	return &pb.GroupList{}, nil
}
