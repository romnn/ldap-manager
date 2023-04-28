package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-ldap/ldap/v3"
	"github.com/jwalton/go-supportscolor"
	pretty "github.com/k0kubun/pp/v3"
)

var (
	// pp for pretty printing
	pp = func() *pretty.PrettyPrinter {
		useColor := supportscolor.Stdout().SupportsColor
		pp := pretty.New()
		pp.SetColoringEnabled(useColor)
		pp.SetExportedOnly(true)
		return pp
	}()
)

const (
	// MinUID for POSIX accounts
	MinUID = 2000
	// MinGID for POSIX accounts, reserved for the users group
	MinGID = 2000
)

// ParseDN parses a DN into its parts.
func ParseDN(dn string) map[string][]string {
	parsed := make(map[string][]string)
	re := regexp.MustCompile("([^,]+)=([^,]+)")
	parts := re.FindAllStringSubmatch(dn, -1)
	for _, part := range parts {
		_, present := parsed[part[1]]
		if !present {
			parsed[part[1]] = []string{}
		}
		parsed[part[1]] = append(parsed[part[1]], part[2])
	}
	return parsed
}

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

// Dedup is a generic function that removes duplicates in a slice.
func Dedup[T comparable](list []T) []T {
	allKeys := make(map[T]bool)
	for _, item := range list {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
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
