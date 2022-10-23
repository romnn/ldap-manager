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

// GetUser gets a user
func (s *LDAPManagerService) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserData, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.User{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.User{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// account, err := s.Manager.GetAccount(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldaperror.Error); safe {
	// 		return &pb.User{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.User{}, status.Error(codes.Internal, "error while getting account")
	// }
	// return account, nil
	return nil, nil
}
