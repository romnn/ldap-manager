package grpc

import (
	"context"
	// "strconv"

	// ldapmanager "github.com/romnn/ldap-manager"
	// ldaperror "github.com/romnn/ldap-manager/pkg/err"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// NewUser creates a new user
func (s *LDAPManagerService) NewUser(ctx context.Context, in *pb.NewUserRequest) (*pb.Empty, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if err := s.Manager.NewAccount(in, pb.HashingAlgorithm_DEFAULT); err != nil {
	// 	if appErr, safe := err.(ldaperror.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while creating new account")
	// }
	return &pb.Empty{}, nil
}
