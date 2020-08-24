package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	ldapmanager "github.com/romnnn/ldap-manager"
	log "github.com/sirupsen/logrus"
)

func (s *LDAPManagerServer) listGroupsHandler(c echo.Context) error {
	var options ldapmanager.ListOptions
	if err := c.Bind(&options); err != nil {
		log.Error(err)
		return err
	}
	groups, err := s.manager.GetGroupList(&ldapmanager.GetGroupListRequest{
		ListOptions: options,
	})
	if err != nil {
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to list groups")
	}
	return c.JSONPretty(http.StatusOK, groups, "  ")
}

type renameGroupRequest struct {
	NewName string `json:"name" xml:"name"`
}

func (s *LDAPManagerServer) renameGroupHandler(c echo.Context) error {
	groupName := c.Param("group")
	var req renameGroupRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	if err := s.manager.RenameGroup(&ldapmanager.RenameGroupRequest{Group: groupName, NewName: req.NewName}); err != nil {
		switch err.(type) {
		case *ldapmanager.GroupValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.ZeroOrMultipleGroupsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleGroupsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get group")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}

func (s *LDAPManagerServer) getGroupHandler(c echo.Context) error {
	groupName := c.Param("group")
	group, err := s.manager.GetGroup(&ldapmanager.GetGroupRequest{Group: groupName})
	if err != nil {
		switch err.(type) {
		case *ldapmanager.GroupValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.ZeroOrMultipleGroupsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleGroupsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get group")
	}
	return c.JSONPretty(http.StatusOK, group, "  ")
}

func (s *LDAPManagerServer) deleteGroupHandler(c echo.Context) error {
	group := c.Param("group")
	if err := s.manager.DeleteGroup(group); err != nil {
		switch err.(type) {
		case *ldapmanager.GroupValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.ZeroOrMultipleGroupsError:
			return echo.NewHTTPError(err.(*ldapmanager.ZeroOrMultipleGroupsError).Status(), err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to delete group")
	}
	return c.JSONPretty(http.StatusOK, nil, "  ")
}

func (s *LDAPManagerServer) newGroupHandler(c echo.Context) error {
	var req ldapmanager.NewGroupRequest
	if err := c.Bind(&req); err != nil {
		log.Error(err)
		return err
	}
	req.Strict = true // enforces all members of the group to already exist
	if err := s.manager.NewGroup(&req); err != nil {
		switch err.(type) {
		case *ldapmanager.GroupValidationError:
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		case *ldapmanager.GroupAlreadyExistsError:
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}
		log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to add new group")
	}
	return nil
}
