package main

import (
	"net/http"
	// "errors"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type loginRequest struct {
	UserID   string `json:"user_id" form:"user_id"`
	Password string `json:"password" form:"password"`
}

type loginResponse struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

func (s *LDAPManagerServer) login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	if req.UserID == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "please provide valid credentials")
	}
	userDN, err := s.AuthenticateUser(req.UserID, req.Password)
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusNotFound, "no such user")
	}
	isMember, err := s.IsGroupMember(s.DefaultAdminGroup, req.UserID)
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

func (s *LDAPManagerServer) logout(c echo.Context) error {
	if s.Service.Healthy {
		c.String(http.StatusOK, "ok")
	} else {
		c.String(http.StatusServiceUnavailable, "service is not available")
	}
	return nil
}
