package pkg

import (
	"github.com/go-ldap/ldap/v3"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
)

// DeleteGroup deletes a group
func (m *LDAPManager) DeleteGroup(req *pb.DeleteGroupRequest) error {
	name := req.GetName()
	if name == "" {
		return &ldaperror.ValidationError{
			Message: "group name can not be empty",
		}
	}
	if m.IsProtectedGroup(name) {
		return &ldaperror.ValidationError{
			Message: "deleting the default user or admin group is not allowed",
		}
	}
	if err := m.ldap.Del(ldap.NewDelRequest(
		m.GroupNamed(name),
		[]ldap.Control{},
	)); err != nil {
		if ldap.IsErrorWithCode(err, ldap.LDAPResultNoSuchObject) {
			return &ZeroOrMultipleGroupsError{
				Group: name,
			}
		}
		return err
	}
	log.Infof("removed group %q", name)
	return nil
}
