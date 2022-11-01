package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetUserList gets a list of users
func (s *LDAPManagerService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.UserList, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	result, err := s.manager.GetUserList(req)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.Error); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error while getting list of accounts")
	}
	return result, nil
}
