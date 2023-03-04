package cli

import (
	"fmt"
	"strings"
)

// EnumValue is a CLI value that can take an enum value
type EnumValue struct {
	Enum      []string
	Default   string
	AllowNone bool
	set       bool
	selected  string
}

// Set sets the enum value and returns an error if the value is not valid
func (enum *EnumValue) Set(value string) error {
	value = strings.TrimSpace(strings.ToLower(value))
	for _, e := range enum.Enum {
		if strings.ToLower(e) == value {
			enum.selected = e
			enum.set = true
			return nil
		}
	}
	if !enum.AllowNone {
		return fmt.Errorf(
			"unknown option %q, must be one of %v",
			value, enum.Enum,
		)
	}
	return nil
}

// String returns the current enum value or the default value otherwise
func (enum *EnumValue) String() string {
	if !enum.set {
		return enum.Default
	}
	return enum.selected
}
