package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// GetGroupList returns a list of groups
func (s *LDAPManagerService) GetGroupList(ctx context.Context, in *pb.GetGroupListRequest) (*pb.GroupList, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.GroupList{}, err
	// }
	// groups, err := s.Manager.GetGroupList(in)
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
