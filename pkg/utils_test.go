package pkg

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
)

// TestParseDN tests parsing a DN into its parts
func TestParseDN(t *testing.T) {
	for _, c := range []struct {
		dn       string
		expected map[string][]string
	}{
		{
			dn: "uid=romnn,ou=users,dc=romnn,dc=com",
			expected: map[string][]string{
				"uid": []string{"romnn"},
				"ou":  []string{"users"},
				"dc":  []string{"romnn", "com"},
			},
		},
	} {
		parsed := ParseDN(c.dn)
		equal := cmp.Equal(parsed, c.expected)
		diff := cmp.Diff(ParseDN(c.dn), c.expected)
		if !equal {
			t.Log(parsed)
			t.Log(c.expected)
			t.Errorf(
				"unexpected parsed parts for dn %q: %s",
				c.dn, diff,
			)
		}
	}
}

// TestEscapeDN tests escaping LDAP DN's to avoid LDAP injection attacks
func TestEscapeDN(t *testing.T) {
	// https://github.com/tcort/ldap-escape/blob/master/test/ldap-escape.test.js
	for _, c := range []struct {
		filter      string
		values      interface{}
		escaped     string
		description string
	}{
		{
			filter: "cn={{.CN}},dc={{.DC}}",
			values: struct {
				CN string
				DC string
			}{
				CN: EscapeDN("alice"),
				DC: EscapeDN("com"),
			},
			escaped:     `cn=alice,dc=com`,
			description: "should work in the base case (no escaping)",
		},
		{
			filter: "cn={{.CN}},dc={{.DC}}",
			values: struct {
				CN string
				DC string
			}{
				CN: EscapeDN(" alice"),
				DC: EscapeDN("com"),
			},
			escaped:     `cn=\ alice,dc=com`,
			description: "should escape a leading space",
		},
		{
			filter: "cn={{.CN}},dc={{.DC}}",
			values: struct {
				CN string
				DC string
			}{
				CN: EscapeDN("#alice"),
				DC: EscapeDN("com"),
			},
			escaped:     `cn=\#alice,dc=com`,
			description: "should escape a leading hash",
		},
		{
			filter: "cn={{.CN}},dc={{.DC}}",
			values: struct {
				CN string
				DC string
			}{
				CN: EscapeDN("# "),
				DC: EscapeDN("com"),
			},
			escaped:     `cn=\#\ ,dc=com`,
			description: "should escape a leading hash and trailing space",
		},
		{
			filter: "cn={{.CN}},dc={{.DC}}",
			values: struct {
				CN string
				DC string
			}{
				CN: EscapeDN("alice "),
				DC: EscapeDN("com"),
			},
			escaped:     `cn=alice\ ,dc=com`,
			description: "should escape a trailing space",
		},
		{
			filter: "cn={{.CN}},dc={{.DC}}",
			values: struct {
				CN string
				DC string
			}{
				CN: EscapeDN("   "),
				DC: EscapeDN("com"),
			},
			escaped:     `cn=\  \ ,dc=com`,
			description: "should escape a dn of just 3 spaces",
		},
		{
			filter: "{{.DN}}",
			values: struct {
				DN string
			}{
				DN: EscapeDN(` Hello\ + , "World" ; `),
			},
			escaped:     `\ Hello\\ \+ \, \"World\" \;\ `,
			description: "should correctly escape the OWASP Christmas Tree Example",
		},
		{
			filter: "cn={{.CN}},ou=West,dc=MyDomain,dc=com",
			values: struct {
				CN string
			}{
				CN: EscapeDN(`Smith, James K.`),
			},
			escaped:     `cn=Smith\, James K.,ou=West,dc=MyDomain,dc=com`,
			description: "should correctly escape the Active Directory Example 1",
		},
		{
			filter: "ou={{.OU}},dc=MyDomain,dc=com",
			values: struct {
				OU string
			}{
				OU: EscapeDN(`Sales\Engineering`),
			},
			escaped:     `ou=Sales\\Engineering,dc=MyDomain,dc=com`,
			description: "should correctly escape the Active Directory Example 2",
		},
		{
			filter: "cn={{.CN}},ou=West,dc=MyDomain,dc=com",
			values: struct {
				CN string
			}{
				CN: EscapeDN(`East#Test + Lab`),
			},
			escaped:     `cn=East#Test \+ Lab,ou=West,dc=MyDomain,dc=com`,
			description: "should correctly escape the Active Directory Example 3",
		},
		{
			filter: "cn={{.CN}},ou=West,dc=MyDomain,dc=com",
			values: struct {
				CN string
			}{
				CN: EscapeDN(` Jim Smith `),
			},
			escaped:     `cn=\ Jim Smith\ ,ou=West,dc=MyDomain,dc=com`,
			description: "should correctly escape the Active Directory Example 4",
		},
	} {
		tmpl, err := template.New("dn").Parse(c.filter)
		if err != nil {
			t.Errorf("failed to parse template: %v", err)
			continue
		}
		var escaped bytes.Buffer
		tmpl.Execute(&escaped, c.values)
		equal := cmp.Equal(escaped.String(), c.escaped)
		diff := cmp.Diff(escaped.String(), c.escaped)
		if !equal {
			t.Log(c.description)
			t.Errorf("unexpected escaped value: %s", diff)
		}
	}
}

// TestEscapeFilter tests escaping LDAP filters to avoid LDAP injection attacks
func TestEscapeFilter(t *testing.T) {
	// https://github.com/tcort/ldap-escape/blob/master/test/ldap-escape.test.js
	for _, c := range []struct {
		filter      string
		values      interface{}
		escaped     string
		description string
	}{
		{
			filter: "(uid={{.UID}})",
			values: struct {
				UID string
			}{
				UID: EscapeFilter("1337"),
			},
			escaped:     "(uid=1337)",
			description: "should work in the base case (no escaping)",
		},
		{
			filter: "(test={{.Test}})",
			values: struct {
				Test string
			}{
				Test: EscapeFilter(`Hi (This) = is * a \ test # ç à ô`),
			},
			escaped:     `(test=Hi \28This\29 = is \2a a \5c test # \c3\a7 \c3\a0 \c3\b4)`,
			description: "should correctly escape the OWASP Christmas Tree Example",
		},
		{
			filter: "{{.Filter}}",
			values: struct {
				Filter string
			}{
				Filter: EscapeFilter("foo=bar(baz)*"),
			},
			escaped:     `foo=bar\28baz\29\2a`,
			description: "should correctly escape the PHP test case",
		},
	} {
		tmpl, err := template.New("filter").Parse(c.filter)
		if err != nil {
			t.Errorf("failed to parse template: %v", err)
			continue
		}
		var escaped bytes.Buffer
		tmpl.Execute(&escaped, c.values)
		equal := cmp.Equal(escaped.String(), c.escaped)
		diff := cmp.Diff(escaped.String(), c.escaped)
		if !equal {
			t.Log(c.description)
			t.Errorf("unexpected escaped value: %s", diff)
		}
	}
}
