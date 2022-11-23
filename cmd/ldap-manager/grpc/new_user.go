package grpc

import (
	"context"

	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewUser creates a new user
func (s *LDAPManagerService) NewUser(ctx context.Context, req *pb.NewUserRequest) (*pb.Empty, error) {
	_, err := s.Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.manager.NewUser(req); err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return nil, appErr.StatusError()
		}
		return nil, status.Error(codes.Internal, "error creating new account")
	}
	return &pb.Empty{}, nil
}
