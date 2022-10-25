package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// GetUserGroups gets the groups an account is member of
func (s *LDAPManagerService) GetUserGroups(ctx context.Context, in *pb.GetUserGroupsRequest) (*pb.GroupList, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.GroupList{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.GroupList{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// groups, err := s.Manager.GetUserGroups(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.GroupList{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.GroupList{}, status.Error(codes.Internal, "error while getting groups")
	// }
	// return groups, nil
	return nil, nil
}
