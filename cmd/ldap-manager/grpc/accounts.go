package grpc

import (
	"context"
	"strconv"

	ldapmanager "github.com/romnnn/ldap-manager"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	// "google.golang.org/grpc/metadata"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

// GetUserList ...
func (s *LDAPManagerServer) GetUserList(ctx context.Context, in *pb.GetUserListRequest) (*pb.UserList, error) {
	_, err := s.authenticate(ctx)
	if err != nil {
		return &pb.UserList{}, err
	}
	result, err := s.Manager.GetUserList(in)
	if err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.UserList{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.UserList{}, status.Error(codes.Internal, "error while getting list of accounts")
	}
	return result, nil
}

// GetAccount ...
func (s *LDAPManagerServer) GetAccount(ctx context.Context, in *pb.GetAccountRequest) (*pb.User, error) {
	claims, err := s.authenticate(ctx)
	if err != nil {
		return &pb.User{}, err
	}
	if !claims.IsAdmin && claims.UID != in.GetUsername() {
		return &pb.User{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	account, err := s.Manager.GetAccount(in)
	if err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.User{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.User{}, status.Error(codes.Internal, "error while getting account")
	}
	return account, nil
}

// NewAccount ...
func (s *LDAPManagerServer) NewAccount(ctx context.Context, in *pb.NewAccountRequest) (*pb.Empty, error) {
	_, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if err := s.Manager.NewAccount(in, pb.HashingAlgorithm_DEFAULT); err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Empty{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while creating new account")
	}
	return &pb.Empty{}, nil
}

// UpdateAccount ...
func (s *LDAPManagerServer) UpdateAccount(ctx context.Context, in *pb.UpdateAccountRequest) (*pb.Token, error) {
	claims, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Token{}, err
	}
	if !claims.IsAdmin && claims.UID != in.GetUsername() {
		return &pb.Token{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	username, uidNumber, err := s.Manager.UpdateAccount(in, pb.HashingAlgorithm_DEFAULT, claims.IsAdmin)
	if err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Token{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Token{}, status.Error(codes.Internal, "error while updating account")
	}
	token, expireSeconds, err := s.Authenticator.Login(&AuthClaims{
		UID:         username,
		UIDNumber:   strconv.Itoa(uidNumber),
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while signing token")
	}
	return &pb.Token{
		Token:       token,
		Username:    username,
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
		Expiration:  expireSeconds,
	}, nil
}

// DeleteAccount ...
func (s *LDAPManagerServer) DeleteAccount(ctx context.Context, in *pb.DeleteAccountRequest) (*pb.Empty, error) {
	claims, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	log.Info(claims.UID, in.GetUsername())
	if !claims.IsAdmin && claims.UID != in.GetUsername() {
		return &pb.Empty{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	allowDeleteOfDefaultGroups := false
	if err := s.Manager.DeleteAccount(in, allowDeleteOfDefaultGroups); err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Empty{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while deleting account")
	}
	return &pb.Empty{}, nil
}

// ChangePassword ...
func (s *LDAPManagerServer) ChangePassword(ctx context.Context, in *pb.ChangePasswordRequest) (*pb.Empty, error) {
	claims, err := s.authenticate(ctx)
	if err != nil {
		return &pb.Empty{}, err
	}
	if !claims.IsAdmin && claims.UID != in.GetUsername() {
		return &pb.Empty{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	}
	if err := s.Manager.ChangePassword(in); err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Empty{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Empty{}, status.Error(codes.Internal, "error while chaning password of account")
	}
	return &pb.Empty{}, nil
}
