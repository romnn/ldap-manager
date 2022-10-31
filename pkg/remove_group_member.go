package pkg

import (
	// "fmt"
	// "sort"
	// "strconv"

	// "github.com/go-ldap/ldap/v3"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
)

// RemoveGroupMember ...
func (m *LDAPManager) RemoveGroupMember(req *pb.GroupMember, allowDeleteDefaultGroups bool) error {
	// if req.GetGroup() == "" || req.GetUsername() == "" {
	// 	return &ValidationError{Message: "group and user name can not be empty"}
	// }
	// if !allowDeleteOfDefaultGroups && m.IsProtectedGroup(req.GetGroup()) {
	// 	return &ValidationError{Message: "deleting members from the default user or admin group is not allowed"}
	// }
	// username := escapeDN(req.GetUsername())
	// if !m.GroupMembershipUsesUID {
	// 	username = m.AccountNamed(req.GetUsername())
	// }
	// modifyRequest := ldap.NewModifyRequest(
	// 	m.GroupNamed(req.GetGroup()),
	// 	[]ldap.Control{},
	// )
	// modifyRequest.Delete(m.GroupMembershipAttribute, []string{username})
	// log.Debugf("DeleteGroupMember: modifyRequest=%v", modifyRequest)
	// if err := m.ldap.Modify(modifyRequest); err != nil {
	// 	if ldap.IsErrorWithCode(err, ldap.LDAPResultObjectClassViolation) {
	// 		return &RemoveLastGroupMemberError{Group: req.GetGroup()}
	// 	}
	// 	if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) || ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchAttribute) {
	// 		return &NoSuchMemberError{Group: req.GetGroup(), Member: req.GetUsername()}
	// 	}
	// 	return err
	// }
	// log.Infof("removed user %q from group %q", username, req.GetGroup())
	return nil
}
