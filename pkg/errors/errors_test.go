package errors

import (
	"testing"
)

var (
	// Accounts
	sampleAccountValidationError      = &AccountValidationError{}
	sampleZeroOrMultipleAccountsError = &ZeroOrMultipleAccountsError{}
	sampleAccountAlreadyExistsError   = &AccountAlreadyExistsError{}

	// Groups
	sampleValidationError           = &ValidationError{}
	sampleZeroOrMultipleGroupsError = &ZeroOrMultipleGroupsError{}
	sampleGroupAlreadyExistsError   = &GroupAlreadyExistsError{}

	// Group members
	sampleNoSuchMemberError          = &NoSuchMemberError{}
	sampleRemoveLastGroupMemberError = &RemoveLastGroupMemberError{}
	sampleMemberAlreadyExistsError   = &MemberAlreadyExistsError{}
)

func toInterface(in interface{}) interface{} {
	return in
}

// Accounts

func TestAccountValidationError(t *testing.T) {
	_, ok := toInterface(sampleAccountValidationError).(Error)
	if !ok {
		t.Errorf("expected AccountValidationError to implement Error interface")
	}
}

func TestZeroOrMultipleAccountsError(t *testing.T) {
	_, ok := toInterface(sampleZeroOrMultipleAccountsError).(Error)
	if !ok {
		t.Errorf("expected ZeroOrMultipleAccountsError to implement Error interface")
	}
}

func TestAccountAlreadyExistsError(t *testing.T) {
	_, ok := toInterface(sampleAccountAlreadyExistsError).(Error)
	if !ok {
		t.Errorf("expected AccountAlreadyExistsError to implement Error interface")
	}
}

// Groups

func TestValidationError(t *testing.T) {
	_, ok := toInterface(sampleValidationError).(Error)
	if !ok {
		t.Errorf("expected ValidationError to implement Error interface")
	}
}

func TestZeroOrMultipleGroupsError(t *testing.T) {
	_, ok := toInterface(sampleZeroOrMultipleGroupsError).(Error)
	if !ok {
		t.Errorf("expected ZeroOrMultipleGroupsError to implement Error interface")
	}
}

func TestGroupAlreadyExistsError(t *testing.T) {
	_, ok := toInterface(sampleGroupAlreadyExistsError).(Error)
	if !ok {
		t.Errorf("expected GroupAlreadyExistsError to implement Error interface")
	}
}

// Groups members

func TestNoSuchMemberError(t *testing.T) {
	_, ok := toInterface(sampleNoSuchMemberError).(Error)
	if !ok {
		t.Errorf("expected NoSuchMemberError to implement Error interface")
	}
}

func TestRemoveLastGroupMemberError(t *testing.T) {
	_, ok := toInterface(sampleRemoveLastGroupMemberError).(Error)
	if !ok {
		t.Errorf("expected RemoveLastGroupMemberError to implement Error interface")
	}
}

func TestMemberAlreadyExistsError(t *testing.T) {
	_, ok := toInterface(sampleMemberAlreadyExistsError).(Error)
	if !ok {
		t.Errorf("expected MemberAlreadyExistsError to implement Error interface")
	}
}
