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
		&cli.IntFlag{
			Name:    "port",
			Value:   80,
			EnvVars: []string{"PORT"},
			Usage:   "service port",
		},
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
		&cli.StringFlag{
			Name:    "openldap-host",
			Value:   "localhost",
			EnvVars: []string{"OPENLDAP_HOST"},
			Usage:   "openldap host",
		},
		&cli.IntFlag{
			Name:    "openldap-port",
			Value:   389,
			EnvVars: []string{"OPENLDAP_PORT"},
			Usage:   "openldap port",
		},
		&cli.StringFlag{
			Name:    "openldap-protocol",
			Value:   "ldap",
			EnvVars: []string{"OPENLDAP_PROTOCOL"},
			Usage:   "openldap protocol",
		},
		&cli.StringFlag{
			Name:    "openldap-admin-password",
			Value:   "admin",
			EnvVars: []string{"OPENLDAP_ADMIN_PASSWORD"},
			Usage:   "openldap admin password",
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
						Host:                 ctx.String("openldap-host"),
						Port:                 ctx.Int("openldap-port"),
						Protocol:             ctx.String("openldap-protocol"),
						Organization:         "Example Inc.",
						Domain:               "example.org",
						BaseDN:               "dc=example,dc=org",
						AdminPassword:        ctx.String("openldap-admin-password"),
						ConfigPassword:       "config",
						ReadonlyUser:         true,
						ReadonlyUserUsername: "readonly",
						ReadonlyUserPassword: "readonly",
						TLS:                  false,
						UseRFC2307BISSchema:  true,
					},
					GroupsOU:                 "groups",
					UsersOU:                  "users",
					GroupsDN:                 "ou=groups,dc=example,dc=org",
					UserGroupDN:              "ou=users,dc=example,dc=org",
					GroupMembershipAttribute: "uniqueMember", // uniquemember or memberUID
					GroupMembershipUsesUID:   false,
					AccountAttribute:         "uid",
					GroupAttribute:           "gid",
					DefaultUserGroup:         "users",
					DefaultAdminGroup:        "admins",
					DefaultUserShell:         "/bin/bash",
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
	// Authentication
	s.echoServer.POST("/api/login", s.login)
	s.echoServer.POST("/api/logout", s.logout)

	// Account management (admin only)
	s.echoServer.POST("/api/accounts/list", s.listAccounts)
	s.echoServer.POST("/api/accounts/new", s.newAccount)

	// Group management (admin only)
	// s.echoServer.POST("/api/groups/list", s.listGroups)
	// s.echoServer.POST("/api/groups/new", s.newGroup)
	// s.echoServer.POST("/api/groups/:gid", s.getGroup)

	// Edit personal account
	s.echoServer.POST("/api/account/:username", s.getAccount)
	s.echoServer.POST("/api/account/delete/:username", s.deleteAccount)

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
