package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUserGroups gets the groups an account is member of
func (s *LDAPManagerService) GetUserGroups(ctx context.Context, req *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	claims, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.IsAdmin && claims.Username != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	groups, err := s.manager.GetUserGroups(req)
	if err != nil {
		if appErr, ok := err.(ldaperror.Error); ok {
			return nil, appErr.StatusError()
		}
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while getting groups")
	}
	return groups, nil
}
