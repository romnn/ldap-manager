package pkg

import (
	"fmt"
	"sort"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// CountUsers counts the number of total users
func (m *LDAPManager) CountUsers() (int, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.UserGroupDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(%s=*)", m.AccountAttribute),
		[]string{"dn"},
		[]ldap.Control{},
	))
	if err != nil {
		return 0, err
	}
	return len(result.Entries), nil
}

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
		m.defaultUserFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	total, err := m.CountUsers()
	if err != nil {
		return nil, err
	}
	users := make(map[string]*pb.UserData)
	for _, entry := range result.Entries {
		if entryKey := entry.GetAttributeValue(req.GetSortKey()); entryKey != "" {
			users[entryKey] = ParseUser(entry)
		}
	}
	// Sort for deterministic clipping
	keys := make([]string, 0, len(users))
	for k := range users {
		keys = append(keys, k)
	}
	// Sort
	sort.Slice(keys, func(i, j int) bool {
		asc := keys[i] < keys[j]
		if req.GetSortOrder() == pb.SortOrder_ASCENDING {
			return !asc
		}
		return asc
	})
	// Clip
	clippedKeys := keys
	if req.GetStart() >= 0 && req.GetEnd() < int32(len(keys)) && req.GetStart() < req.GetEnd() {
		clippedKeys = keys[req.GetStart():req.GetEnd()]
	}
	// total := 0
	clipped := &pb.UserList{Total: int64(total)}
	for _, key := range clippedKeys {
		clipped.Users = append(clipped.Users, users[key])
	}
	return clipped, nil
}
