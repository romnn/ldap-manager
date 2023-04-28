package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/romnn/go-service/pkg/auth"
	flags "github.com/romnn/ldap-manager/cmd/ldap-manager/cli"
	ldapgrpc "github.com/romnn/ldap-manager/cmd/ldap-manager/grpc"
	ldaphttp "github.com/romnn/ldap-manager/cmd/ldap-manager/http"

	ldapmanager "github.com/romnn/ldap-manager/pkg"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

// Version is set during build
var Version = ""

// Rev is set during build
var Rev = ""

// newLDAPManager configures the LDAP manager based on the CLI config
func newLDAPManager(ctx *cli.Context) ldapmanager.LDAPManager {
	hasReadOnlyUser := ctx.String(flags.LdapReadOnlyUser.Name) != ""
	baseDN := ctx.String(flags.LdapBaseDn.Name)
	groupsOU := ctx.String(flags.GroupsOu.Name)
	usersOU := ctx.String(flags.UsersOu.Name)

	groupsDN := ctx.String(flags.GroupsDn.Name)
	if groupsDN == "" {
		groupsDN = fmt.Sprintf(
			"ou=%s,%s",
			groupsOU, baseDN,
		)
	}
	userGroupDN := ctx.String(flags.UsersDn.Name)
	if userGroupDN == "" {
		userGroupDN = fmt.Sprintf(
			"ou=%s,%s",
			usersOU, baseDN,
		)
	}

	config := ldapconfig.Config{
		Host:                ctx.String(flags.LdapHost.Name),
		Port:                ctx.Int(flags.LdapPort.Name),
		Protocol:            ctx.String(flags.LdapProtocol.Name),
		Organization:        ctx.String(flags.LdapOrganization.Name),
		Domain:              ctx.String(flags.LdapDomain.Name),
		BaseDN:              baseDN,
		AdminUsername:       ctx.String(flags.LdapAdminUsername.Name),
		AdminPassword:       ctx.String(flags.LdapAdminPassword.Name),
		ReadOnlyUser:        hasReadOnlyUser,
		ReadOnlyUsername:    ctx.String(flags.LdapReadOnlyUser.Name),
		ReadOnlyPassword:    ctx.String(flags.LdapReadOnlyPassword.Name),
		TLS:                 ctx.Bool(flags.LdapTLS.Name),
		UseRFC2307BISSchema: ctx.Bool(flags.LdapUseRfc2307Bis.Name),
	}

	return ldapmanager.LDAPManager{
		Config:                   config,
		GroupsOU:                 groupsOU,
		UsersOU:                  usersOU,
		GroupsDN:                 groupsDN,
		UserGroupDN:              userGroupDN,
		GroupMembershipAttribute: ctx.String(flags.GroupMembershipAttribute.Name),
		GroupMembershipUsesUID:   ctx.Bool(flags.GroupMembershipUsesUID.Name),
		AccountAttribute:         ctx.String(flags.AccountAttribute.Name),
		GroupAttribute:           ctx.String(flags.GroupAttribute.Name),
		DefaultUserGroup:         ctx.String(flags.DefaultUserGroup.Name),
		DefaultAdminGroup:        ctx.String(flags.DefaultAdminGroup.Name),
		DefaultUserShell:         ctx.String(flags.DefaultLoginShell.Name),
		DefaultAdminUsername:     ctx.String(flags.DefaultAdminUsername.Name),
		DefaultAdminPassword:     ctx.String(flags.DefaultAdminPassword.Name),
		ForceCreateAdmin:         ctx.Bool(flags.ForceCreateAdmin.Name),
	}
}

// newAuthenticator configures the authenticator based on the CLI config
func newAuthenticator(ctx *cli.Context) auth.Authenticator {
	expiresAfter, _ := time.ParseDuration(
		ctx.String(flags.ExpirationTime.Name),
	)
	return auth.Authenticator{
		ExpiresAfter: expiresAfter,
		Issuer:       ctx.String(flags.Issuer.Name),
		Audience:     ctx.String(flags.Audience.Name),
	}
}

// newAuthenticator crates a new auth key config based on the CLI config
func newAuthKeyConfig(ctx *cli.Context) auth.KeyConfig {
	return auth.KeyConfig{
		Jwks:     ctx.String(flags.Jwks.Name),
		JwksFile: ctx.String(flags.JwksFile.Name),
		Key:      ctx.String(flags.Key.Name),
		KeyFile:  ctx.String(flags.KeyFile.Name),
		Generate: ctx.Bool(flags.Generate.Name),
	}
}

func versionString(version string, buildTime string) string {
	if buildTime != "" {
		return fmt.Sprintf(
			"%s (built on %s)",
			version, buildTime,
		)
	}
	return version
}

func setupLogging(cliCtx *cli.Context) {
	level := log.InfoLevel
	switch strings.ToLower(cliCtx.String(flags.LogLevel.Name)) {
	case "debug":
		level = log.DebugLevel
		break
	case "info":
		level = log.InfoLevel
		break
	case "warn":
		level = log.WarnLevel
		break
	case "fatal":
		level = log.FatalLevel
		break
	case "trace":
		level = log.TraceLevel
		break
	case "error":
		level = log.ErrorLevel
		break
	case "panic":
		level = log.PanicLevel
		break
	default:
		break
	}
	log.SetLevel(level)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:  cliCtx.Bool(flags.ForceColors.Name),
		DisableQuote: cliCtx.Bool(flags.DisableQuote.Name),
	})
}

func serve(cliCtx *cli.Context) error {
	setupLogging(cliCtx)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	grpcAddr := fmt.Sprintf(":%d", cliCtx.Int(flags.GRPCPort.Name))
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	httpAddr := fmt.Sprintf(":%d", cliCtx.Int(flags.HTTPPort.Name))
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	keyConfig := newAuthKeyConfig(cliCtx)
	authenticator := newAuthenticator(cliCtx)
	if err := authenticator.SetupKeys(&keyConfig); err != nil {
		return err
	}

	manager := newLDAPManager(cliCtx)
	log.Debug(ldapmanager.PrettyPrint(manager))
	log.Debug(ldapmanager.PrettyPrint(manager.Config))

	if err := manager.Setup(); err != nil {
		log.Error(err)
		return err
	}

	serveErrChan := make(chan error, 2)
	ctx, cancel := context.WithDeadline(
		context.Background(),
		time.Now().Add(1*time.Minute),
	)
	grpcService := ldapgrpc.NewLDAPManagerService(
		ctx, manager, authenticator,
	)

	// if setup succeeds, we are serving
	grpcService.SetHealthy(true)

	go func() {
		log.Infof("grpc listening on: %v", grpcListener.Addr())
		serveErrChan <- grpcService.Serve(grpcListener)
		log.Infof("shutdown grpc service")
	}()

	var httpService *ldaphttp.LDAPManagerService

	go func() {
		upstream, err := grpc.DialContext(
			ctx,
			grpcListener.Addr().String(),
			grpc.WithInsecure(),
		)
		if err != nil {
			serveErrChan <- fmt.Errorf(
				"failed to dial grpc upstream: %v",
				err,
			)
			return
		}
		defer upstream.Close()

		httpService, err = ldaphttp.NewLDAPManagerService(
			ctx,
			upstream,
			&ldaphttp.Config{
				ServeStatic: !cliCtx.Bool(flags.NoStatic.Name),
				StaticPath:  cliCtx.String(flags.StaticRoot.Name),
			},
		)
		if err != nil {
			serveErrChan <- fmt.Errorf(
				"failed to start http service: %v",
				err,
			)
			return
		}

		httpService.SetHealthy(true)
		log.Infof("http listening on: %v", httpListener.Addr())
		serveErrChan <- httpService.Serve(httpListener)
		log.Infof("shutdown http service")
	}()

	var once sync.Once

	// shutdown should only be called once
	shutdown := func() {
		log.Warnf("shutdown ...")
		// cancel setup
		cancel()
		// shutdown http service first,
		// as it keeps a connection to the GRPC service
		if httpService != nil {
			log.Warnf("shutting down http ...")
			httpService.Shutdown()
		}
		log.Warnf("shutting down grpc ...")
		grpcService.Shutdown()
	}

	go func() {
		<-shutdownChan
		once.Do(shutdown)
	}()

	// await completion or error
	var serveErr error
	for i := 0; i < 2; i++ {
		err := <-serveErrChan
		grpcService.SetHealthy(false)
		// todo: shutdown anyways?
		if err != nil && errors.Is(err, context.Canceled) {
			serveErr = err
			go once.Do(shutdown)
		}
	}
	return serveErr
}

func main() {
	version := versionString(Version, Rev)
	log.Infof("LDAPManager v%s", version)

	app := &cli.App{
		Name:    "LDAPManager",
		Usage:   "service for managing LDAP",
		Version: version,
		Flags: append(
			flags.LdapConfigFlags,
			flags.LdapFlags...,
		),
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "serves the LDAPManager API",
				Flags: append(
					flags.ServiceFlags,
					flags.AuthFlags...,
				),
				Action: serve,
			},
			// TODO: Implement CLI interface with more commands
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
