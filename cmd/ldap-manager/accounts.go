package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ldapmanager "github.com/romnnn/ldap-manager"
	log "github.com/sirupsen/logrus"
)

func (s *LDAPManagerServer) listAccountsHandler(c echo.Context) error {
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

func (s *LDAPManagerServer) getAccountHandler(c echo.Context) error {
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

func (s *LDAPManagerServer) updateAccountHandler(c echo.Context) error {
	var req ldapmanager.NewAccountRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	// Delete existing account
	if err := s.deleteAccount(c, c.Param("username")); err != nil {
		return err
	}
	// Insert the updated account
	return s.newAccount(c, &req)
}

func (s *LDAPManagerServer) deleteAccount(c echo.Context, username string) error {
	if err := s.manager.DeleteAccount(&ldapmanager.DeleteAccountRequest{Username: username}); err != nil {
		switch err.(type) {
		case *ldapmanager.ZeroOrMultipleAccountsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleAccountsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete account")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}

func (s *LDAPManagerServer) deleteAccountHandler(c echo.Context) error {
	return s.deleteAccount(c, c.Param("username"))
}

func (s *LDAPManagerServer) newAccount(c echo.Context, req *ldapmanager.NewAccountRequest) error {
	if err := s.manager.NewAccount(req); err != nil {
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

func (s *LDAPManagerServer) newAccountHandler(c echo.Context) error {
	var req ldapmanager.NewAccountRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	return s.newAccount(c, &req)
}
