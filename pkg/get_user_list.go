package pkg

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// // CountUsers counts the number of total users
// func (m *LDAPManager) CountUsers() (int, error) {
// 	result, err := m.ldap.Search(ldap.NewSearchRequest(
// 		m.UserGroupDN,
// 		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
// 		fmt.Sprintf("(%s=*)", m.AccountAttribute),
// 		[]string{"dn"},
// 		[]ldap.Control{},
// 	))
// 	if err != nil {
// 		return 0, err
// 	}
// 	return len(result.Entries), nil
// }

// GetUserList gets a list of all users
func (m *LDAPManager) GetUserList(req *pb.GetUserListRequest) (*pb.UserList, error) {
	if req.GetSortKey() == "" {
		req.SortKey = m.AccountAttribute
	}
	filter := ParseFilter(req.Filter)
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(%s=*)%s)", m.AccountAttribute, filter),
		m.userFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	// total, err := m.CountUsers()
	// if err != nil {
	// 	return nil, err
	// }
	users := make(map[string]*pb.User)
	for _, entry := range result.Entries {
		if entryKey := entry.GetAttributeValue(req.GetSortKey()); entryKey != "" {
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
