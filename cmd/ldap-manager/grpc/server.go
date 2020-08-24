package grpc

import (
	"github.com/neko-neko/echo-logrus/v2/log"
	ldapbase "github.com/romnnn/ldap-manager/cmd/ldap-manager/base"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	"github.com/urfave/cli/v2"
)

// LDAPManagerServer ...
type LDAPManagerServer struct {
	pb.UnimplementedLDAPManagerServer
	*ldapbase.LDAPManagerServer
}

// NewGRPCLDAPManagerServer ...
func NewGRPCLDAPManagerServer(base *ldapbase.LDAPManagerServer) *LDAPManagerServer {
	return &LDAPManagerServer{
		LDAPManagerServer: base,
	}
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	if s.LDAPManagerServer != nil {
		s.LDAPManagerServer.Shutdown()
	}
}

// Serve ...
func (s *LDAPManagerServer) Serve(ctx *cli.Context) error {
	if err := s.Service.BootstrapGrpc(ctx, nil); err != nil {
		return err
	}
	go s.Connect(ctx)
	pb.RegisterLDAPManagerServer(s.Service.GrpcServer, s)
	if err := s.Service.ServeGrpc(s.Listener); err != nil {
		return err
	}
	log.Info("closing socket")
	s.Listener.Close()
	return nil
}
