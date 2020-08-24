package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/neko-neko/echo-logrus/v2/log"
	ldapmanager "github.com/romnnn/ldap-manager"
	pb "github.com/romnnn/ldap-manager/grpc/ldap-manager"
)

func (s *LDAPManagerServer) updatePasswordHandler(c echo.Context) error {
	var req pb.ChangePasswordRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}

	if err := s.Manager.ChangePassword(&req); err != nil {
		switch err.(type) {
		case *ldapmanager.ZeroOrMultipleAccountsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleAccountsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update password")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}
