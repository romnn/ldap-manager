package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// DeleteGroupMember removes a member of a group
func (s *LDAPManagerService) DeleteGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.Empty{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// allowDeleteOfDefaultGroups := claims.IsAdmin
	// if err := s.Manager.DeleteGroupMember(in, allowDeleteOfDefaultGroups); err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while deleting group member")
	// }
	return &pb.Empty{}, nil
}
