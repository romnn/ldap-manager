package cli

import (
	"fmt"
	"strings"
	"time"
)

// DurationValue is a CLI value that represents a duration
type DurationValue struct {
	duration time.Duration
	Default  time.Duration
	set      bool
}

// Set sets the duration value and returns an error if the value is not valid
func (dur *DurationValue) Set(value string) error {
	value = strings.ToLower(value)
	value = strings.ReplaceAll(value, " ", "")
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("value %q is not a valid duration (e.g. 5h30m40s): %v", value, err)
	}
	dur.set = true
	dur.duration = duration
	return nil
}

// String returns the current duration value or the default duration otherwise
func (dur *DurationValue) String() string {
	if !dur.set {
		return fmt.Sprintf("%v", dur.Default)
	}
	return fmt.Sprintf("%v", dur.duration)
}
