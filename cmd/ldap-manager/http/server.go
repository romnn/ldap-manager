package http

import (
	"context"
	"net"
	"net/http"
	"sync"

	ldapbase "github.com/romnnn/ldap-manager/cmd/ldap-manager/base"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/romnnn/ldap-manager/grpc/ldap-manager"
	"google.golang.org/grpc"
)

// LDAPManagerServer ...
type LDAPManagerServer struct {
	*ldapbase.LDAPManagerServer
	Listener     net.Listener
	GRPCListener net.Listener
	Mux          *runtime.ServeMux
}

// NewHTTPLDAPManagerServer ...
func NewHTTPLDAPManagerServer(base *ldapbase.LDAPManagerServer, listener, grpcListener net.Listener) *LDAPManagerServer {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "X-Custom-Header2":
				return "custom-header2", true
			default:
				return key, false
			}
		}),
	)
	return &LDAPManagerServer{
		LDAPManagerServer: base,
		Listener:          listener,
		GRPCListener:      grpcListener,
		Mux:               mux,
	}
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	if s.LDAPManagerServer != nil {
		s.LDAPManagerServer.Shutdown()
	}
}

// Connect ...
func (s *LDAPManagerServer) Connect(ctx *cli.Context, listener net.Listener) {
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	if err := gw.RegisterLDAPManagerHandlerFromEndpoint(context.Background(), s.Mux, s.GRPCListener.Addr().String(), opts); err != nil {
		log.Error(err)
		s.Shutdown()
	}
	s.LDAPManagerServer.Connect(ctx, listener)
}

// Serve ...
func (s *LDAPManagerServer) Serve(wg *sync.WaitGroup, ctx *cli.Context) error {
	defer wg.Done()
	s.Service.HTTPServer = &http.Server{Handler: s.Mux}
	if err := s.Service.Bootstrap(ctx); err != nil {
		return err
	}

	go s.Connect(ctx, s.Listener)
	if err := s.Service.HTTPServer.Serve(s.Listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("closing socket")
	s.Listener.Close()
	return nil
}
