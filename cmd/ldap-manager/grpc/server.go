package grpc

import (
	"context"
	"net"
	"sync"

	gogrpcservice "github.com/romnn/go-grpc-service"
	ldapmanager "github.com/romnn/ldap-manager"
	ldapbase "github.com/romnn/ldap-manager/cmd/ldap-manager/base"
	pb "github.com/romnn/ldap-manager/grpc/ldap-manager"
	log "github.com/sirupsen/logrus"

	"google.golang.org/grpc/status"
)

// LDAPManagerServer ...
type LDAPManagerServer struct {
	pb.UnimplementedLDAPManagerServer
	*ldapbase.LDAPManagerServer
	Listener net.Listener
	SetupMux sync.Mutex
}

// NewGRPCLDAPManagerServer ...
func NewGRPCLDAPManagerServer(base *ldapbase.LDAPManagerServer, listener net.Listener) *LDAPManagerServer {
	return &LDAPManagerServer{
		LDAPManagerServer: base,
		Listener:          listener,
	}
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	s.SetupMux.Lock()
	defer s.SetupMux.Unlock()
	if s.LDAPManagerServer != nil {
		s.LDAPManagerServer.Shutdown()
	}
}

// Serve ...
func (s *LDAPManagerServer) Serve(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	s.SetupMux.Lock()
	if err := s.Service.BootstrapGrpc(ctx, nil, &gogrpcservice.BootstrapGrpcOptions{}); err != nil {
		return err
	}
	s.LDAPManagerServer.Connect(ctx, s.Listener)
	pb.RegisterLDAPManagerServer(s.Service.GrpcServer, s)
	s.SetupMux.Unlock()
	if err := s.Service.ServeGrpc(s.Listener); err != nil {
		return err
	}
	log.Info("closing socket")
	s.Listener.Close()
	return nil
}

func toStatus(e ldapmanager.Error) error {
	return status.Error(e.Code(), e.Error())
}
