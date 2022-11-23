package cli

import (
	"testing"
	"time"
)

func TestDurationValueDefault(t *testing.T) {
	defaultDur := 24 * time.Hour
	value := DurationValue{
		Default: defaultDur,
	}

	dur, err := time.ParseDuration(value.String())
	if err != nil {
		t.Fatalf("failed to parse %q to duration: %v", value.String(), err)
	}
	if dur != defaultDur {
		t.Fatalf("expected duration %#v but got %#v", defaultDur, dur)
	}
}

func TestDurationValue(t *testing.T) {
	defaultDur := 24 * time.Hour
	for _, c := range []struct {
		input    string
		expected time.Duration
	}{
		{
			input:    "1h20m",
			expected: 1*time.Hour + 20*time.Minute,
		},
		{
			input:    "1h 20m",
			expected: 1*time.Hour + 20*time.Minute,
		},
		{
			input:    "1h 20m3s",
			expected: 1*time.Hour + 20*time.Minute + 3*time.Second,
		},
		{
			input:    "20m1h",
			expected: 1*time.Hour + 20*time.Minute,
		},
		{
			input:    "0s",
			expected: 0 * time.Second,
		},
		{
			input:    "5m",
			expected: 5 * time.Minute,
		},
	} {
		value := DurationValue{
			Default: defaultDur,
		}
		if err := value.Set(c.input); err != nil {
			t.Errorf("failed to set duration to %q: %v", c.input, err)
		}
		dur, err := time.ParseDuration(value.String())
		if err != nil {
			t.Errorf("failed to parse %q to duration: %v", value.String(), err)
		}
		if dur != c.expected {
			t.Errorf("expected duration %v but got %v", c.expected, dur)
		}
	}
}
