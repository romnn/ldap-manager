package pkg

import (
	"fmt"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// UpdateGroup updates a group
func (m *LDAPManager) UpdateGroup(req *pb.UpdateGroupRequest) error {
	groupName := req.GetName()
	if groupName == "" {
		return &ldaperror.ValidationError{
			Message: "group name must not be empty",
		}
	}

	conn, err := m.Pool.Get()
	if err != nil {
		return err
	}
	defer conn.Close()

	newGroupName := req.GetNewName()
	if newGroupName != "" && newGroupName != groupName {
		modifyRequest := &ldap.ModifyDNRequest{
			DN:           m.GroupDN(groupName),
			NewRDN:       fmt.Sprintf("cn=%s", newGroupName),
			DeleteOldRDN: true,
			NewSuperior:  "",
		}
		log.Debug(PrettyPrint(modifyRequest))

		if err := conn.ModifyDN(modifyRequest); err != nil {
			return fmt.Errorf(
				"failed to rename group %q to %q",
				groupName, newGroupName,
			)
		}
		log.Infof(
			"renamed group from %q to %q",
			groupName, newGroupName,
		)
		groupName = newGroupName
	}

	modifyGroupRequest := ldap.NewModifyRequest(
		m.GroupDN(groupName),
		[]ldap.Control{},
	)
	// update GID
	if req.GetGID() >= MinGID {
		GID := strconv.Itoa(int(req.GetGID()))
		modifyGroupRequest.Replace("gidNumber", []string{GID})
	}
	if err := conn.Modify(modifyGroupRequest); err != nil {
		return fmt.Errorf(
			"failed to modify group %q: %v",
			groupName, err,
		)
	}
	updated := len(modifyGroupRequest.Changes)
	log.Infof(
		"updated %d attributes of group %q",
		updated, groupName,
	)
	return nil
}
