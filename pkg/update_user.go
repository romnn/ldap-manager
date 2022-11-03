package pkg

import (
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

func buildModifyRequest(dn string, user *pb.User, update *pb.NewUserRequest) *ldap.ModifyRequest {
	req := ldap.NewModifyRequest(
		dn,
		[]ldap.Control{},
	)
	oldFirstName := user.GetFirstName()
	oldLastName := user.GetLastName()

	// upate first name
	if firstName := update.GetFirstName(); firstName != "" {
		oldFirstName = firstName
		req.Replace("givenName", []string{firstName})
	}
	// update last name
	if lastName := update.GetLastName(); lastName != "" {
		oldLastName = lastName
		req.Replace("sn", []string{lastName})
	}
	// update full name
	if update.GetFirstName() != "" || update.GetLastName() != "" {
		fullName := fmt.Sprintf("%s %s", oldFirstName, oldLastName)
		req.Replace("displayName", []string{fullName})
		req.Replace("cn", []string{fullName})
	}
	// update login shell
	if loginShell := update.GetLoginShell(); loginShell != "" {
		req.Replace("loginShell", []string{loginShell})
	}
	// update home directory
	if homeDir := update.GetHomeDirectory(); homeDir != "" {
		req.Replace("homeDirectory", []string{homeDir})
	}
	// update email
	if mail := update.GetEmail(); mail != "" {
		req.Replace("mail", []string{mail})
	}
	return req
}

// MigrateUserGroups migrates a user in each group it is a member of
func (m *LDAPManager) MigrateUserGroups(username, newUsername string) error {
	groups, err := m.GetUserGroups(&pb.GetUserGroupsRequest{
		Username: username,
	})
	if err != nil {
		return fmt.Errorf(
			"failed to get list of groups: %v",
			err,
		)
	}
	for _, group := range groups.GetGroups() {
		// add first, to avoid removing the only member of a group
		allowMissing := false
		if err := m.AddGroupMember(&pb.GroupMember{
			Group:    group.GetName(),
			Username: newUsername,
		}, allowMissing); err != nil {
			_, exists := err.(*MemberAlreadyExistsError)
			if !exists {
				return fmt.Errorf(
					"failed to add renamed user (%q -> %q) to group %q: %v",
					username, newUsername, group.GetName(), err,
				)
			}
		}
		// remove after
		allowRemoveFromDefaultGroups := true
		if err := m.RemoveGroupMember(&pb.GroupMember{
			Group:    group.GetName(),
			Username: username,
		}, allowRemoveFromDefaultGroups); err != nil {
			_, nomember := err.(*NoSuchMemberError)
			if !nomember {
				return fmt.Errorf(
					"failed to remove renamed user (%q -> %q) from group %q: %v",
					username, newUsername, group.GetName(), err,
				)
			}
		}
		log.Infof(
			"Migrated member %q to %q in group %q",
			username, newUsername, group.GetName(),
		)
	}
	return nil
}

// UpdateUser updates a user
func (m *LDAPManager) UpdateUser(req *pb.UpdateUserRequest, isAdmin bool) (string, error) {
	username := req.GetUsername()
	newUsername := req.GetUsername()
	update := req.GetUpdate()

	// check that the user exists
	user, err := m.GetUser(username)
	if err != nil {
		return "", err
	}
	// check if the username should be changed
	if update.GetUsername() != "" && update.GetUsername() != username {
		newUsername = update.GetUsername()
		// make sure the new username is not taken
		_, err := m.GetUser(newUsername)
		if err == nil {
			return "", &UserAlreadyExistsError{
				Username: newUsername,
			}
		}

		modifyRequest := &ldap.ModifyDNRequest{
			DN: m.UserDN(username),
			NewRDN: fmt.Sprintf(
				"%s=%s",
				m.AccountAttribute,
				newUsername,
			),
			DeleteOldRDN: true,
			NewSuperior:  "",
		}
		log.Debug(PrettyPrint(modifyRequest))
		if err := m.ldap.ModifyDN(modifyRequest); err != nil {
			return "", err
		}
		log.Infof(
			"renamed user from %q to %q",
			username, newUsername,
		)

		if err := m.MigrateUserGroups(username, newUsername); err != nil {
			return "", err
		}
	}

	modifyUserReq := buildModifyRequest(
		m.UserDN(newUsername),
		user,
		update,
	)
	if isAdmin {
		// Only the admin is allowed to change these,
		// because they identify a unique user (username + uidNumber)
		if uid := update.GetUID(); uid >= MinUID {
			modifyUserReq.Replace("uidNumber", []string{
				strconv.Itoa(int(uid)),
			})
		}
		if gid := update.GetGID(); gid >= MinGID {
			modifyUserReq.Replace("gidNumber", []string{
				strconv.Itoa(int(gid)),
			})
		}
	}

	log.Debug(PrettyPrint(modifyUserReq))
	if err := m.ldap.Modify(modifyUserReq); err != nil {
		return "", fmt.Errorf(
			"failed to modify existing user: %v",
			err,
		)
	}
	updated := len(modifyUserReq.Changes)
	log.Infof(
		"updated %d attributes of user %q",
		updated, username,
	)

	// change password
	if update.GetPassword() != "" {
		if err := m.ChangePassword(&pb.ChangePasswordRequest{
			Username: newUsername,
			Password: update.GetPassword(),
		}); err != nil {
			return "", fmt.Errorf(
				"failed to change user password: %v",
				err,
			)
		}
	}
	return newUsername, nil
}
