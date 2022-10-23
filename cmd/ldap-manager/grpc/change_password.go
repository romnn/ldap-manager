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

// ChangePassword changes the password for an account
func (s *LDAPManagerService) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.Empty{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// if err := s.Manager.ChangePassword(in); err != nil {
	// 	if appErr, safe := err.(ldaperror.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while chaning password of account")
	// }
	return &pb.Empty{}, nil
}
