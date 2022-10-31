package pkg

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

const (
	PagingSize = 256
)

// CountGroups counts the number of groups
// func (m *LDAPManager) CountGroups() (int, error) {
// 	result, err := m.ldap.Search(ldap.NewSearchRequest(
// 		m.GroupsDN,
// 		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
// 		"(objectClass=posixGroup)",
// 		[]string{"cn"},
// 		[]ldap.Control{},
// 	))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return len(result.Entries), nil
// }

// GetGroupList gets a list of all groups
func (m *LDAPManager) GetGroupList(req *pb.GetGroupListRequest) (*pb.GroupList, error) {
	filter := ParseFilter(req.Filter)
	result, err := m.ldap.SearchWithPaging(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=posixGroup)%s)", filter),
		[]string{},
		[]ldap.Control{},
	), PagingSize)
	if err != nil {
		return nil, err
	}
	groups := []*pb.Group{}
	for _, entry := range result.Entries {
		group, err := m.ParseGroup(entry)
		if err != nil {
			log.Errorf("failed to parse group %v: %v", entry, err)
			continue
		}
		groups = append(groups, group)
	}

	// Sort
	sort.Slice(groups, func(i, j int) bool {
		asc := groups[i].GetName() < groups[j].GetName()
		if req.GetSortOrder() == pb.SortOrder_DESCENDING {
			return !asc
		}
		return asc
	})

	// Clip
	validStart := req.GetStart() >= 0
	validEnd := req.GetEnd() < int32(len(groups))
	validRange := req.GetStart() < req.GetEnd()
	if validRange && validStart && validEnd {
		groups = groups[req.GetStart():req.GetEnd()]
	}
	return &pb.GroupList{
		Groups: groups,
		Total:  int64(len(groups)),
	}, nil
}
