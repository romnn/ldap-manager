package pkg

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// GetUserGroups gets the groups a user is member of
func (m *LDAPManager) GetUserGroups(req *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	username := EscapeDN(req.GetUsername())
	if !m.GroupMembershipUsesUID {
		username = m.UserNamed(req.GetUsername())
	}
	filter := fmt.Sprintf(
		"(&(objectClass=posixGroup)(%s=%s))",
		m.GroupMembershipAttribute, username,
	)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		m.groupFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	var groups []*pb.Group
	for _, entry := range result.Entries {
		group, err := m.parseGroup(entry)
		if err != nil {
			log.Warnf(
				"failed to parse group %s: %v",
				PrettyPrint(entry), err,
			)
		} else {
			groups = append(groups, group)
		}
	}

	// No sorting and clipping here
	return &pb.GroupList{
		Groups: groups,
		Total:  int64(len(groups)),
	}, nil
}
