package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UpdateUser updates a user
func (s *LDAPManagerService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.Token, error) {
	claims, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.IsAdmin && claims.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	username, err := s.manager.UpdateUser(req, claims.IsAdmin)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error updating user")
	}
	user, err := s.manager.GetUser(username)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error getting user")
	}

	return s.SignUserToken(&AuthClaims{
		Username:    user.GetUsername(),
		UID:         user.GetUID(),
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
	})
}
