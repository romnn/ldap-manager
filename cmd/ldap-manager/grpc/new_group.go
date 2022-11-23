package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewGroup adds a new LDAP group
func (s *LDAPManagerService) NewGroup(ctx context.Context, req *pb.NewGroupRequest) (*pb.Empty, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.manager.NewGroup(req, false); err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while creating group")
	}
	return &pb.Empty{}, nil
}
