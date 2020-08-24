package base

import (
	"context"

	"github.com/neko-neko/echo-logrus/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "google.golang.org/grpc/metadata"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// GetUserList ...
func (s *LDAPManagerServer) GetUserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.UserList, error) {
	result, err := s.Manager.GetUserList(in)
	if err != nil {
		log.Error(err)
		return &pb.UserList{}, status.Error(codes.Internal, "error while getting list of accounts")
	}
	return result, nil
}

// AuthenticateUser ...
func (s *LDAPManagerServer) AuthenticateUser(ctx context.Context, in *pb.AuthenticateUserRequest) (*pb.Empty, error) {
	if err := s.Manager.AuthenticateUser(in); err != nil {
		return &pb.Empty{}, status.Error(codes.Internal, "error while getting list of accounts")
	}
	return &pb.Empty{}, nil
}

// GetAccount ...
func (s *LDAPManagerServer) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.User, error) {
	var result pb.User
	log.Info("GetAccount")
	return &result, nil
}

// NewAccount ...
func (s *LDAPManagerServer) NewAccount(ctx context.Context, in *pb.NewAccountRequest) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("NewAccount")
	return &result, nil
}

// DeleteAccount ...
func (s *LDAPManagerServer) DeleteAccount(ctx context.Context, in *pb.DeleteAccountRequest) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("DeleteAccount")
	return &result, nil
}

// ChangePassword ...
func (s *LDAPManagerServer) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {
	var result pb.Empty
	log.Info("ChangePassword")
	return &result, nil
}
