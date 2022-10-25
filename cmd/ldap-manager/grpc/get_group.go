package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	// ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// GetGroup gets a group
func (s *LDAPManagerService) GetGroup(ctx context.Context, in *pb.GetGroupRequest) (*pb.Group, error) {
	claims, err := s.authenticate(ctx)
  log.Infof("claims: %v", claims)
	if err != nil {
		return &pb.Group{}, err
	}
	// group, err := s.manager.GetGroup(in)
	// if err != nil {
	// 	if err, safe := err.(ldaperror.Error); safe {
	// 		return &pb.Group{}, toStatus(err)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Group{}, status.Error(codes.Internal, "error while getting group")
	// }
	// return group, nil
	return &pb.Group{}, nil
}
