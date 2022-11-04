package pkg

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// GetUserList gets a list of all users
func (m *LDAPManager) GetUserList(req *pb.GetUserListRequest) (*pb.UserList, error) {
	if req.GetSortKey() == "" {
		req.SortKey = m.AccountAttribute
	}
	conn, err := m.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	result, err := conn.SearchWithPaging(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf(
			"(&(%s=*)%s)",
			m.AccountAttribute,
			BuildFilter(req.GetFilter()),
		),
		m.userFields(),
		[]ldap.Control{},
	), pagingSize)
	if err != nil {
		return nil, err
	}
	users := make(map[string]*pb.User)
	for _, entry := range result.Entries {
		sortKey := req.GetSortKey()
		if entryKey := entry.GetAttributeValue(sortKey); entryKey != "" {
			users[entryKey] = m.ParseUser(entry)
		}
	}

	// sort for deterministic clipping
	keys := make([]string, 0, len(users))
	for k := range users {
		keys = append(keys, k)
	}

	// sort
	sort.Slice(keys, func(i, j int) bool {
		asc := keys[i] < keys[j]
		if req.GetSortOrder() == pb.SortOrder_ASCENDING {
			return !asc
		}
		return asc
	})

	// clip
	clippedKeys := keys
	validStart := req.GetStart() >= 0
	validEnd := req.GetEnd() < int32(len(keys))
	validRange := req.GetStart() < req.GetEnd()
	if validStart && validEnd && validRange {
		clippedKeys = keys[req.GetStart():req.GetEnd()]
	}
	var clipped []*pb.User
	for _, key := range clippedKeys {
		clipped = append(clipped, users[key])
	}
	return &pb.UserList{
		Users: clipped,
		Total: int64(len(users)),
	}, nil
}
