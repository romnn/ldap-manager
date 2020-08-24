package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2/log"
	ldapmanager "github.com/romnnn/ldap-manager"
)

type loginRequest struct {
	UserID   string `json:"user_id" form:"user_id"`
	Password string `json:"password" form:"password"`
}

type loginResponse struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func (s *LDAPManagerServer) loginHandler(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	if req.UserID == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please provide valid credentials")
	}
	userDN, err := s.Manager.AuthenticateUser(&ldapmanager.AuthenticateUserRequest{Username: req.UserID, Password: req.Password})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "no such user")
	}
	isMember, err := s.Manager.IsGroupMember(&ldapmanager.IsGroupMemberRequest{Username: req.UserID, Group: s.Manager.DefaultAdminGroup})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to check admin status of user")
	}
	log.Infof("user %q (%s) logged in (admin=%t)", req.UserID, userDN, isMember)
	// TODO: set_passkey_cookie($user_auth,$is_admin);
	u := &loginResponse{
		Name:  "Jon",
		Email: "jon@labstack.com",
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}

func (s *LDAPManagerServer) logoutHandler(c echo.Context) error {
	if s.Service.Healthy {
		c.String(http.StatusOK, "ok")
	} else {
		c.String(http.StatusServiceUnavailable, "service is not available")
	}
	return nil
}
