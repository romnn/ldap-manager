package harbor

import (
	// "encoding/json"
	"fmt"
	// "io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	// "github.com/go-ldap/ldap/v3"
	"github.com/romnn/ldap-manager/pkg"
	pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

// Test wraps a pre-configured harbor and LDAP Manager setup
type Test struct {
	LMTest pkg.Test
	Client http.Client
	// Container *Container
	// Manager   *LDAPManager
}

// Start starts the containers
func (test *Test) Start(t *testing.T) *Test {
	// var err error
	// t.Parallel()
	test.LMTest.Start(t)

	// // start OpenLDAP container
	// options := ContainerOptions{
	// 	Config: ldapconfig.NewConfig(),
	// }
	// container, err := StartOpenLDAP(context.Background(), options)
	// if err != nil {
	// 	t.Fatalf("failed to start OpenLDAP container: %v", err)
	// }
	// test.Container = &container

	// // create and setup the LDAP Manager service
	// test.Manager = NewLDAPManager(test.Container.Config)
	// test.Manager.DefaultAdminUsername = "ldapadmin"
	// test.Manager.DefaultAdminPassword = "123456"
	// if err := test.Manager.Connect(); err != nil {
	// 	t.Fatalf("failed to connect to OpenLDAP: %v", err)
	// }
	return test
}

// Setup runs the setup
func (test *Test) Setup(t *testing.T) *Test {
	test.LMTest.Setup(t)
	// if test.Manager == nil {
	// 	t.Fatal("must call test.Start(..) before running setup")
	// }
	// if err := test.Manager.Setup(); err != nil {
	// 	t.Fatalf("failed to setup manager: %v", err)
	// }
	return test
}

// Teardown stops the container
func (test *Test) Teardown() {
	test.LMTest.Teardown()
	// if test.Container != nil {
	// 	test.Container.Terminate(context.Background())
	// }
}

// TestHarborIntegration tests integration with the Harbor container registry
func TestHarborIntegration(t *testing.T) {
	test := new(Test).Start(t).Setup(t)
	defer test.Teardown()

	manager := test.LMTest.Manager
	config := manager.Config

	// ldapURI := config.URI()
	// t.Logf("dial %s", ldapURI)

	// // go func() {
	// // for {
	// // bind for the config CN to apply ACL rules
	// configDN := "cn=config"
	// configDN = "cn=admin,cn=config"
	// configDN = "cn=admin,dc=example,dc=org"
	// // configDN = "olcDatabase={1}mdb,cn=config"

	// // olcRootDN: cn=admin,cn=config
	// // olcRootPW: {SSHA}ImWUms6GLxtm4tNoEGRMsxRvFgq19GVI

	// // olcRootDN: cn=admin,dc=example,dc=org
	// // olcRootPW: {SSHA}jmzeBy7KhOsnV2dBOt6D7jUQZIRc7wz/

	// // configDN := fmt.Sprintf(
	// // 	"cn=%s,cn=config",
	// // 	"cn=config",
	// // 	// m.Config.AdminUsername,
	// // )
	// configPassword := "blabla123"
	// configPassword = "config"
	// configPassword = config.AdminPassword

	// // "config"
	// if err := ldapClient.Bind(configDN, configPassword); err != nil {
	// 	t.Errorf(
	// 		"unable to bind as %q with password %q: %v",
	// 		configDN, configPassword, err,
	// 	)
	// }
	// // time.Sleep(10 * time.Second)
	// // }
	// // }()

	// // "blabla123"
	// // time.Sleep(10 * time.Minute)
	// // t.Fatal("exit")
	// if err := manager.SetupReadOnlyUser(); err != nil {
	// 	t.Fatal(err)
	// }

	// assert that binding as the read-only search user works
	// searchUserDN := manager.UserDN(config.ReadOnlyUsername)
	searchUserDN := fmt.Sprintf(
		"cn=%s,%s",
		config.ReadOnlyUsername,
		config.BaseDN,
	)
	// searchUserDN = fmt.Sprintf("cn=%s", config.ReadOnlyUsername)
	searchPass := config.ReadOnlyPassword

	// time.Sleep(30 * time.Second)
	// t.Log("start")
	// go func() {
	// 	newLdapClient, err := ldap.DialURL(ldapURI)
	// 	if err != nil {
	// 		t.Errorf("failed to dial LDAP (%s): %v", ldapURI, err)
	// 	}

	// 	for {

	// 		// ldapClient.Unbind()
	// 		// if err := ldapClient.MD5Bind(searchUserDN, searchUserDN, searchPass); err != nil {
	// 		// 	t.Errorf(
	// 		// 		"failed to bind as user %q with password %q: %v",
	// 		// 		searchUserDN, searchPass, err,
	// 		// 	)
	// 		// }

	// 		// ldapClient.Unbind()
	// 		t.Log("attempt")
	// 		if err := newLdapClient.Bind(searchUserDN, searchPass); err != nil {
	// 			t.Errorf(
	// 				"failed to bind as user %q with password %q: %v",
	// 				searchUserDN, searchPass, err,
	// 			)
	// 		}
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 10 * time.Second,
	}, 10)

	// var ldapClient *ldap.Conn
	err := backoff.Retry(func() error {
		// var err error
		// ldapClient, err = ldap.DialURL(ldapURI)

		conn, err := manager.Pool.Get()
		if err != nil {
			t.Fatalf("failed to connect to LDAP: %v", err)
		}
		defer conn.Close()

		// if err := conn.Bind(searchUserDN, searchPass); err != nil {
		//   t.Fatalf(
		//     "warning: failed to bind as user %q with password %q: %v",
		//     searchUserDN, searchPass, err,
		//   )
		// }

		// if err != nil {
		// t.Logf("warning: failed to dial LDAP (%s): %v", ldapURI, err)
		// 	return err
		// }

		err = conn.Bind(searchUserDN, searchPass)
		if err != nil {
			t.Logf(
				"warning: failed to bind as user %q with password %q: %v",
				searchUserDN, searchPass, err,
			)
			// ldapClient.Close()
			// return err
		}
		return err
	}, b)

	if err != nil {
		t.Fatalf(
			"failed to bind as user %q with password %q: %v",
			searchUserDN, searchPass, err,
		)
	}

	// if err := ldapClient.Bind(searchUserDN, searchPass); err != nil {
	// time.Sleep(10 * time.Minute)
	// t.Log("done")
	// t.Log(ldapClient)
	// t.Fatal("exit")

	// "auth_mode": "ldap_auth",
	// "ldap_url": "{{ .Values.ldapmanager.ldap.protocol }}://{{ .Values.ldapmanager.ldap.host }}:{{ .Values.ldapmanager.ldap.port }}",
	// "ldap_base_dn": {{ .Values.ldapmanager.ldap.baseDN | quote }},
	// "ldap_search_dn": "cn={{ .Values.ldapmanager.ldap.readonly.user }},{{ .Values.ldapmanager.ldap.baseDN }}",
	// "ldap_search_password": {{ .Values.ldapmanager.ldap.readonly.password | quote }},
	// "ldap_uid": {{ .Values.ldapmanager.accountAttribute | quote }},
	// "ldap_scope": 2,
	// "ldap_filter": "objectclass=posixAccount",
	// "ldap_group_base_dn": "ou={{ .Values.ldapmanager.groupsOU }},{{ .Values.ldapmanager.ldap.baseDN }}",
	// "ldap_group_search_filter": "objectclass=posixGroup",
	// "ldap_group_search_scope": 2,
	// "ldap_group_attribute_name": "cn",
	// "ldap_group_admin_dn": "cn={{ .Values.ldapmanager.defaultAdminGroup }},ou={{ .Values.ldapmanager.groupsOU }},{{ .Values.ldapmanager.ldap.baseDN }}",
	// "ldap_group_membership_attribute": {{ .Values.ldapmanager.groupMembershipAttribute | quote }},
	// "self_registration": false

	req := UpdateConfigurationRequest{
		AuthMode: "ldap_auth",
		LdapURL: fmt.Sprintf(
			"%s://%s:%d",
			config.Protocol,
			"docker.for.mac.localhost",
			config.Port,
		),
		LdapBaseDN:   config.BaseDN,
		LdapSearchDN: searchUserDN,
		// fmt.Sprintf(
		// "cn=%s,%s",
		// config.ReadOnlyUserUsername,
		// config.BaseDN,
		// ),
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
		//   #ldap_filter: "{{ harbor_ldap_filter }}"
		//   #ldap_scope: "{{harbor_ldap_scope}}"
		//   #ldap_timeout: "{{ harbor_ldap_timeout }}"
		//   #ldap_uid: "{{ harbor_ldap_uid }}"
		//   #ldap_verify_cert: "{{ harbor_ldap_verify_cert }}"
		//   #ldap_group_admin_dn: "{{ harbor_ldap_group_admin_dn }}"
		//   #ldap_group_attribute_name: "{{ harbor_ldap_group_attribute_name }}"
		//   #ldap_group_base_dn: "{{ harbor_ldap_group_base_dn }}"
		//   #ldap_group_search_filter: "{{ harbor_ldap_group_search_filter }}"
		//   #ldap_group_search_scope: "{{ harbor_group_search_scope }}"
		//   #ldap_group_membership_attribute: "{{ harbor_group_membership_attribute }}"

	}
	body, err := ToJson(&req)
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
		&Auth{
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
		&Auth{
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
		&Auth{
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
		&Auth{
			Username: username, // "random",
			Password: password, // "shit",
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response.Status)
	t.Log(response.Body)

	// // ping the LDAP server
	// ldapPingURL, err := url.JoinPath(harborURL, "/api/v2.0/ldap/ping")
	// response, err = test.post(
	// 	ldapPingURL,
	// 	strings.NewReader(""),
	// 	&Auth{
	// 		Username: harborAdminUsername,
	// 		Password: harborAdminPassword,
	// 	},
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(response.Status)
	// t.Log(response.Body)

	// wait so that we can inspect
	time.Sleep(10 * time.Minute)

	// // get projects for the admin user
	// response, err = test.get(
	// 	projectsURL,
	// 	&Auth{
	// 		Username: harborAdminUsername,
	// 		Password: harborAdminPassword,
	// 	},
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(response.Status)
	// t.Log(response.Body)

	// // get projects for a user that does not exist
	// response, err = test.get(
	// 	projectsURL,
	// 	&Auth{
	// 		Username: "this user does not exist",
	// 		Password: "the password also makes no sense",
	// 	},
	// )
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(response.Status)
	// t.Log(response.Body)

	// // check if we can authenticate as the user
	// user, err := test.Manager.GetUser(username)
	// if err != nil {
	// t.Fatalf("failed to get user: %v", err)
	// }
	// t.Log(PrettyPrint(user))
}
