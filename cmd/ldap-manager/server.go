package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// "crypto/tls"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	ldapmanager "github.com/romnnn/ldap-manager"
	ldapconfig "github.com/romnnn/ldap-manager/config"

	"github.com/romnnn/flags4urfavecli/flags"
	// "github.com/romnnn/flags4urfavecli/values"

	gogrpcservice "github.com/romnnn/go-grpc-service"
	"github.com/romnnn/go-grpc-service/versioning"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Rev is set on build time to the git HEAD
var Rev = ""

// Version is incremented using bump2version
const Version = "0.0.1"

var server LDAPManagerServer

// LDAPManagerServer ...
type LDAPManagerServer struct {
	gogrpcservice.Service
	echoServer *echo.Echo
	manager    *ldapmanager.LDAPManager
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	s.Service.GracefulStop()
	s.echoServer.Shutdown(context.Background())
	if s.manager != nil {
		s.manager.Close()
	}
}

func main() {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		server.Shutdown()
	}()

	cliFlags := []cli.Flag{
		&flags.LogLevelFlag,
		/*
			&cli.GenericFlag{
				Name: "format",
				Value: &values.EnumValue{
					Enum:    []string{"json", "xml", "csv"},
					Default: "xml",
				},
				EnvVars: []string{"FILEFORMAT"},
				Usage:   "input file format",
			},
		*/
		&cli.StringFlag{
			Name:    "ldap-uri",
			Value:   "ldap://localhost:389",
			EnvVars: []string{"LDAP_URI", "LDAP_CONNECTION_URI"},
			Usage:   "ldap connection URI",
		},
		&cli.StringFlag{
			Name:    "user-group-dn",
			Value:   "ldap://localhost:389",
			EnvVars: []string{"LDAP_URI", "LDAP_CONNECTION_URI"},
			Usage:   "ldap connection URI",
		},
		&cli.IntFlag{
			Name:    "port",
			Value:   80,
			EnvVars: []string{"PORT"},
			Usage:   "service port",
		},
	}

	name := "ldap manager service"

	app := &cli.App{
		Name:    name,
		Version: versioning.BinaryVersion(Version, Rev),
		Usage:   "manages ldap user accounts",
		Flags:   cliFlags,
		Action: func(ctx *cli.Context) error {
			server = LDAPManagerServer{
				Service: gogrpcservice.Service{
					Name:               name,
					Version:            Version,
					BuildTime:          Rev,
					HTTPHealthCheckURL: "health/healthz",
				},
				manager: &ldapmanager.LDAPManager{
					OpenLDAPConfig: ldapconfig.OpenLDAPConfig{
						Host:          "localhost",
						Port:          ctx.Int("port"),
						Protocol:      "ldap",
						AdminPassword: "admin",
					},
					GroupsOU: "groups",
					UsersOU:  "users",
					// BaseDN:                   "dc=example,dc=org",
					GroupsDN:                 "ou=groups,dc=example,dc=org",
					UserGroupDN:              "ou=users,dc=example,dc=org",
					GroupMembershipAttribute: "uniqueMember", // uniquemember or memberUID
					GroupMembershipUsesUID:   false,
					// AdminBindUsername:        "admin",
					// AdminBindPassword:        "admin",
					// ReadonlyBindUsername:     "readonly",
					// ReadonlyBindPassword:     "readonly",
					AccountAttribute:  "uid",
					GroupAttribute:    "gid",
					DefaultUserGroup:  "users",
					DefaultAdminGroup: "admins",
					DefaultUserShell:  "/bin/bash",
					// RequireStartTLS:          false,
				},
			}

			port := fmt.Sprintf(":%d", ctx.Int("port"))
			listener, err := net.Listen("tcp", port)
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}

			server.BootstrapHTTP(ctx)
			return server.Serve(ctx, listener)
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// BootstrapHTTP prepares an http service
func (s *LDAPManagerServer) BootstrapHTTP(cliCtx *cli.Context) error {
	s.echoServer = echo.New()
	// s.echoServer.Use(middleware.Logger())
	s.echoServer.Use(middleware.Recover())

	s.echoServer.GET("/healthz", func(c echo.Context) error {
		if s.Service.Healthy {
			c.String(http.StatusOK, "ok")
		} else {
			c.String(http.StatusServiceUnavailable, "service is not available")
		}
		return nil
	})
	s.echoServer.POST("/api/login", s.login)
	s.echoServer.POST("/api/logout", s.logout)
	s.echoServer.Static("/", "./frontend/dist")

	return s.Service.Bootstrap(cliCtx)
}

// Setup prepares the service
func (s *LDAPManagerServer) Setup(ctx *cli.Context) error {
	if err := s.manager.Setup(); err != nil {
		return err
	}
	return nil
}

// Serve starts the service
func (s *LDAPManagerServer) Serve(ctx *cli.Context, listener net.Listener) error {

	go func() {
		log.Info("connecting...")
		if err := s.Setup(ctx); err != nil {
			log.Error(err)
			s.Shutdown()
			return
		}
		s.Service.Ready = true
		s.Service.SetHealthy(true)
		log.Infof("%s ready at %s", s.Service.Name, listener.Addr())
	}()

	s.echoServer.Listener = listener
	s.echoServer.TLSListener = listener
	if err := s.echoServer.Start(""); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("closing socket")
	listener.Close()
	return nil
}
