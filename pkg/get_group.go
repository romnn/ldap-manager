package pkg

import (
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// A ZeroOrMultipleGroupsError is returned when zero or multiple
// groups are found
type ZeroOrMultipleGroupsError struct {
	Group string
	GID   int
	Count int
}

func (e *ZeroOrMultipleGroupsError) groupName() string {
	if e.Group != "" {
		return fmt.Sprintf("name %q", e.Group)
	}
	return fmt.Sprintf("GID %d", e.GID)
}

func (e *ZeroOrMultipleGroupsError) Error() string {
	if e.Count > 1 {
		return fmt.Sprintf(
			"multiple (%d) groups with %s",
			e.Count, e.groupName(),
		)
	}
	return fmt.Sprintf(
		"no group with %s",
		e.groupName(),
	)
}

// StatusError returns the GRPC status error for this error
func (e *ZeroOrMultipleGroupsError) StatusError() error {
	if e.Count > 1 {
		return status.Errorf(codes.Internal, e.Error())
	}
	return status.Errorf(codes.NotFound, e.Error())
}

const (
	groupGidNumber = "gidNumber"
	groupCN        = "cn"
)

// ParseGroup parses an ldap.Entry as a group
func (m *LDAPManager) parseGroup(entry *ldap.Entry) (*pb.Group, error) {
	var members []string
	for _, member := range entry.GetAttributeValues(m.GroupMembershipAttribute) {
		members = append(members, member)
	}
	GID, err := strconv.Atoi(entry.GetAttributeValue(groupGidNumber))
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse group GID as integer: %v",
			err,
		)
	}
	return &pb.Group{
		Members: members,
		Name:    entry.GetAttributeValue(groupCN),
		GID:     int64(GID),
	}, nil
}

func (m *LDAPManager) groupFields() []string {
	return []string{
		m.GroupMembershipAttribute,
		groupGidNumber,
		groupCN,
	}
}

// GetGroupByGID gets a group by its GID
func (m *LDAPManager) GetGroupByGID(GID int) (*pb.Group, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(gid=%d)", GID),
		m.groupFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleGroupsError{
			GID:   GID,
			Count: len(result.Entries),
		}
	}
	return m.parseGroup(result.Entries[0])
}

// GetGroupByName gets a group by its name
func (m *LDAPManager) GetGroupByName(name string) (*pb.Group, error) {
	result, err := m.ldap.Search(ldap.NewSearchRequest(
		m.GroupsDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(cn=%s)", EscapeFilter(name)),
		m.groupFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	if len(result.Entries) != 1 {
		return nil, &ZeroOrMultipleGroupsError{
			Group: name,
			Count: len(result.Entries),
		}
	}
	return m.parseGroup(result.Entries[0])
}
