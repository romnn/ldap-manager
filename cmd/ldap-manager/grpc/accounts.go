package grpc

import (
	"context"

	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// GetUserList ...
func (s *LDAPManagerServer) GetUserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.UserList, error) {
	return s.LDAPManagerServer.GetUserList(ctx, in)
}

// AuthenticateUser ...
func (s *LDAPManagerServer) AuthenticateUser(ctx context.Context, in *pb.AuthenticateUserRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.AuthenticateUser(ctx, in)
}

// GetAccount ...
func (s *LDAPManagerServer) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.User, error) {
	return s.LDAPManagerServer.GetAccount(ctx, in)
}

// NewAccount ...
func (s *LDAPManagerServer) NewAccount(ctx context.Context, in *pb.NewAccountRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.NewAccount(ctx, in)
}

// DeleteAccount ...
func (s *LDAPManagerServer) DeleteAccount(ctx context.Context, in *pb.DeleteAccountRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.DeleteAccount(ctx, in)
}

// ChangePassword ...
func (s *LDAPManagerServer) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {
	return s.LDAPManagerServer.ChangePassword(ctx, in)
}
