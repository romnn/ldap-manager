package pkg

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

const (
	pagingSize = 256
)

// GetGroupList gets a list of all groups
func (m *LDAPManager) GetGroupList(req *pb.GetGroupListRequest) (*pb.GroupList, error) {
	conn, err := m.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	result, err := conn.SearchWithPaging(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(
			"(&(objectClass=posixGroup)%s)",
			BuildFilter(req.GetFilter()),
		),
		[]string{},
		[]ldap.Control{},
	), pagingSize)
	if err != nil {
		return nil, err
	}
	groups := []*pb.Group{}
	for _, entry := range result.Entries {
		group, err := m.parseGroup(entry)
		if err != nil {
			log.Errorf(
				"failed to parse group %v: %v",
				entry, err,
			)
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
