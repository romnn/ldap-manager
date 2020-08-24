package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2/log"
	ldapmanager "github.com/romnnn/ldap-manager"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

func (s *LDAPManagerServer) listAccountsHandler(c echo.Context) error {
	var options pb.ListOptions
	if err := c.Bind(&options); err != nil {
		log.Error(err)
		return err
	}
	users, err := s.Manager.GetUserList(&pb.GetUserListRequest{
		Options: &options,
	})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list accounts")
	}
	return c.JSONPretty(http.StatusOK, users, "  ")
}

func (s *LDAPManagerServer) getAccountHandler(c echo.Context) error {
	username := c.Param("username")
	user, err := s.Manager.GetAccount(&pb.GetAccountRequest{Username: username})
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
	var req pb.NewAccountRequest
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
	leaveGroups := false
	if err := s.Manager.DeleteAccount(&pb.DeleteAccountRequest{Username: username}, leaveGroups); err != nil {
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

func (s *LDAPManagerServer) newAccount(c echo.Context, req *pb.NewAccountRequest) error {
	if err := s.Manager.NewAccount(req); err != nil {
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
	var req pb.NewAccountRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	return s.newAccount(c, &req)
}
