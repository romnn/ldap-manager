package grpc

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"

	"github.com/romnn/go-service/pkg/grpc/reflect"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// Checks if the GRPC method requires authentication
func MethodRequiresAdmin(ctx context.Context) (bool, error) {
	if info, ok := reflect.GetMethodInfo(ctx); ok {
		methOptions := info.Method().Options()
		if requiresAdmin, ok := proto.GetExtension(methOptions, pb.E_RequireAdmin).(bool); ok {
			return requiresAdmin, nil
		}
	}

	return true, errors.New("route has no or insufficient authentication policy")
}

// Signs an authentication claim and returns a user token
func (s *LDAPManagerService) SignUserToken(claims *AuthClaims) (*pb.Token, error) {
	token, err := s.authenticator.SignJwtClaims(claims)
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while signing token")
	}
	expirationTime := time.Now().Add(s.authenticator.ExpiresAfter)
	return &pb.Token{
		Token:       token,
		Username:    claims.UID,
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
		Expires:     timestamppb.New(expirationTime),
	}, nil
}

// Attempts to authenticate a user if a token is supplied in the request context
func (s *LDAPManagerService) Authenticate(ctx context.Context) (*AuthClaims, error) {
	requireAdmin, err := MethodRequiresAdmin(ctx)
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
	user, err := s.manager.AuthenticateUser(in)
	if err != nil {
		if appErr, safe := err.(ldaperror.Error); safe {
			return &pb.Token{}, appErr.StatusError()
		}
		log.Error(err)
		return &pb.Token{}, status.Error(codes.Unauthenticated, "unauthorized")
	}
	uid := user.GetAttributeValue(s.manager.AccountAttribute)
	uidNumber := user.GetAttributeValue("uidNumber")
	if uid == "" || uidNumber == "" {
		return &pb.Token{}, status.Error(codes.NotFound, "user is invalid")
	}

	adminMemberStatus, err := s.manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: uid,
		Group:    s.manager.DefaultAdminGroup,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while checking user member status")
	}
	isAdmin := adminMemberStatus.GetIsMember()
	displayName := user.GetAttributeValue("displayName")
	return s.SignUserToken(&AuthClaims{
		UID:         uid,
		UIDNumber:   uidNumber,
		IsAdmin:     isAdmin,
		DisplayName: displayName,
	})
}
