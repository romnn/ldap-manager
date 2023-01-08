package harbor

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/romnn/ldap-manager/pkg"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

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

	// retry 10 times
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
		}
		return err
	}, b)

	if err != nil {
		t.Fatalf(
			"failed to bind as user %q with password %q: %v",
			searchUserDN, searchPass, err,
		)
	}

	harborHost, err := test.HarborCoreContainer.Host(context.TODO())
	if err != nil {
		t.Fatalf("failed to get container host: %v", err)
	}

	harborPort, _ := nat.NewPort("", strconv.Itoa(harborCorePort))
	realHarborPort, err := test.HarborCoreContainer.MappedPort(context.TODO(), harborPort)
	if err != nil {
		t.Fatalf("failed to get exposed container port: %v", err)
	}

	harborURL := fmt.Sprintf("http://%s:%d", harborHost, realHarborPort.Int())
	// harborURL = "http://localhost:80"

	test.configureHarborLDAPAuthentication(t, harborURL, searchUserDN)

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
	test.getProjectsForLDAPUsers(t, harborURL, username, password)

	// get LDAP users
	test.getLDAPUsers(t, harborURL, username, password)
}

func (test *Test) configureHarborLDAPAuthentication(t *testing.T, harborURL, searchUserDN string) {
	manager := test.LMTest.Manager
	config := manager.Config

	ldapIP, err := test.LMTest.Container.Container.ContainerIP(context.TODO())
	if err != nil {
		t.Fatalf("failed to get ldap container ip: %v", err)
	}

	ldapURL := fmt.Sprintf(
		"%s://%s:%d",
		config.Protocol,
		ldapIP,
		pkg.OpenLDAPPort,
	)
	req := updateConfigurationRequest{
		AuthMode:           "ldap_auth",
		LdapURL:            ldapURL,
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

	body, err := toJSON(&req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(body)

	configURL, err := url.JoinPath(harborURL, "/api/v2.0/configurations")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(configURL)

	// configure LDAP authentication
	response, err := test.put(
		configURL,
		strings.NewReader(body),
		&auth{
			Username: test.HarborAdminUsername,
			Password: test.HarborAdminPassword,
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status)
	t.Log(response.DebugBody)
	if response.StatusCode != 200 {
		t.Fatalf("configuring LDAP returned bad status code %s", response.Status)
	}
}

func (test *Test) getProjectsForLDAPUsers(t *testing.T, harborURL, username, password string) {
	projectsURL, err := url.JoinPath(harborURL, "/api/v2.0/projects")
	response, err := test.get(
		projectsURL,
		&auth{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(response.Status)
	t.Log(response.DebugBody)
	if response.StatusCode != 200 {
		t.Errorf("user query for user %q returned bad status code %s", username, response.Status)
	}

	// get projects for a random LDAP user
	randomUsername := "random"
	projectsURL, err = url.JoinPath(harborURL, "/api/v2.0/projects")
	response, err = test.get(
		projectsURL,
		&auth{
			Username: randomUsername,
			Password: "shit",
		},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(response.Status)
	t.Log(response.DebugBody)
	if response.StatusCode != 200 {
		t.Errorf("user query for user %q returned bad status code %s", randomUsername, response.Status)
	}
}

func (test *Test) getLDAPUsers(t *testing.T, harborURL, username, password string) {
	query := url.Values{}
	query.Set("username", username)
	usersURL, err := url.JoinPath(
		harborURL,
		"/api/v2.0/ldap/users/search",
	)
	response, err := test.get(
		usersURL,
		&auth{
			Username: test.HarborAdminUsername,
			Password: test.HarborAdminPassword,
		},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(response.Status)
	t.Log(response.DebugBody)
	if response.StatusCode != 200 {
		t.Errorf(
			"user query for harbor admin %q returned bad status code %s",
			test.HarborAdminUsername,
			response.Status,
		)
	}
	var queryUsers []map[string]interface{}
	if err := json.Unmarshal([]byte(response.Body), &queryUsers); err != nil {
		t.Errorf("failed to parse users JSON: %v", err)
	}
	if len(queryUsers) != 2 {
		t.Errorf(
			"user query for harbor admin %q returned %d users, expected 2",
			test.HarborAdminUsername,
			len(queryUsers),
		)
	}

	response, err = test.get(
		usersURL,
		&auth{
			Username: username,
			Password: password,
		},
	)
	if err != nil {
		t.Error(err)
	}
	t.Log(response.Status)
	t.Log(response.DebugBody)
	if response.StatusCode != 403 {
		t.Errorf("user query for LDAP user %q returned bad status code %s", username, response.Status)
	}
}
