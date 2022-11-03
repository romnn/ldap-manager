package pkg

import (
	"fmt"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/k0kubun/pp/v3"
)

const (
	// MinUID for POSIX accounts
	MinUID = 2000
	// MinGID for POSIX accounts
	MinGID = 2000
)

// GroupDN returns the full group DN for a group name
func (m *LDAPManager) GroupDN(name string) string {
	return fmt.Sprintf(
		"cn=%s,%s",
		EscapeDN(name),
		m.GroupsDN,
	)
}

// UserDN returns the full user DN for a user name
func (m *LDAPManager) UserDN(name string) string {
	return fmt.Sprintf(
		"%s=%s,%s",
		m.AccountAttribute,
		EscapeDN(name),
		m.UserGroupDN,
	)
}

// PrettyPrint formats an interface into a human readable string
func PrettyPrint(m interface{}) string {
	return pp.Sprint(m)
}

// Contains is a generic function that checks if a collection contains a value
func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

// EscapeFilter escapes an LDAP filter to avoid LDAP injection attacks
func EscapeFilter(s string) string {
	return ldap.EscapeFilter(s)
}

// EscapeDN escapes an LDAP DN to avoid LDAP injection attacks
//
// source: https://github.com/go-ldap/ldap/blob/master/ldap.go
// Note: for the next ldap release, we can directly use `ldap.EscapeDN`
func EscapeDN(dn string) string {
	if dn == "" {
		return ""
	}

	builder := strings.Builder{}

	for i, r := range dn {
		// Escape leading and trailing spaces
		if (i == 0 || i == len(dn)-1) && r == ' ' {
			builder.WriteRune('\\')
			builder.WriteRune(r)
			continue
		}

		// Escape leading '#'
		if i == 0 && r == '#' {
			builder.WriteRune('\\')
			builder.WriteRune(r)
			continue
		}

		// Escape characters as defined in RFC4514
		switch r {
		case '"', '+', ',', ';', '<', '>', '\\':
			builder.WriteRune('\\')
			builder.WriteRune(r)
		case '\x00': // Null byte may not be escaped by a leading backslash
			builder.WriteString("\\00")
		default:
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// BuildFilter escapes and concatenates multiple filter expressions
func BuildFilter(filters []string) string {
	var filter string
	for _, f := range filters {
		if pair := strings.Split(f, "="); len(pair) == 2 {
			attr := strings.ToLower(pair[0])
			filter += fmt.Sprintf(
				"(%s=*%s*)",
				attr, EscapeFilter(pair[1]),
			)
		}
	}
	return filter
}
