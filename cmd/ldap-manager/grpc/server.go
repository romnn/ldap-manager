package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/romnn/go-service/pkg/auth"
	"github.com/romnn/go-service/pkg/grpc/reflect"

	ldapmanager "github.com/romnn/ldap-manager/pkg"
	ldaperror "github.com/romnn/ldap-manager/pkg/err"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/status"
)

func toStatus(e ldaperror.Error) error {
	return status.Error(e.Code(), e.Error())
}

// LDAPManagerService implements the GRPC service
type LDAPManagerService struct {
	pb.UnimplementedLDAPManagerServer

	manager       ldapmanager.LDAPManager
	authenticator auth.Authenticator

	server   *grpc.Server
	health   *health.Server
	registry reflect.Registry
}

// Shutdown gracefully stops the service
func (s *LDAPManagerService) Shutdown() {
	s.server.GracefulStop()
}

// NewLDAPManagerService builds the service
func NewLDAPManagerService(ctx context.Context, manager ldapmanager.LDAPManager, authenticator auth.Authenticator) LDAPManagerService {
	megabyte := 1024 * 1024
	maxMsgSize := 10 * megabyte

	registry := reflect.NewRegistry()
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		reflect.UnaryServerInterceptor(registry),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		reflect.StreamServerInterceptor(registry),
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
		grpc.MaxRecvMsgSize(maxMsgSize),
		grpc.MaxSendMsgSize(maxMsgSize),
	)

	health := health.NewServer()
	service := LDAPManagerService{
		manager:       manager,
		authenticator: authenticator,
		health:        health,
		server:        server,
		registry:      registry,
	}

	pb.RegisterLDAPManagerServer(server, &service)
	healthpb.RegisterHealthServer(server, health)
	registry.Load(server)

	return service
}

// Serve serves the service on a listener
func (s *LDAPManagerService) Serve(listener net.Listener) error {
	defer listener.Close()

	if err := s.server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve GRPC service: %v", err)
	}
	return nil
}
