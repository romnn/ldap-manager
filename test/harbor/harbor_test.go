package harbor

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/romnn/ldap-manager/pkg"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// Test wraps a pre-configured harbor and LDAP Manager setup
type Test struct {
	LMTest pkg.Test
	Client http.Client
}

// Start starts the containers
func (test *Test) Start(t *testing.T) *Test {
	test.LMTest.Start(t)
	return test
}

// Setup runs the setup
func (test *Test) Setup(t *testing.T) *Test {
	test.LMTest.Setup(t)
	return test
}

// Teardown stops the container
func (test *Test) Teardown() {
	test.LMTest.Teardown()
}

// TestHarborIntegration tests integration with the Harbor container registry
func TestHarborIntegration(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	manager := test.LMTest.Manager
	config := manager.Config

	// assert that binding as the read-only search user works
	searchUserDN := fmt.Sprintf(
		"cn=%s,%s",
		config.ReadOnlyUsername,
		config.BaseDN,
	)
	searchPass := config.ReadOnlyPassword

	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 10 * time.Second,
	}, 10)

	err := backoff.Retry(func() error {
		conn, err := manager.Pool.Get()
		if err != nil {
			t.Logf("failed to connect to LDAP: %v", err)
			return err
		}
		defer conn.Close()

		err = conn.Bind(searchUserDN, searchPass)
		if err != nil {
			t.Logf(
				"warning: failed to bind as user %q with password %q: %v",
				searchUserDN, searchPass, err,
			)
			return err
		}
		return nil
	}, b)

	if err != nil {
		t.Fatalf(
			"failed to bind as user %q with password %q: %v",
			searchUserDN, searchPass, err,
		)
	}

	req := updateConfigurationRequest{
		AuthMode: "ldap_auth",
		LdapURL: fmt.Sprintf(
			"%s://%s:%d",
			config.Protocol,
			"docker.for.mac.localhost",
			config.Port,
		),
		LdapBaseDN:         config.BaseDN,
		LdapSearchDN:       searchUserDN,
		LdapSearchPassword: config.ReadOnlyPassword,
		LdapUID:            manager.AccountAttribute,
		LdapScope:          2,
		LdapFilter:         "objectclass=posixAccount",
		LdapGroupBaseDN: fmt.Sprintf(
			"ou=%s,%s",
			manager.GroupsOU,
			config.BaseDN,
		),
		LdapGroupSearchFilter:  "objectclass=posixGroup",
		LdapGroupSearchScope:   2,
		LdapGroupAttributeName: "cn",
		LdapGroupAdminDN: fmt.Sprintf(
			"cn=%s,ou=%s,%s",
			manager.DefaultAdminGroup,
			manager.GroupsOU,
			config.BaseDN,
		),
		SelfRegistration: false,
	}

	body, err := toJson(&req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(body)

	// configure LDAP
	harborAdminUsername := "admin"
	harborAdminPassword := "Harbor12345"

	harborURL := "http://localhost:80"
	configURL, err := url.JoinPath(harborURL, "/api/v2.0/configurations")
	if err != nil {
		t.Fatal(err)
	}

	response, err := test.put(
		configURL,
		strings.NewReader(body),
		&auth{
			Username: harborAdminUsername,
			Password: harborAdminPassword,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status)
	t.Log(response.Body)

	// create new LDAP user
	username := "romnn"
	password := "hallo welt"
	if err := manager.NewUser(&pb.NewUserRequest{
		Username:  username,
		Password:  password,
		Email:     "a@b.de",
		FirstName: "roman",
		LastName:  "d",
	}); err != nil {
		t.Fatalf("failed to add user: %v", err)
	}

	// get projects for the new LDAP user
	projectsURL, err := url.JoinPath(harborURL, "/api/v2.0/projects")
	response, err = test.get(
		projectsURL,
		&auth{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status)
	t.Log(response.Body)

	// get projects for a random LDAP user
	projectsURL, err = url.JoinPath(harborURL, "/api/v2.0/projects")
	response, err = test.get(
		projectsURL,
		&auth{
			Username: "random",
			Password: "shit",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status)
	t.Log(response.Body)

	// get LDAP users
	usersURL, err := url.JoinPath(
		harborURL,
		fmt.Sprintf("/api/v2.0/ldap/users/search?username=%s", username),
	)
	t.Log(usersURL)
	response, err = test.get(
		usersURL,
		&auth{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status)
	t.Log(response.Body)

	// // ping the LDAP server (THIS WILL CRASH HARBOR)
	// ldapPingURL, err := url.JoinPath(harborURL, "/api/v2.0/ldap/ping")
	// response, err = test.post(
	// 	ldapPingURL,
	// 	strings.NewReader(""),
	// 	&auth{
	// 		Username: harborAdminUsername,
	// 		Password: harborAdminPassword,
	// 	},
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(response.Status)
	// t.Log(response.Body)

	// wait to inspect the containers
	// time.Sleep(10 * time.Minute)

	// // get projects for the admin user
	// response, err = test.get(
	// 	projectsURL,
	// 	&auth{
	// 		Username: harborAdminUsername,
	// 		Password: harborAdminPassword,
	// 	},
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(response.Status)
	// t.Log(response.Body)
}
