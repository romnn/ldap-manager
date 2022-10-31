package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ChangePassword changes the password for an account
func (s *LDAPManagerService) ChangePassword(ctx context.Context, req *pb.ChangePasswordRequest) (*pb.Empty, error) {
	claims, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.IsAdmin && claims.UID != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	if err := s.manager.ChangePassword(req); err != nil {
		if appErr, safe := err.(ldaperror.Error); safe {
			return nil, appErr.StatusError()
		}
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while chaning password of account")
	}
	return &pb.Empty{}, nil
}
