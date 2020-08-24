package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ldapmanager "github.com/romnnn/ldap-manager"
	log "github.com/sirupsen/logrus"
)

type groupMemberRequest struct {
	Username string `json:"username" xml:"username"`
}

func (s *LDAPManagerServer) addGroupMemberHandler(c echo.Context) error {
	group := c.Param("group")
	var req groupMemberRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	if err := s.manager.AddGroupMember(&ldapmanager.AddGroupMemberRequest{Group: group, Username: req.Username}); err != nil {
		switch err.(type) {
		case *ldapmanager.GroupValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.ZeroOrMultipleGroupsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleGroupsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to add group member")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}

func (s *LDAPManagerServer) removeGroupMemberHandler(c echo.Context) error {
	group := c.Param("group")
	var req groupMemberRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	if err := s.manager.DeleteGroupMember(&ldapmanager.DeleteGroupMemberRequest{Group: group, Username: req.Username}); err != nil {
		switch err.(type) {
		case *ldapmanager.GroupValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.ZeroOrMultipleGroupsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleGroupsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete group member")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}
