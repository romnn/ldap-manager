package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/romnn/go-grpc-service/auth"
	ldapbase "github.com/romnn/ldap-manager/cmd/ldap-manager/base"
	ldapgrpc "github.com/romnn/ldap-manager/cmd/ldap-manager/grpc"
	ldaphttp "github.com/romnn/ldap-manager/cmd/ldap-manager/http"

	"github.com/romnn/flags4urfavecli/flags"
	"github.com/romnn/flags4urfavecli/values"
	"github.com/romnn/go-grpc-service/versioning"
	ldapmanager "github.com/romnn/ldap-manager"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var (
	grpcServer LDAPManagerServer
	httpServer LDAPManagerServer
)

// LDAPManagerServer ...
type LDAPManagerServer interface {
	Serve(context.Context, *sync.WaitGroup) error
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
		&cli.BoolFlag{
			Name:    "no-static",
			Value:   false,
			Aliases: []string{"disable-serve-static"},
			EnvVars: []string{"NO_STATIC", "DISABLE_SERVE_STATIC"},
			Usage:   "disable serving of the static frontend",
		},
		&cli.StringFlag{
			Name:    "static-root",
			Value:   "./frontend/dist",
			EnvVars: []string{"STATIC_DIR", "STATIC_ROOT"},
			Usage:   "root source directory of the static files to be served",
		},
	}

	jwtAuthFlags := auth.DefaultCLIFlags(&auth.DefaultCLIFlagsOptions{
		Issuer:    "issuer@example.org",
		Audience:  "example.org",
		ExpireSec: 1 * 24 * 60 * 60,
	})

	ldapConfigFlags := []cli.Flag{
		// Connection
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
		&cli.StringFlag{
			Name:    "openldap-config-password",
			Value:   "config",
			EnvVars: []string{"OPENLDAP_CONFIG_PASSWORD"},
			Usage:   "openldap config password",
		},
		&cli.StringFlag{
			Name:    "openldap-readonly-user",
			Value:   "", // no readonly user
			EnvVars: []string{"OPENLDAP_READONLY_USER"},
			Usage:   "openldap readonly user",
		},
		&cli.StringFlag{
			Name:    "openldap-readonly-password",
			Value:   "", // no readonly user
			EnvVars: []string{"OPENLDAP_READONLY_PASSWORD"},
			Usage:   "openldap readonly password",
		},
		&cli.StringFlag{
			Name:    "openldap-organization",
			Value:   "Example Inc.",
			EnvVars: []string{"OPENLDAP_ORGANIZATION"},
			Usage:   "openldap organization",
		},
		&cli.StringFlag{
			Name:    "openldap-domain",
			Value:   "example.org",
			EnvVars: []string{"OPENLDAP_DOMAIN"},
			Usage:   "openldap domain",
		},
		&cli.StringFlag{
			Name:    "openldap-base-dn",
			Value:   "dc=example,dc=org",
			EnvVars: []string{"OPENLDAP_BASE_DN"},
			Usage:   "openldap base DN",
		},
		&cli.BoolFlag{
			Name:    "openldap-tls",
			Value:   false,
			EnvVars: []string{"OPENLDAP_TLS"},
			Usage:   "openldap tls",
		},
		&cli.BoolFlag{
			Name:    "openldap-use-rfc2307bis",
			Value:   true,
			EnvVars: []string{"OPENLDAP_USE_RFC2307BIS"},
			Usage:   "openldap use RFC2307BIS schema",
		},
	}

	ldapManagerFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "groups-ou",
			Value:   "groups",
			EnvVars: []string{"GROUPS_OU"},
			Usage:   "group organizational unit",
		},
		&cli.StringFlag{
			Name:    "users-ou",
			Value:   "users",
			EnvVars: []string{"USERS_OU"},
			Usage:   "user organizational unit",
		},
		&cli.StringFlag{
			Name:    "groups-dn",
			Value:   "", // default is ou=GROUPS_OU,BASE_DN
			EnvVars: []string{"GROUPS_DN"},
			Usage:   "groups DN (default is ou=$GROUPS_OU,$BASE_DN)",
		},
		&cli.StringFlag{
			Name:    "users-dn",
			Value:   "", // default is ou=USERS_DN,BASE_DN
			EnvVars: []string{"USERS_DN"},
			Usage:   "users DN (default is ou=$USERS_DN,$BASE_DN)",
		},
		&cli.GenericFlag{
			Name: "group-membership-attribute",
			Value: &values.EnumValue{
				Enum:    []string{"uniqueMember", "memberUID"},
				Default: "uniqueMember",
			},
			EnvVars: []string{"GROUP_MEMBERSHIP_ATTRIBUTE"},
			Usage:   "group membership attribute (e.g. uniqueMember)",
		},
		&cli.BoolFlag{
			Name:    "group-membership-uses-uid",
			Value:   false,
			EnvVars: []string{"GROUP_MEMBERSHIP_USES_UID"},
			Usage:   "group membership uses UID only instead of full DN",
		},
		&cli.StringFlag{
			Name:    "account-attribute",
			Value:   "uid",
			EnvVars: []string{"ACCOUNT_ATTRIBUTE"},
			Usage:   "account attribute",
		},
		&cli.StringFlag{
			Name:    "group-attribute",
			Value:   "gid",
			EnvVars: []string{"GROUP_ATTRIBUTE"},
			Usage:   "group attribute",
		},
		&cli.StringFlag{
			Name:    "default-user-group",
			Value:   "users",
			EnvVars: []string{"DEFAULT_USER_GROUP"},
			Usage:   "default user group",
		},
		&cli.StringFlag{
			Name:    "default-admin-group",
			Value:   "admins",
			EnvVars: []string{"DEFAULT_ADMIN_GROUP"},
			Usage:   "default admin group",
		},
		&cli.StringFlag{
			Name:    "default-login-shell",
			Value:   "/bin/bash",
			EnvVars: []string{"DEFAULT_LOGIN_SHELL"},
			Usage:   "default login shell",
		},
		&cli.StringFlag{
			Name:    "default-admin-username",
			Value:   "admin",
			EnvVars: []string{"DEFAULT_ADMIN_USERNAME"},
			Usage:   "default admin username",
		},
		&cli.StringFlag{
			Name:    "default-admin-password",
			Value:   "admin",
			EnvVars: []string{"DEFAULT_ADMIN_PASSWORD"},
			Usage:   "default admin password",
		},
		&cli.BoolFlag{
			Name:    "force-create-admin",
			Value:   false,
			EnvVars: []string{"FORCE_CREATE_ADMIN"},
			Usage:   "force creation of the admin user even if there is a different user in the admin group",
		},
	}

	name := "ldap manager service"
	log.Infof("%s v%s", name, versioning.BinaryVersion(ldapmanager.Version, ldapbase.Rev))

	app := &cli.App{
		Name:    name,
		Version: versioning.BinaryVersion(ldapmanager.Version, ldapbase.Rev),
		Usage:   "manages ldap user accounts",
		Flags:   append(ldapConfigFlags, ldapManagerFlags...),
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "serve ldap manager service",
				Flags: append(serverFlags, jwtAuthFlags...),
				Action: func(cliCtx *cli.Context) error {
					grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cliCtx.Int("grpc-port")))
					if err != nil {
						return fmt.Errorf("failed to listen: %v", err)
					}
					httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cliCtx.Int("http-port")))
					if err != nil {
						return fmt.Errorf("failed to listen: %v", err)
					}

					base := ldapbase.NewLDAPManagerServer(cliCtx)
					var wg sync.WaitGroup
					wg.Add(2)
					ctx := context.Background()

					grpcServer = ldapgrpc.NewGRPCLDAPManagerServer(base, grpcListener)
					go grpcServer.Serve(ctx, &wg)

					upstream, err := grpc.DialContext(ctx, grpcListener.Addr().String(), grpc.WithInsecure(), grpc.WithBlock())
					if err != nil {
						grpcServer.Shutdown()
						return err
					}

					httpServer = ldaphttp.NewHTTPLDAPManagerServer(base, httpListener, upstream)
					go httpServer.Serve(ctx, &wg)
					wg.Wait()
					return nil
				},
			},
			// TODO: Implement CLI interface with more commands
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
