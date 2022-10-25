package grpc

import (
	"context"

	// ldapmanager "github.com/romnn/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/status"
)

// AddGroupMember adds a new member to a group
func (s *LDAPManagerService) AddGroupMember(ctx context.Context, in *pb.GroupMember) (*pb.Empty, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Empty{}, err
	// }
	// if err := s.Manager.AddGroupMember(in, false); err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Empty{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Empty{}, status.Error(codes.Internal, "error while adding group member")
	// }
	// return &pb.Empty{}, nil
	return nil, nil
}
