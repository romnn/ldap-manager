package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetGroupList returns a list of groups
func (s *LDAPManagerService) GetGroupList(ctx context.Context, req *pb.GetGroupListRequest) (*pb.GroupList, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	groups, err := s.manager.GetGroupList(req)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while getting groups")
	}
	return groups, nil
}
