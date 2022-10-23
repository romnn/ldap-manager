package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	// "sync"

	// ldapbase "github.com/romnn/ldap-manager/cmd/ldap-manager/base"
	// log "github.com/sirupsen/logrus"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	gw "github.com/romnn/ldap-manager/pkg/grpc/gen"
	"google.golang.org/grpc"
)

// LDAPManagerServer serves
type LDAPManagerService struct {
	// *ldapbase.LDAPManagerServer
	// Listener net.Listener
	upstream   *grpc.ClientConn
	gatewayMux *runtime.ServeMux
	server     *http.Server
	// SetupMux sync.Mutex
}

// Shutdown gracefully shuts down the service
func (s *LDAPManagerService) Shutdown() {
	s.server.Shutdown(context.Background())
	// s.SetupMux.Lock()
	// defer s.SetupMux.Unlock()
	// s.Upstream.Close()
	// if s.Listener != nil {
	// 	s.Listener.Close()
	// }
}

// func NewHTTPLDAPManagerServer(base *ldapbase.LDAPManagerServer, listener net.Listener, grpcConn *grpc.ClientConn) *LDAPManagerServer {

type Config struct {
	ServeStatic bool
	StaticPath  string
}

// NewLDAPManagerService builds the service
func NewLDAPManagerService(ctx context.Context, upstream *grpc.ClientConn, config *Config) (*LDAPManagerService, error) {
	gatewayMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case "X-User-Token":
				return "X-User-Token", true
			default:
				return key, false
			}
		}),
	)

	if err := gw.RegisterLDAPManagerHandler(ctx, gatewayMux, upstream); err != nil {
		return nil, err
	}

	rootMux := http.NewServeMux()
	if config.ServeStatic {
		// static frontend
		fileServer := http.FileServer(http.Dir(config.StaticPath))

		rootMux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := filepath.Join(config.StaticPath, r.URL.Path)
			_, err := os.Stat(path)
			if os.IsNotExist(err) {
				// file does not exist, serve index.html
				http.ServeFile(w, r, filepath.Join(config.StaticPath, "index.html"))
				return
			} else if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			fileServer.ServeHTTP(w, r)
		}))
	}

	// // health check
	// if s.LDAPManagerServer.Service.HTTPHealthCheckURL != "" {
	// 	rootMux.HandleFunc(s.LDAPManagerServer.Service.HTTPHealthCheckURL, func(w http.ResponseWriter, r *http.Request) {
	// 		if s.LDAPManagerServer.Service.Healthy {
	// 			w.WriteHeader(http.StatusOK)
	// 			w.Write([]byte("ok"))
	// 		} else {
	// 			w.WriteHeader(http.StatusServiceUnavailable)
	// 			w.Write([]byte("service is not available"))
	// 		}
	// 	})
	// }

	// gateway grpc API
	rootMux.Handle("/api/", http.StripPrefix("/api", gatewayMux))

	server := &http.Server{Handler: rootMux}

	return &LDAPManagerService{
		upstream:   upstream,
		gatewayMux: gatewayMux,
		server:     server,
	}, nil
}

// func (s *LDAPManagerServer) bootstrapHTTP(ctx context.Context) *http.ServeMux {
// 	rootMux := http.NewServeMux()
// 	if s.Static {
// 		// static frontend
// 		fileServer := http.FileServer(http.Dir(s.StaticRoot))
// 		rootMux.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			path := filepath.Join(s.StaticRoot, r.URL.Path)
// 			_, err := os.Stat(path)
// 			if os.IsNotExist(err) {
// 				// file does not exist, serve index.html
// 				http.ServeFile(w, r, filepath.Join(s.StaticRoot, "index.html"))
// 				return
// 			} else if err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}

// 			fileServer.ServeHTTP(w, r)
// 		}))
// 	}
// 	// health check
// 	if s.LDAPManagerServer.Service.HTTPHealthCheckURL != "" {
// 		rootMux.HandleFunc(s.LDAPManagerServer.Service.HTTPHealthCheckURL, func(w http.ResponseWriter, r *http.Request) {
// 			if s.LDAPManagerServer.Service.Healthy {
// 				w.WriteHeader(http.StatusOK)
// 				w.Write([]byte("ok"))
// 			} else {
// 				w.WriteHeader(http.StatusServiceUnavailable)
// 				w.Write([]byte("service is not available"))
// 			}
// 		})
// 	}
// 	// gateway grpc API
// 	rootMux.Handle("/api/", http.StripPrefix("/api", s.Mux))
// 	return rootMux
// }

// Serve serves the service on a listener
func (s *LDAPManagerService) Serve(listener net.Listener) error {
	defer listener.Close()

	// log.Printf("listening on: %v", listener.Addr())
	err := s.server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve HTTP service: %v", err)
	}
	return nil

	// s.SetupMux.Lock()
	// s.Service.HTTPServer = &http.Server{Handler: s.bootstrapHTTP(ctx)}
	// if err := gw.RegisterLDAPManagerHandler(ctx, s.Mux, s.Upstream); err != nil {
	// 	return err
	// }

	// s.Service.Ready = true
	// s.Service.SetHealthy(true)
	// log.Infof("%s ready at %s", s.Service.Name, s.Listener.Addr())
	// s.SetupMux.Unlock()

	// if err := s.Service.HTTPServer.Serve(s.Listener); err != nil && err != http.ErrServerClosed {
	// 	return err
	// }
	// log.Info("closing socket")
	// s.Listener.Close()
	// return nil
}
