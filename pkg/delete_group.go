package pkg

import (
	// "fmt"
	// "sort"
	// "strconv"
	// "strings"

	// "google.golang.org/grpc/codes"

	// "github.com/go-ldap/ldap/v3"
	// "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
)

// DeleteGroup ...
func (m *LDAPManager) DeleteGroup(req *pb.DeleteGroupRequest) error {
	// if req.GetName() == "" {
	// 	return &ValidationError{Message: "group name can not be empty"}
	// }
	// if m.IsProtectedGroup(req.GetName()) {
	// 	return &ValidationError{Message: "deleting the default user or admin group is not allowed"}
	// }
	// if err := m.ldap.Del(ldap.NewDelRequest(
	// 	m.GroupNamed(req.GetName()),
	// 	[]ldap.Control{},
	// )); err != nil {
	// 	if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
	// 		return &ZeroOrMultipleGroupsError{Group: req.GetName()}
	// 	}
	// 	return err
	// }
	// log.Infof("removed group %q", req.GetName())
	return nil
}
