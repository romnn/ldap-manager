package pkg

import (
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GroupAlreadyExistsError is returned when a group already exists
type GroupAlreadyExistsError struct {
	ldaperror.ApplicationError
	Group string
}

func (e *GroupAlreadyExistsError) Error() string {
	return fmt.Sprintf(
		"group %q already exists",
		e.Group,
	)
}

// StatusError returns the GRPC status error for this error
func (e *GroupAlreadyExistsError) StatusError() error {
	return status.Errorf(codes.AlreadyExists, e.Error())
}

// NewGroup creates a new group
func (m *LDAPManager) NewGroup(req *pb.NewGroupRequest, strict bool) error {
	groupName := req.GetName()
	if groupName == "" {
		return &ldaperror.ValidationError{
			Message: "group name can not be empty",
		}
	}
	_, err := m.GetGroupByName(groupName)
	if _, notfound := err.(*ZeroOrMultipleGroupsError); !notfound {
		return &GroupAlreadyExistsError{
			Group: groupName,
		}
	}
	// we define the user group to always be MinGID
	GID := MinGID
	if groupName != m.DefaultUserGroup {
		GID, err = m.GetHighestGID()
		if err != nil {
			return err
		}
	}

	var memberList []string
	for _, username := range req.GetMembers() {
		if strict {
			memberStatus, err := m.IsGroupMember(&pb.IsGroupMemberRequest{
				Username: username,
				Group:    m.DefaultUserGroup,
			})
			if err != nil {
				return fmt.Errorf(
					"failed to check if member %q exists: %v",
					username, err,
				)
			}
			if !memberStatus.GetIsMember() {
				log.Warnf(
					"skip adding user %q to group %q (not in user group %q)",
					username, groupName, m.DefaultUserGroup,
				)
				continue
			}
		}
		memberDN := m.GroupMemberDN(username)
		memberList = append(memberList, memberDN)
	}

	var objectClass []string
	if !m.UseRFC2307BISSchema {
		objectClass = []string{
			"top",
			"posixGroup",
		}
	} else {
		if len(memberList) < 1 {
			return &ldaperror.ValidationError{
				Message: "must specify at least one existing group member when using RFC2307BIS (not NIS)",
			}
		}
		objectClass = []string{
			"top",
			"groupOfUniqueNames",
			"posixGroup",
		}
	}

	groupAttributes := []ldap.Attribute{
		{Type: "objectClass", Vals: objectClass},
		{Type: "cn", Vals: []string{EscapeDN(groupName)}},
		{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		{Type: m.GroupMembershipAttribute, Vals: memberList},
	}

	addGroupRequest := &ldap.AddRequest{
		DN:         m.GroupDN(groupName),
		Attributes: groupAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(PrettyPrint(addGroupRequest))

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()
	if err := conn.Add(addGroupRequest); err != nil {
		return err
	}
	if err := m.updateLastID("lastGID", GID+1); err != nil {
		return err
	}
	log.Infof(
		"added new group %q with %d members (gid=%d)",
		groupName, len(memberList), GID,
	)
	return nil
}
