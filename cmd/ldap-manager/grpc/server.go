package grpc

import (
	"net"
	"sync"

	ldapmanager "github.com/romnnn/ldap-manager"
	ldapbase "github.com/romnnn/ldap-manager/cmd/ldap-manager/base"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc/status"
)

// LDAPManagerServer ...
type LDAPManagerServer struct {
	pb.UnimplementedLDAPManagerServer
	*ldapbase.LDAPManagerServer
	Listener net.Listener
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
	if s.LDAPManagerServer != nil {
		s.LDAPManagerServer.Shutdown()
	}
}

func toStatus(e ldapmanager.Error) error {
	return status.Error(e.Code(), e.Error())
}

// Serve ...
func (s *LDAPManagerServer) Serve(wg *sync.WaitGroup, ctx *cli.Context) error {
	defer wg.Done()
	if err := s.Service.BootstrapGrpc(ctx, nil); err != nil {
		return err
	}
	go s.Connect(ctx, s.Listener)
	pb.RegisterLDAPManagerServer(s.Service.GrpcServer, s)
	if err := s.Service.ServeGrpc(s.Listener); err != nil {
		return err
	}
	log.Info("closing socket")
	s.Listener.Close()
	return nil
}
