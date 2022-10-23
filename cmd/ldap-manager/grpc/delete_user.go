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

// DeleteUser deletes an account
func (s *LDAPManagerService) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.Empty, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// log.Info(claims.UID, in.GetUsername())
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.Empty{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// allowDeleteOfDefaultGroups := false
	// if err := s.Manager.DeleteAccount(in, allowDeleteOfDefaultGroups); err != nil {
	// 	if appErr, safe := err.(ldaperror.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while deleting account")
	// }
	return &pb.Empty{}, nil
}
