package pkg

import (
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// UpdateGroup ...
func (m *LDAPManager) UpdateGroup(req *pb.UpdateGroupRequest) error {
	// if req.GetName() == "" {
	// 	return &ValidationError{Message: "group name can not be empty"}
	// }

	// groupName := req.GetName()
	// if req.GetNewName() != "" && req.GetNewName() != groupName {
	// 	modifyRequest := &ldap.ModifyDNRequest{
	// 		DN:           m.GroupNamed(groupName),
	// 		NewRDN:       fmt.Sprintf("cn=%s", req.GetNewName()),
	// 		DeleteOldRDN: true,
	// 		NewSuperior:  "",
	// 	}
	// 	log.Debugf("UpdateGroup modifyRequest=%v", modifyRequest)
	// 	if err := m.ldap.ModifyDN(modifyRequest); err != nil {
	// 		return err
	// 	}
	// 	log.Infof("renamed group from %q to %q", req.GetName(), req.GetNewName())
	// 	groupName = req.GetNewName()
	// }

	// modifyGroupRequest := ldap.NewModifyRequest(
	// 	m.GroupNamed(groupName),
	// 	[]ldap.Control{},
	// )
	// if req.GetGid() >= MinGID {
	// 	modifyGroupRequest.Replace("gidNumber", []string{strconv.Itoa(int(req.GetGid()))})
	// }
	// if err := m.ldap.Modify(modifyGroupRequest); err != nil {
	// 	return fmt.Errorf("failed to modify group %q: %v", groupName, err)
	// }
	// log.Infof("updated %d attributes of group %q", len(modifyGroupRequest.Changes), groupName)
	return nil
}
