package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateGroup updates an LDAP group
func (s *LDAPManagerService) UpdateGroup(ctx context.Context, in *pb.UpdateGroupRequest) (*pb.Empty, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.manager.UpdateGroup(in); err != nil {
		log.Error(err)
		if appErr, safe := err.(ldaperror.Error); safe {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while updating group")
	}
	return &pb.Empty{}, nil
}
