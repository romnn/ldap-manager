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

// UpdateUser updates a user
func (s *LDAPManagerService) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.Token, error) {
	// claims, err := s.authenticate(ctx)
	// if err != nil {
	// 	return &pb.Token{}, err
	// }
	// if !claims.IsAdmin && claims.UID != in.GetUsername() {
	// 	return &pb.Token{}, status.Error(codes.PermissionDenied, "requires admin privileges")
	// }
	// username, uidNumber, err := s.Manager.UpdateAccount(in, pb.HashingAlgorithm_DEFAULT, claims.IsAdmin)
	// if err != nil {
	// 	if appErr, safe := err.(ldaperror.Error); safe {
	// 		return &pb.Token{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Token{}, status.Error(codes.Internal, "error while updating account")
	// }
	// token, expireSeconds, err := s.Authenticator.Login(&AuthClaims{
	// 	UID:         username,
	// 	UIDNumber:   strconv.Itoa(uidNumber),
	// 	IsAdmin:     claims.IsAdmin,
	// 	DisplayName: claims.DisplayName,
	// })
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, status.Error(codes.Internal, "error while signing token")
	// }
	// return &pb.Token{
	// 	Token:       token,
	// 	Username:    username,
	// 	IsAdmin:     claims.IsAdmin,
	// 	DisplayName: claims.DisplayName,
	// 	Expiration:  expireSeconds,
	// }, nil
	return nil, nil
}
