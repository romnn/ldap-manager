package pkg

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GroupAlreadyExistsError ...
type GroupAlreadyExistsError struct {
	error
	Group string
}

func (e *GroupAlreadyExistsError) Error() string {
	return fmt.Sprintf(
		"group %q already exists",
		e.Group,
	)
}

func (e *GroupAlreadyExistsError) StatusError() error {
	return status.Errorf(codes.AlreadyExists, e.Error())
}

// IsProtectedGroup ...
func (m *LDAPManager) IsProtectedGroup(group string) bool {
	isAdminGroup := strings.ToLower(group) == strings.ToLower(m.DefaultAdminGroup)
	isUserGroup := strings.ToLower(group) == strings.ToLower(m.DefaultUserGroup)
	return isAdminGroup || isUserGroup
}

// GroupNamed ...
func (m *LDAPManager) GroupNamed(name string) string {
	return fmt.Sprintf(
		"cn=%s,%s",
		EscapeDN(name), m.GroupsDN,
	)
}

// UserNamed ...
func (m *LDAPManager) UserNamed(name string) string {
	return fmt.Sprintf(
		"%s=%s,%s",
		m.AccountAttribute,
		EscapeDN(name),
		m.UserGroupDN,
	)
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
	GID, err := m.GetHighestGID()
	if err != nil {
		return err
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
		member := EscapeDN(username)
		if !m.GroupMembershipUsesUID {
			member = m.UserNamed(username)
		}
		memberList = append(memberList, member)
	}

	var groupAttributes []ldap.Attribute
	if !m.UseRFC2307BISSchema {
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{
				"top",
				"posixGroup",
			}},
			{Type: "cn", Vals: []string{EscapeDN(groupName)}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		}
	} else {
		if len(memberList) < 1 {
			return &ldaperror.ValidationError{
				Message: "must specify at least one existing group member when using RFC2307BIS (not NIS)",
			}
		}
		groupAttributes = []ldap.Attribute{
			{Type: "objectClass", Vals: []string{
				"top",
				"groupOfUniqueNames",
				"posixGroup",
			}},
			{Type: "cn", Vals: []string{EscapeDN(groupName)}},
			{Type: "gidNumber", Vals: []string{strconv.Itoa(GID)}},
		}
	}

	groupAttributes = append(groupAttributes, ldap.Attribute{
		Type: m.GroupMembershipAttribute,
		Vals: memberList,
	})

	addGroupRequest := &ldap.AddRequest{
		DN:         m.GroupNamed(groupName),
		Attributes: groupAttributes,
		Controls:   []ldap.Control{},
	}
	log.Debug(PrettyPrint(addGroupRequest))
	if err := m.ldap.Add(addGroupRequest); err != nil {
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
