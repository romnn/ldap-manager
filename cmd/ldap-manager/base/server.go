package base

import (
	"fmt"
	"net"

	"github.com/neko-neko/echo-logrus/v2/log"
	gogrpcservice "github.com/romnnn/go-grpc-service"
	ldapmanager "github.com/romnnn/ldap-manager"
	ldapconfig "github.com/romnnn/ldap-manager/config"
	"github.com/urfave/cli/v2"
)

// Rev is set on build time to the git HEAD
var Rev = ""

// Version is incremented using bump2version
const Version = "0.0.1"

// LDAPManagerServer ...
type LDAPManagerServer struct {
	gogrpcservice.Service
	Manager  *ldapmanager.LDAPManager
	Listener net.Listener
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	s.Service.GracefulStop()
	if s.Manager != nil {
		s.Manager.Close()
	}
}

// NewLDAPManagerServer ...
func NewLDAPManagerServer(ctx *cli.Context) (*LDAPManagerServer, error) {
	var err error
	s := &LDAPManagerServer{
		Service: gogrpcservice.Service{
			Name:               "ldap manager service",
			Version:            Version,
			BuildTime:          Rev,
			HTTPHealthCheckURL: "health/healthz",
		},
		Manager: &ldapmanager.LDAPManager{
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
	s.Listener, err = net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}
	return s, nil
}

// Setup prepares the service
func (s *LDAPManagerServer) Setup(ctx *cli.Context) error {
	if err := s.Manager.Setup(); err != nil {
		return err
	}
	return nil
}

// Connect starts the service
func (s *LDAPManagerServer) Connect(ctx *cli.Context) {
	log.Info("connecting...")
	if err := s.Setup(ctx); err != nil {
		log.Error(err)
		s.Shutdown()
		return
	}
	s.Service.Ready = true
	s.Service.SetHealthy(true)
	log.Infof("%s ready at %s", s.Service.Name, s.Listener.Addr())
}
