package grpc

import (
	"context"
	"errors"

	"github.com/dgrijalva/jwt-go"
	gogrpcservice "github.com/romnnn/go-grpc-service"
	ldapmanager "github.com/romnnn/ldap-manager"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// AuthClaims ...
type AuthClaims struct {
	UID         string `json:"uid"`
	UIDNumber   string `json:"uid_num"`
	IsAdmin     bool   `json:"is_admin"`
	DisplayName string `json:"display_name"`
	jwt.StandardClaims
}

// GetStandardClaims ...
func (claims *AuthClaims) GetStandardClaims() *jwt.StandardClaims {
	// the authenticator will use this method to get the standard claims that will be set based on the config
	// note that it is important to return a pointer to the current claims' standard claims and not any ones
	return &claims.StandardClaims
}

func routeRequiresAdmin(ctx context.Context) (bool, error) {
	if methodDesc, ok := ctx.Value(gogrpcservice.GrpcMethodDescriptor).(pref.MethodDescriptor); ok {
		if requireAdmin, ok := proto.GetExtension(methodDesc.Options(), pb.E_RequireAdmin).(bool); ok {
			return requireAdmin, nil
		}
	}
	return true, errors.New("route has no or insufficient authentication policy")
}

// Login logs in a user
func (s *LDAPManagerServer) authenticate(ctx context.Context) (*AuthClaims, error) {
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
	valid, token, err := s.Authenticator.Validate(tokens[0], &AuthClaims{})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
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
func (s *LDAPManagerServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.Token, error) {
	user, err := s.Manager.AuthenticateUser(in)
	if err != nil {
		if appErr, safe := err.(ldapmanager.Error); safe {
			return &pb.Token{}, toStatus(appErr)
		}
		log.Error(err)
		return &pb.Token{}, status.Error(codes.Unauthenticated, "unauthorized")
	}
	uid := user.GetAttributeValue(s.Manager.AccountAttribute)
	uidNumber := user.GetAttributeValue("uidNumber")
	if uid == "" || uidNumber == "" {
		return &pb.Token{}, status.Error(codes.NotFound, "user is invalid")
	}

	adminMemberStatus, err := s.Manager.IsGroupMember(&pb.IsGroupMemberRequest{
		Username: uid,
		Group:    s.Manager.DefaultAdminGroup,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while checking user member status")
	}
	isAdmin := adminMemberStatus.GetIsMember()
	displayName := user.GetAttributeValue("displayName")
	token, expireSeconds, err := s.Authenticator.Login(&AuthClaims{
		UID:         uid,
		UIDNumber:   uidNumber,
		IsAdmin:     isAdmin,
		DisplayName: displayName,
	})
	if err != nil {
		log.Error(err)
		return nil, status.Error(codes.Internal, "error while signing token")
	}
	return &pb.Token{
		Token:       token,
		Username:    uid,
		IsAdmin:     isAdmin,
		DisplayName: displayName,
		Expiration:  expireSeconds,
	}, nil
}
