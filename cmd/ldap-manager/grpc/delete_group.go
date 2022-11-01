package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteGroup deletes a group
func (s *LDAPManagerService) DeleteGroup(ctx context.Context, req *pb.DeleteGroupRequest) (*pb.Empty, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.manager.DeleteGroup(req); err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.Error); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while deleting group")
	}
	return &pb.Empty{}, nil
}
