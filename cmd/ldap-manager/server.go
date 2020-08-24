package main

import (
	// "context"
	// "fmt"
	// "net"
	// "net/http"
	"os"
	"os/signal"
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
	"github.com/neko-neko/echo-logrus/v2/log"
	"github.com/romnnn/go-grpc-service/versioning"
	"github.com/urfave/cli/v2"
)

var server LDAPManagerServer

// LDAPManagerServer ...
type LDAPManagerServer interface {
	Serve(*cli.Context) error
	Shutdown()
}

func main() {
	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		server.Shutdown()
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
			Name:    "port",
			Value:   80,
			EnvVars: []string{"PORT"},
			Usage:   "service port",
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
				Subcommands: []*cli.Command{
					{
						Name:  "grpc",
						Usage: "serve ldap manager grpc service",
						Action: func(ctx *cli.Context) error {
							base, err := ldapbase.NewLDAPManagerServer(ctx)
							if err != nil {
								return err
							}
							server = ldapgrpc.NewGRPCLDAPManagerServer(base)
							return server.Serve(ctx)
						},
					},
					{
						Name:  "http",
						Usage: "serve ldap manager http service",
						Action: func(ctx *cli.Context) error {
							base, err := ldapbase.NewLDAPManagerServer(ctx)
							if err != nil {
								return err
							}
							server = ldaphttp.NewHTTPLDAPManagerServer(base)
							return server.Serve(ctx)
						},
					},
				},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
