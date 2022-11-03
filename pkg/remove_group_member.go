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

// A RemoveLastGroupMemberError is returned when attempting
// to remove the only member of a group
type RemoveLastGroupMemberError struct {
	error
	Group string
}

func (err *RemoveLastGroupMemberError) Error() string {
	return fmt.Sprintf("cannot remove the only remaining group member from group %q, consider deleting the group first", err.Group)
}

func (err *RemoveLastGroupMemberError) StatusError() error {
	return status.Errorf(codes.FailedPrecondition, err.Error())
}

// RemoveGroupMember removes a group member from a group
func (m *LDAPManager) RemoveGroupMember(req *pb.GroupMember, allowRemoveFromDefaultGroups bool) error {
	username := req.GetUsername()
	group := req.GetGroup()
	if group == "" {
		return &ldaperror.ValidationError{Message: "group must not be empty"}
	}
	if username == "" {
		return &ldaperror.ValidationError{Message: "username must not be empty"}
	}
	protected := m.IsProtectedGroup(group)
	if !allowRemoveFromDefaultGroups && protected {
		return &ldaperror.ValidationError{
			Message: "removing members from default user or admin groups not allowed"}
	}
	username = EscapeDN(req.GetUsername())
	if !m.GroupMembershipUsesUID {
		username = m.UserNamed(req.GetUsername())
	}
	modifyRequest := ldap.NewModifyRequest(
		m.GroupNamed(group),
		[]ldap.Control{},
	)
	modifyRequest.Delete(m.GroupMembershipAttribute, []string{username})
	log.Debugf("DeleteGroupMember: modifyRequest=%v", modifyRequest)

	if err := m.ldap.Modify(modifyRequest); err != nil {
		violation := ldap.IsErrorWithCode(err, ldap.LDAPResultObjectClassViolation)
		if violation {
			return &RemoveLastGroupMemberError{
				Group: group,
			}
		}
		notFound := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject)
		noAttribute := ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchAttribute)
		if notFound || noAttribute {
			return &NoSuchMemberError{
				Group:  group,
				Member: username,
			}
		}
		return err
	}
	log.Infof("removed user %q from group %q", username, group)
	return nil
}