package pkg

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// GetUserGroups gets the groups a user is member of
func (m *LDAPManager) GetUserGroups(req *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	username := EscapeDN(req.GetUsername())
	if !m.GroupMembershipUsesUID {
		username = m.UserNamed(req.GetUsername())
	}
	filter := fmt.Sprintf("(&(objectClass=posixGroup)(%s=%s))", m.GroupMembershipAttribute, username)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"cn"},
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	// groupList := &pb.GroupList{Total: int64(len(result.Entries))}
	var groups []*pb.Group
	for _, entry := range result.Entries {
		if group, err := m.parseGroup(entry); err == nil {
			groups = append(groups, group)
		}
		// if cn := group.GetAttributeValue("cn"); cn != "" {
		// 	groupList.Groups = append(groupList.Groups, cn)
		// }
	}

	// No sorting and clipping here
	return &pb.GroupList{
		Groups: groups,
		Total:  int64(len(groups)),
	}, nil
}
