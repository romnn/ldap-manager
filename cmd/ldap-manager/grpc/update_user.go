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
	username, uidNumber, err := s.manager.UpdateUser(req, claims.IsAdmin)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.Error); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while updating account")
	}
	return s.SignUserToken(&AuthClaims{
		Username:    username,
		UID:         uidNumber,
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
	})
}
