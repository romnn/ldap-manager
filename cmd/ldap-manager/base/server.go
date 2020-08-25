package base

import (
	"net"

	gogrpcservice "github.com/romnnn/go-grpc-service"
	ldapmanager "github.com/romnnn/ldap-manager"
	ldapconfig "github.com/romnnn/ldap-manager/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Rev is set on build time to the git HEAD
var Rev = ""

// LDAPManagerServer ...
type LDAPManagerServer struct {
	gogrpcservice.Service
	Manager *ldapmanager.LDAPManager
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	s.Service.GracefulStop()
	if s.Manager != nil {
		s.Manager.Close()
	}
}

// NewLDAPManagerServer ...
func NewLDAPManagerServer(ctx *cli.Context) *LDAPManagerServer {
	return &LDAPManagerServer{
		Service: gogrpcservice.Service{
			Name:               "ldap manager service",
			Version:            ldapmanager.Version,
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
}

// Setup prepares the service
func (s *LDAPManagerServer) Setup(ctx *cli.Context) error {
	if err := s.Manager.Setup(); err != nil {
		return err
	}
	return nil
}

// Connect starts the service
func (s *LDAPManagerServer) Connect(ctx *cli.Context, listener net.Listener) {
	log.Info("connecting...")
	if err := s.Setup(ctx); err != nil {
		log.Error(err)
		s.Shutdown()
		return
	}
	s.Service.Ready = true
	s.Service.SetHealthy(true)
	log.Infof("%s ready at %s", s.Service.Name, listener.Addr())
}
