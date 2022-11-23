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
	Username    string `json:"username"`
	UID         int32  `json:"uid"`
	IsAdmin     bool   `json:"is_admin"`
	DisplayName string `json:"display_name"`
	jwt.RegisteredClaims
}

// GetRegisteredClaims returns the common registered claims
func (claims *AuthClaims) GetRegisteredClaims() *jwt.RegisteredClaims {
	return &claims.RegisteredClaims
}

// MethodRequiresAdmin checks if the GRPC method requires authentication
func MethodRequiresAdmin(ctx context.Context) (bool, error) {
	if info, ok := reflect.GetMethodInfo(ctx); ok {
		methOptions := info.Method().Options()
		if requiresAdmin, ok := proto.GetExtension(
			methOptions,
			pb.E_RequireAdmin,
		).(bool); ok {
			return requiresAdmin, nil
		}
	}

	return true, errors.New("route has no or insufficient authentication policy")
}

// SignUserToken signs an authentication claim and returns it as a JWT token
func (s *LDAPManagerService) SignUserToken(claims *AuthClaims) (*pb.Token, error) {
	token, err := s.authenticator.SignJwtClaims(claims)
	if err != nil {
		log.Error(err)
		return nil, status.Error(
			codes.Internal,
			"error while signing token",
		)
	}
	expirationTime := time.Now().Add(s.authenticator.ExpiresAfter)
	// log.Debugf("signed token for %q: %s", claims.Username, token)
	return &pb.Token{
		Token:       token,
		Username:    claims.Username,
		UID:         claims.UID,
		IsAdmin:     claims.IsAdmin,
		DisplayName: claims.DisplayName,
		Expires:     timestamppb.New(expirationTime),
	}, nil
}

// Authenticate attempts to authenticate a user.
// It checks if a token is supplied in the request context
func (s *LDAPManagerService) Authenticate(ctx context.Context) (*AuthClaims, error) {
	requireAdmin, err := MethodRequiresAdmin(ctx)
	if err != nil {
		return nil, err
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(
			codes.Internal,
			"failed to extract authentication metadata",
		)
	}
	tokens := md.Get("x-user-token")
	if len(tokens) < 1 {
		return nil, status.Error(
			codes.Unauthenticated,
			"missing authentication token",
		)
	}
	valid, token, err := s.authenticator.Validate(
		tokens[0],
		&AuthClaims{},
	)
	if err != nil {
		log.Errorf("token=%s", tokens[0])
		log.Errorf("token validation failed: %v", err)
		return nil, status.Error(
			codes.Unauthenticated,
			"token validation failed",
		)
	}
	if claims, ok := token.Claims.(*AuthClaims); ok && valid {
		if requireAdmin && !claims.IsAdmin {
			return nil, status.Error(
				codes.PermissionDenied,
				"requires admin priviledges",
			)
		}
		// authenticated
		return claims, nil
	}
	return nil, status.Error(
		codes.Unauthenticated,
		"invalid token",
	)
}

// Login logs in a user
func (s *LDAPManagerService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.Token, error) {
	user, err := s.manager.AuthenticateUser(req)
	if err != nil {
		log.Error(err)
		if appErr, ok := err.(ldaperror.ApplicationError); ok {
			return &pb.Token{}, appErr.StatusError()
		}
		return nil, status.Error(
			codes.Unauthenticated, "unauthorized")
	}
	username := user.GetUsername()
	UID := user.GetUID()
	if username == "" || UID == 0 {
		return nil, status.Error(
			codes.NotFound,
			"user is invalid",
		)
	}

	memberStatus, err := s.manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: username,
		Group:    s.manager.DefaultAdminGroup,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Error(
			codes.Internal,
			"error checking user member status",
		)
	}
	return s.SignUserToken(&AuthClaims{
		Username:    username,
		UID:         UID,
		IsAdmin:     memberStatus.GetIsMember(),
		DisplayName: user.GetDisplayName(),
	})
}
