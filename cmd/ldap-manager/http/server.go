package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/romnn/go-service/pkg/http/health"
	gw "github.com/romnn/ldap-manager/pkg/grpc/gen"

	"google.golang.org/grpc"
)

// LDAPManagerService implements the HTTP service
type LDAPManagerService struct {
	upstream *grpc.ClientConn

	gatewayMux *runtime.ServeMux
	server     *http.Server
	health     *health.Health
}

// Shutdown gracefully stops the service
func (s *LDAPManagerService) Shutdown() {
	s.server.Shutdown(context.Background())
}

// Config defines configuration options for the HTTP service
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

	router := http.NewServeMux()
	server := &http.Server{Handler: router}
	health := &health.Health{}
	service := &LDAPManagerService{
		upstream:   upstream,
		gatewayMux: gatewayMux,
		server:     server,
		health:     health,
	}

	// serve static files
	if config.ServeStatic {
		fileServer := http.FileServer(http.Dir(config.StaticPath))

		router.HandleFunc("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// health check
	router.Handle("/healthz", health)

	// grpc gateway
	router.Handle("/api/", http.StripPrefix("/api", gatewayMux))

	return service, nil
}

// Serve serves the service on a listener
func (s *LDAPManagerService) Serve(listener net.Listener) error {
	defer listener.Close()

	err := s.server.Serve(listener)
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to serve HTTP service: %v", err)
	}
	return nil
}
