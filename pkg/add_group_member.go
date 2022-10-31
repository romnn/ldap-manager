package pkg

import (
	"fmt"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MemberAlreadyExistsError ...
type MemberAlreadyExistsError struct {
	error
	Group, Member string
}

// Error ...
func (e *MemberAlreadyExistsError) Error() string {
	return fmt.Sprintf("member %q is already a member of group %q", e.Member, e.Group)
}

// StatusError...
func (e *MemberAlreadyExistsError) StatusError() error {
	return status.Errorf(codes.AlreadyExists, e.Error())
}

// AddGroupMember ...
func (m *LDAPManager) AddGroupMember(req *pb.GroupMember, allowNonExistent bool) error {
	group := req.GetGroup()
	if group == "" {
		return &ldaperror.ValidationError{Message: "group name must not be empty"}
	}
	username := req.GetUsername()
	if username == "" {
		return &ldaperror.ValidationError{Message: "username must not be empty"}
	}
	if !allowNonExistent && !m.IsProtectedGroup(group) {
		memberStatus, err := m.IsGroupMember(&pb.IsGroupMemberRequest{
			Username: username,
			Group:    m.DefaultUserGroup,
		})
		if err != nil {
			return fmt.Errorf("failed to check if member %q exists: %v", req.GetUsername(), err)
		}
		if !memberStatus.GetIsMember() {
			return &ZeroOrMultipleUsersError{
				Username: req.GetUsername(),
			}
		}
	}

	member := EscapeDN(username)
	if !m.GroupMembershipUsesUID {
		member = m.UserNamed(username)
	}
	modifyRequest := ldap.NewModifyRequest(
		m.GroupNamed(group),
		[]ldap.Control{},
	)
	modifyRequest.Add(m.GroupMembershipAttribute, []string{member})
	log.Debugf("AddGroupMember: modifyRequest=%v", modifyRequest)
	if err := m.ldap.Modify(modifyRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultAttributeOrValueExists) {
			return &MemberAlreadyExistsError{Member: username, Group: group}
		}
		return err
	}
	log.Infof("added member %q to group %q", member, group)
	return nil
}
