package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	// "io"
	// "encoding/json"

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
			case "X-User-Token":
				return "X-User-Token", true
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

	s.Service.Ready = true
	s.Service.SetHealthy(true)
	log.Infof("%s ready at %s", s.Service.Name, listener.Addr())
}

func (s *LDAPManagerServer) bootstrapHTTP(ctx *cli.Context) *http.ServeMux {
	rootMux := http.NewServeMux()
	if ctx.Bool("static") {
		// static frontend
		staticRoot := ctx.String("static-root")
		fileServer := http.FileServer(http.Dir(staticRoot))
		rootMux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(staticRoot, r.URL.Path)
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				// file does not exist, serve index.html
				http.ServeFile(w, r, filepath.Join(staticRoot, "index.html"))
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileServer.ServeHTTP(w, r)
		}))
	}
	// health check
	rootMux.HandleFunc(s.LDAPManagerServer.Service.HTTPHealthCheckURL, func(w http.ResponseWriter, r *http.Request) {
		if s.LDAPManagerServer.Service.Healthy {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("service is not available"))
		}
	})
	// gateway grpc API
	rootMux.Handle("/api/", http.StripPrefix("/api", s.Mux))
	return rootMux
}

// Serve ...
func (s *LDAPManagerServer) Serve(wg *sync.WaitGroup, ctx *cli.Context) error {
	defer wg.Done()
	s.Service.HTTPServer = &http.Server{Handler: s.bootstrapHTTP(ctx)}

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
