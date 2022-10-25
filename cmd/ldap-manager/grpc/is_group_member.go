package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// IsGroupMember checks if an account is member of a group
func (s *LDAPManagerService) IsGroupMember(ctx context.Context, in *pb.IsGroupMemberRequest) (*pb.GroupMemberStatus, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.GroupMemberStatus{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.GroupMemberStatus{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// memberStatus, err := s.Manager.IsGroupMember(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.GroupMemberStatus{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.GroupMemberStatus{}, status.Error(codes.Internal, "error while checking if user is member")
	// }
	// return memberStatus, nil
	return nil, nil
}
