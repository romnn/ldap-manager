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

// A NoSuchMemberError is returned when the group does not contain the member
type NoSuchMemberError struct {
	error
	Group, Member string
}

// Error ...
func (err *NoSuchMemberError) Error() string {
	return fmt.Sprintf("no member %q in group %q", err.Member, err.Group)
}

// StatusError ...
func (err *NoSuchMemberError) StatusError() error {
	return status.Errorf(codes.NotFound, err.Error())
}

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

// DeleteUser deletes a user
func (m *LDAPManager) DeleteUser(req *pb.DeleteUserRequest, keepGroups bool) error {
	username := req.GetUsername()
	if username == "" {
		return &ldaperror.ValidationError{Message: "username must not be empty"}
	}
	if !keepGroups {
		// delete the account from all its groups
		groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{
			Username: username,
		})
		if err != nil {
			return fmt.Errorf("failed to get list of groups: %v", err)
		}
		for _, group := range groups.GetGroups() {
			allowDeleteOfDefaultGroups := true
			if err := m.RemoveGroupMember(&pb.GroupMember{
				Group:    group.GetName(),
				Username: username,
			}, allowDeleteOfDefaultGroups); err != nil {
				if _, ok := err.(*RemoveLastGroupMemberError); ok {
					return err
				}
				if _, ok := err.(*NoSuchMemberError); !ok {
					return err
					// fmt.Errorf("failed to remove user %q from group %q: %v", username, group, err)
				}
			}
		}
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		fmt.Sprintf("%s=%s,%s", m.AccountAttribute, EscapeDN(username), m.UserGroupDN),
		[]ldap.Control{},
	)); err != nil {
		return err
	}
	log.Infof("removed account %q", username)
	return nil
}
