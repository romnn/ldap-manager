package ldapmanager

import (
	log "github.com/sirupsen/logrus"
)

func (m *LDAPManager) examples() {
	/*
		if err := s.BindReadOnly(); err != nil {
			return err
		}
	*/
	/* Test checking group members
	isMember, err := s.IsGroupMember("billy", "admins")
	if err != nil {
		return err
	}
	log.Info(isMember)
	*/

	// Add admin to admin group
	/*
		members, err := s.AddGroupMember()
		if err != nil {
			return err
		}
		log.Info(members)

		members, err := s.GetGroupMembers(s.DefaultAdminGroup, 0, 0, "")
		if err != nil {
			return err
		}
		log.Info(members)
	*/

	// Add a sample user
	if err := m.NewAccount(&NewAccountRequest{
		Username: "romnn",
		Password: "Hallo Welt",
	}); err != nil {
		log.Error(err)
	}
}
