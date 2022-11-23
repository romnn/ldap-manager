package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetGroup gets a group
func (s *LDAPManagerService) GetGroup(ctx context.Context, req *pb.GetGroupRequest) (*pb.Group, error) {
	claims, err := s.Authenticate(ctx)
	log.Infof("claims: %v", claims)
	if err != nil {
		return nil, err
	}
	groupName := req.GetName()
	group, err := s.manager.GetGroupByName(groupName)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Errorf(codes.Internal, "failed to get group %q", groupName)
	}
	return group, nil
}
