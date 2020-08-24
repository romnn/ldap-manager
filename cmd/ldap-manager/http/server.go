package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2/log"
	ldapbase "github.com/romnnn/ldap-manager/cmd/ldap-manager/base"
	"github.com/urfave/cli/v2"
)

// LDAPManagerServer ...
type LDAPManagerServer struct {
	*ldapbase.LDAPManagerServer
}

// NewHTTPLDAPManagerServer ...
func NewHTTPLDAPManagerServer(base *ldapbase.LDAPManagerServer) *LDAPManagerServer {
	return &LDAPManagerServer{
		LDAPManagerServer: base,
	}
}

// Shutdown ...
func (s *LDAPManagerServer) Shutdown() {
	if s.LDAPManagerServer != nil {
		s.LDAPManagerServer.Shutdown()
	}
}

func (s *LDAPManagerServer) setupRouter() *echo.Echo {
	e := echo.New()

	// Authentication
	e.POST("/api/login", s.loginHandler)
	e.POST("/api/logout", s.logoutHandler)

	// Account management (admin only)
	e.GET("/api/accounts", s.listAccountsHandler)
	e.PUT("/api/accounts", s.newAccountHandler)

	// Group management (admin only)
	e.GET("/api/groups", s.listGroupsHandler)
	e.DELETE("/api/group/:group", s.deleteGroupHandler)
	e.PUT("/api/groups", s.newGroupHandler)
	e.POST("/api/group/:group/add", s.addGroupMemberHandler)
	e.POST("/api/group/:group/remove", s.removeGroupMemberHandler)
	e.POST("/api/group/:group/rename", s.renameGroupHandler)
	e.GET("/api/group/:group", s.getGroupHandler)

	// Edit personal account
	e.GET("/api/account/:username", s.getAccountHandler)
	e.DELETE("/api/account/:username", s.deleteAccountHandler)
	e.PUT("/api/account/:username", s.updateAccountHandler)
	e.PUT("/api/account/:username/password", s.updatePasswordHandler)

	e.Static("/", "./frontend/dist")
	return e
}

// Serve ...
func (s *LDAPManagerServer) Serve(ctx *cli.Context) error {
	if err := s.Service.BootstrapHTTP(ctx, s.setupRouter(), nil); err != nil {
		return err
	}
	go s.Connect(ctx)
	if err := s.Service.HTTPServer.Serve(s.Listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	log.Info("closing socket")
	s.Listener.Close()
	return nil
}
