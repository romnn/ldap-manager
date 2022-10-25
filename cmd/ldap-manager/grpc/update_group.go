package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// UpdateGroup updates an LDAP group
func (s *LDAPManagerService) UpdateGroup(ctx context.Context, in *pb.UpdateGroupRequest) (*pb.Empty, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if err := s.Manager.UpdateGroup(in); err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while updating group")
	// }
	return &pb.Empty{}, nil
}
