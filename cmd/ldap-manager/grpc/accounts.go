package grpc

import (
	"context"

	log "github.com/sirupsen/logrus"
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
	account, err := s.Manager.GetAccount(in)
	if err != nil {
		log.Error(err)
		return &pb.User{}, status.Error(codes.Internal, "error while getting account")
	}
	return account, nil
}

// NewAccount ...
func (s *LDAPManagerServer) NewAccount(ctx context.Context, in *pb.NewAccountRequest) (*pb.Empty, error) {
	if err := s.Manager.NewAccount(in); err != nil {
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while creating new account")
	}
	return &pb.Empty{}, nil
}

// DeleteAccount ...
func (s *LDAPManagerServer) DeleteAccount(ctx context.Context, in *pb.DeleteAccountRequest) (*pb.Empty, error) {
	if err := s.Manager.DeleteAccount(in, false); err != nil {
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while deleting account")
	}
	return &pb.Empty{}, nil
}

// ChangePassword ...
func (s *LDAPManagerServer) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {
	if err := s.Manager.ChangePassword(in); err != nil {
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while chaning password of account")
	}
	return &pb.Empty{}, nil
}
