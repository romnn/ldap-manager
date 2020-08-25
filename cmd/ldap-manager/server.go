package main

import (
	// "context"
	// "fmt"
	// "net"
	// "net/http"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	// "time"

	// "crypto/tls"

	// "github.com/labstack/echo/v4"
	ldapbase "github.com/romnnn/ldap-manager/cmd/ldap-manager/base"
	ldapgrpc "github.com/romnnn/ldap-manager/cmd/ldap-manager/grpc"
	ldaphttp "github.com/romnnn/ldap-manager/cmd/ldap-manager/http"

	// ldapmanager "github.com/romnnn/ldap-manager"
	// ldapconfig "github.com/romnnn/ldap-manager/config"

	"github.com/romnnn/flags4urfavecli/flags"
	// "github.com/romnnn/flags4urfavecli/values"

	// gogrpcservice "github.com/romnnn/go-grpc-service"
	// log "github.com/sirupsen/logrus"
	"github.com/romnnn/go-grpc-service/versioning"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	grpcServer LDAPManagerServer
	httpServer LDAPManagerServer
)

// LDAPManagerServer ...
type LDAPManagerServer interface {
	Serve(*sync.WaitGroup, *cli.Context) error
	Shutdown()
}

func main() {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		grpcServer.Shutdown()
		httpServer.Shutdown()
	}()

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

	serverFlags := []cli.Flag{
		&flags.LogLevelFlag,
		&cli.IntFlag{
			Name:    "grpc-port",
			Value:   9090,
			EnvVars: []string{"GRPC_PORT"},
			Usage:   "grpc service port",
		},
		&cli.IntFlag{
			Name:    "http-port",
			Value:   80,
			Aliases: []string{"port"},
			EnvVars: []string{"HTTP_PORT", "PORT"},
			Usage:   "http service port",
		},
	}

	configFlags := []cli.Flag{
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
		Version: versioning.BinaryVersion(ldapbase.Version, ldapbase.Rev),
		Usage:   "manages ldap user accounts",
		Flags:   configFlags,
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "serve ldap manager service",
				Flags: serverFlags,
				Action: func(ctx *cli.Context) error {
					grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctx.Int("grpc-port")))
					if err != nil {
						return fmt.Errorf("failed to listen: %v", err)
					}
					httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", ctx.Int("http-port")))
					if err != nil {
						return fmt.Errorf("failed to listen: %v", err)
					}

					base := ldapbase.NewLDAPManagerServer(ctx)
					grpcServer = ldapgrpc.NewGRPCLDAPManagerServer(base, grpcListener)
					httpServer = ldaphttp.NewHTTPLDAPManagerServer(base, httpListener, grpcListener)
					var wg sync.WaitGroup
					wg.Add(2)
					go grpcServer.Serve(&wg, ctx)
					go httpServer.Serve(&wg, ctx)
					wg.Wait()
					return nil
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
