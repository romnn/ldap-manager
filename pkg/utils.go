package pkg

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/k0kubun/pp/v3"
	"github.com/go-ldap/ldap/v3"
)

const (
	// MinUID for POSIX accounts
	MinUID = 2000
	// MinGID for POSIX accounts
	MinGID = 2000

	// SortAscending ...
	SortAscending = "asc"
	// SortDescending ...
	SortDescending = "desc"

	hex = "0123456789abcdef"
)

// var (
// 	validAttributes = []string{
// 		"uid",
// 		"cn",
// 		"uidNumber",
// 		"gidNumber",
// 		"mail",
// 		"sn",
// 		"givenName",
// 		"displayName",
// 		"loginShell",
// 		"homeDirectory",
// 	}
// )

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

func IsValidAttribute(attr string) bool {
  return true
	// return Contains(validAttributes, attr)
}

func EscapeFilter(s string) string {
	return ldap.EscapeFilter(s)
}

func MustEscapeFilter(c byte) bool {
	if c > 0x7f {
		return true
	}
	switch c {
	case ',', ';', '(', ')', '\\', '*', '"', '#', '=', '+', '<', '>', 0:
		return true
	}
	return false
}

func EscapeDN(dn string) string {
	// escapes https://ldapwiki.com/wiki/DN%20Escape%20Values
	escape := 0
	for i := 0; i < len(dn); i++ {
		if MustEscapeFilter(dn[i]) {
			escape++
		}
	}
	if escape == 0 {
		return dn
	}
	buf := make([]byte, len(dn)+escape*2)
	for i, j := 0, 0; i < len(dn); i++ {
		c := dn[i]
		if MustEscapeFilter(c) {
			buf[j+0] = '\\'
			buf[j+1] = hex[c>>4]
			buf[j+2] = hex[c&0xf]
			j += 3
		} else {
			buf[j] = c
			j++
		}
	}
	return string(buf)
}

func ParseFilter(filters []string) string {
	var filter string
	for _, f := range filters {
		if pair := strings.Split(f, "="); len(pair) == 2 {
			if attr := strings.ToLower(pair[0]); IsValidAttribute(attr) {
				filter += fmt.Sprintf("(%s=*%s*)", attr, EscapeFilter(pair[1]))
			}
		}
	}
	return filter
}

func ExtractAttribute(dn string, attribute string) (string, error) {
	re, err := regexp.Compile(fmt.Sprintf("%s=(?P<Attribute>.*?),", attribute))
	if err != nil {
		return "", err
	}
	if matches := re.FindStringSubmatch(dn); len(matches) > 1 {
		if match := matches[1]; match != "" {
			return match, nil
		}
	}
	return "", fmt.Errorf("could not find attribute %q in %q", attribute, dn)
}
