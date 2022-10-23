package grpc

import (
	"context"
	"errors"

	// "github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v4"
	// gogrpcservice "github.com/romnn/go-grpc-service"
	"github.com/romnn/go-service/pkg/grpc/reflect"
	// ldapmanager "github.com/romnn/ldap-manager"
	// ldapcli "github.com/romnn/ldap-manager/cmd/ldap-manager"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	// log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	// pref "google.golang.org/protobuf/reflect/protoreflect"
)

// AuthClaims encode authentication JWT claims
type AuthClaims struct {
	UID         string `json:"uid"`
	UIDNumber   string `json:"uid_num"`
	IsAdmin     bool   `json:"is_admin"`
	DisplayName string `json:"display_name"`
	jwt.RegisteredClaims
}

// GetRegisteredClaims ...
func (claims *AuthClaims) GetRegisteredClaims() *jwt.RegisteredClaims {
	return &claims.RegisteredClaims
}

func routeRequiresAdmin(ctx context.Context) (bool, error) {
	if info, ok := reflect.GetMethodInfo(ctx); ok {
		methOptions := info.Method().Options()
		if requiresAdmin, ok := proto.GetExtension(methOptions, pb.E_RequireAdmin).(bool); ok {
			return requiresAdmin, nil
		}
	}

	return true, errors.New("route has no or insufficient authentication policy")
}

// Login logs in a user
func (s *LDAPManagerService) authenticate(ctx context.Context) (*AuthClaims, error) {
	requireAdmin, err := routeRequiresAdmin(ctx)
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Internal, "failed to extract authentication metadata")
	}
	tokens := md.Get("x-user-token")
	if len(tokens) < 1 {
		return nil, status.Error(codes.Unauthenticated, "missing authentication token")
	}
	valid, token, err := s.authenticator.Validate(tokens[0], &AuthClaims{})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "token validation failed")
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && valid {
		if requireAdmin && !claims.IsAdmin {
			return nil, status.Error(codes.PermissionDenied, "requires admin priviledges")
		}
		// authenticated
		return claims, nil
	}
	return nil, status.Error(codes.Unauthenticated, "invalid token")
}

// Login logs in a user
func (s *LDAPManagerService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.Token, error) {
	// user, err := s.Manager.AuthenticateUser(in)
	// if err != nil {
	// 	if appErr, safe := err.(ldapmanager.Error); safe {
	// 		return &pb.Token{}, toStatus(appErr)
	// 	}
	// 	log.Error(err)
	// 	return &pb.Token{}, status.Error(codes.Unauthenticated, "unauthorized")
	// }
	// uid := user.GetAttributeValue(s.Manager.AccountAttribute)
	// uidNumber := user.GetAttributeValue("uidNumber")
	// if uid == "" || uidNumber == "" {
	// 	return &pb.Token{}, status.Error(codes.NotFound, "user is invalid")
	// }

	// adminMemberStatus, err := s.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
	// 	Username: uid,
	// 	Group:    s.Manager.DefaultAdminGroup,
	// })
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, status.Error(codes.Internal, "error while checking user member status")
	// }
	// isAdmin := adminMemberStatus.GetIsMember()
	// displayName := user.GetAttributeValue("displayName")
	// token, expireSeconds, err := s.Authenticator.Login(&AuthClaims{
	// 	UID:         uid,
	// 	UIDNumber:   uidNumber,
	// 	IsAdmin:     isAdmin,
	// 	DisplayName: displayName,
	// })
	// if err != nil {
	// 	log.Error(err)
	// 	return nil, status.Error(codes.Internal, "error while signing token")
	// }
	// return &pb.Token{
	// 	Token:       token,
	// 	Username:    uid,
	// 	IsAdmin:     isAdmin,
	// 	DisplayName: displayName,
	// 	Expiration:  expireSeconds,
	// }, nil
	return nil, nil
}
