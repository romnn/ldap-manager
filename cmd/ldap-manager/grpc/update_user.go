package grpc

import (
	"context"
	"strconv"

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
	if !claims.IsAdmin && claims.UID != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	username, uidNumber, err := s.manager.UpdateUser(req, pb.HashingAlgorithm_DEFAULT, claims.IsAdmin)
	if err != nil {
		log.Error(err)
		if appErr, safe := err.(ldaperror.Error); safe {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while updating account")
	}
	return s.SignUserToken(&AuthClaims{
		UID:         username,
		UIDNumber:   strconv.Itoa(uidNumber),
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
	})
}
