package pkg

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// IsGroupMember checks if a user is member of a group
func (m *LDAPManager) IsGroupMember(req *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	if _, err := m.GetGroupByName(req.GetGroup()); err != nil {
		return nil, err
	}
	conn, err := m.Pool.Get()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	filter := fmt.Sprintf(
		"(&(objectClass=posixAccount)(%s=%s)(memberOf=%s))",
		m.AccountAttribute,
		EscapeFilter(req.GetUsername()),
		m.GroupDN(req.GetGroup()),
	)
	log.Info(filter)
	log.Info(m.BaseDN)
	result, err := conn.Search(ldap.NewSearchRequest(
		m.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		m.userFields(),
		[]ldap.Control{},
	))
	if err != nil {
		return nil, err
	}
	log.Info(result.Entries)
	return &pb.GroupMemberStatus{
		IsMember: len(result.Entries) > 0,
	}, nil
}
