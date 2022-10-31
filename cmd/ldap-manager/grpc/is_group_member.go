package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IsGroupMember checks if an account is member of a group
func (s *LDAPManagerService) IsGroupMember(ctx context.Context, req *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	claims, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if !claims.IsAdmin && claims.UID != req.GetUsername() {
		return nil, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	memberStatus, err := s.manager.IsGroupMember(req)
	if err != nil {
		if appErr, safe := err.(ldaperror.Error); safe {
			return nil, appErr.StatusError()
		}
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while checking if user is member")
	}
	return memberStatus, nil
}
