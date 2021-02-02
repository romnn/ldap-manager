package http

import (
	"context"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	ldapbase "github.com/romnn/ldap-manager/cmd/ldap-manager/base"
	log "github.com/sirupsen/logrus"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/romnn/ldap-manager/grpc/ldap-manager"
	"google.golang.org/grpc"
)

// LDAPManagerServer ...
type LDAPManagerServer struct {
	*ldapbase.LDAPManagerServer
	Listener net.Listener
	Upstream *grpc.ClientConn
	Mux      *runtime.ServeMux
	SetupMux sync.Mutex
}

// NewHTTPLDAPManagerServer ...
func NewHTTPLDAPManagerServer(base *ldapbase.LDAPManagerServer, listener net.Listener, grpcConn *grpc.ClientConn) *LDAPManagerServer {
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
		Upstream:          grpcConn,
		Mux:               mux,
	}
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	s.SetupMux.Lock()
	defer s.SetupMux.Unlock()
	s.Upstream.Close()
	if s.Listener != nil {
		s.Listener.Close()
	}
}

func (s *LDAPManagerServer) bootstrapHTTP(ctx context.Context) *http.ServeMux {
	rootMux := http.NewServeMux()
	if s.Static {
		// static frontend
		fileServer := http.FileServer(http.Dir(s.StaticRoot))
		rootMux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(s.StaticRoot, r.URL.Path)
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				// file does not exist, serve index.html
				http.ServeFile(w, r, filepath.Join(s.StaticRoot, "index.html"))
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileServer.ServeHTTP(w, r)
		}))
	}
	// health check
	if s.LDAPManagerServer.Service.HTTPHealthCheckURL != "" {
		rootMux.HandleFunc(s.LDAPManagerServer.Service.HTTPHealthCheckURL, func(w http.ResponseWriter, r *http.Request) {
			if s.LDAPManagerServer.Service.Healthy {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok"))
			} else {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte("service is not available"))
			}
		})
	}
	// gateway grpc API
	rootMux.Handle("/api/", http.StripPrefix("/api", s.Mux))
	return rootMux
}

// Serve ...
func (s *LDAPManagerServer) Serve(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	s.SetupMux.Lock()
	s.Service.HTTPServer = &http.Server{Handler: s.bootstrapHTTP(ctx)}
	if err := gw.RegisterLDAPManagerHandler(ctx, s.Mux, s.Upstream); err != nil {
		return err
	}

	s.Service.Ready = true
	s.Service.SetHealthy(true)
	log.Infof("%s ready at %s", s.Service.Name, s.Listener.Addr())
	s.SetupMux.Unlock()

	if err := s.Service.HTTPServer.Serve(s.Listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("closing socket")
	s.Listener.Close()
	return nil
}
