package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ldapmanager "github.com/romnnn/ldap-manager"
	log "github.com/sirupsen/logrus"
)

func (s *LDAPManagerServer) updatePasswordHandler(c echo.Context) error {
	var req ldapmanager.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}

	if err := s.manager.ChangePassword(&req); err != nil {
		switch err.(type) {
		case *ldapmanager.ZeroOrMultipleAccountsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleAccountsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update password")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}
