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

// GetUserList gets a list of users
func (s *LDAPManagerService) GetUserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.UserList, error) {
	// _, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.UserList{}, err
	// }
	// result, err := s.Manager.GetUserList(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldaperror.Error); safe {
	// 		return &pb.UserList{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.UserList{}, status.Error(codes.Internal, "error while getting list of accounts")
	// }
	// return result, nil
	return nil, nil
}
