package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUser gets a user
func (s *LDAPManagerService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	claims, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	username := req.GetUsername()
	if !claims.IsAdmin && claims.Username != username {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	user, err := s.manager.GetUser(username)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while getting account")
	}
	return user, nil
}
