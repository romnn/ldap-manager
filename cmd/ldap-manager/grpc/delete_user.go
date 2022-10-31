package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DeleteUser deletes an account
func (s *LDAPManagerService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.Empty, error) {
	claims, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	log.Info(claims.UID, req.GetUsername())

	if !claims.IsAdmin && claims.UID != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	allowDeleteOfDefaultGroups := false
	if err := s.manager.DeleteUser(req, allowDeleteOfDefaultGroups); err != nil {
		log.Error(err)
		if appErr, safe := err.(ldaperror.Error); safe {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while deleting account")
	}
	return &pb.Empty{}, nil
}
