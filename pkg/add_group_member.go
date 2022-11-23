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

// MemberAlreadyExistsError is returned when a user is already a group member
type MemberAlreadyExistsError struct {
	ldaperror.ApplicationError
	Group, Member string
}

func (e *MemberAlreadyExistsError) Error() string {
	return fmt.Sprintf(
		"member %q is already a member of group %q",
		e.Member, e.Group,
	)
}

// StatusError returns the GRPC status error for this error
func (e *MemberAlreadyExistsError) StatusError() error {
	return status.Errorf(codes.AlreadyExists, e.Error())
}

// GroupMemberDN gets the distinguished name of a group member
func (m *LDAPManager) GroupMemberDN(username string) string {
	if !m.GroupMembershipUsesUID {
		return m.UserDN(username)
	}
	return EscapeDN(username)
}

// AddGroupMember adds a user as a group member
func (m *LDAPManager) AddGroupMember(req *pb.GroupMember, allowNonExistent bool) error {
	group := req.GetGroup()
	if group == "" {
		return &ldaperror.ValidationError{
			Message: "group name must not be empty",
		}
	}
	username := req.GetUsername()
	if username == "" {
		return &ldaperror.ValidationError{
			Message: "username must not be empty",
		}
	}
	if !allowNonExistent && !m.IsProtectedGroup(group) {
		memberStatus, err := m.IsGroupMember(&pb.IsGroupMemberRequest{
			Username: username,
			Group:    m.DefaultUserGroup,
		})
		if err != nil {
			return fmt.Errorf(
				"failed to check if member %q exists: %v",
				req.GetUsername(), err,
			)
		}
		if !memberStatus.GetIsMember() {
			return &ZeroOrMultipleUsersError{
				Username: req.GetUsername(),
			}
		}
	}

	memberDN := m.GroupMemberDN(username)
	modifyRequest := ldap.NewModifyRequest(
		m.GroupDN(group),
		[]ldap.Control{},
	)
	modifyRequest.Add(m.GroupMembershipAttribute, []string{memberDN})
	log.Debug(PrettyPrint(modifyRequest))

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Modify(modifyRequest); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultAttributeOrValueExists) {
			return &MemberAlreadyExistsError{
				Member: username,
				Group:  group,
			}
		}
		return err
	}
	log.Infof(
		"added member %q to group %q",
		memberDN, group,
	)
	return nil
}
