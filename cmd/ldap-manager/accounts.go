package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ldapmanager "github.com/romnnn/ldap-manager"
	log "github.com/sirupsen/logrus"
)

func (s *LDAPManagerServer) listAccounts(c echo.Context) error {
	var options ldapmanager.ListOptions
	if err := c.Bind(&options); err != nil {
		log.Error(err)
		return err
	}
	users, err := s.manager.GetUserList(&ldapmanager.GetUserListRequest{
		ListOptions: options,
	})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list accounts")
	}
	return c.JSONPretty(http.StatusOK, users, "  ")
}

func (s *LDAPManagerServer) getAccount(c echo.Context) error {
	username := c.Param("username")
	user, err := s.manager.GetAccount(username)
	if err != nil {
		switch err.(type) {
		case *ldapmanager.ZeroOrMultipleAccountsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleAccountsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get account")
	}
	return c.JSONPretty(http.StatusOK, user, "  ")
}

func (s *LDAPManagerServer) deleteAccount(c echo.Context) error {
	username := c.Param("username")
	if err := s.manager.DeleteAccount(username); err != nil {
		switch err.(type) {
		case *ldapmanager.ZeroOrMultipleAccountsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleAccountsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get account")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}

func (s *LDAPManagerServer) newAccount(c echo.Context) error {
	var req ldapmanager.NewAccountRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	if err := s.manager.NewAccount(&req); err != nil {
		switch err.(type) {
		case *ldapmanager.AccountValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.AccountAlreadyExistsError:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to add new account")
	}
	return nil
}
