package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/romnn/go-service/pkg/auth"
	"github.com/romnn/go-service/pkg/grpc/reflect"
	ldapmanager "github.com/romnn/ldap-manager/pkg"

	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

// ServiceName is the name of the service used for health checking
const ServiceName = "ldap-manager"

// LDAPManagerService implements the GRPC service
type LDAPManagerService struct {
	pb.UnimplementedLDAPManagerServer

	manager       ldapmanager.LDAPManager
	authenticator auth.Authenticator

	server   *grpc.Server
	health   *health.Server
	registry reflect.Registry
}

// SetHealthy sets the health state for the service
func (s *LDAPManagerService) SetHealthy(healthy bool) {
	if s.health == nil {
		return
	}

	// assumes SetServingStatus is thread-safe
	if healthy {
		s.health.SetServingStatus(ServiceName, healthpb.HealthCheckResponse_SERVING)
	} else {
		s.health.SetServingStatus(ServiceName, healthpb.HealthCheckResponse_NOT_SERVING)
	}
}

// Shutdown gracefully stops the service
func (s *LDAPManagerService) Shutdown() {
	s.health.Shutdown()
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
