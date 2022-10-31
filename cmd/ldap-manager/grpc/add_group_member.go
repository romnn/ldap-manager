package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AddGroupMember adds a new member to a group
func (s *LDAPManagerService) AddGroupMember(ctx context.Context, req *pb.GroupMember) (*pb.Empty, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.manager.AddGroupMember(req, false); err != nil {
		log.Error(err)
		if appErr, safe := err.(ldaperror.Error); safe {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while adding group member")
	}
	return &pb.Empty{}, nil
}
