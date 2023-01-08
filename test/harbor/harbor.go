package harbor

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/k0kubun/pp/v3"
	"github.com/romnn/ldap-manager/pkg"
	ldapconfig "github.com/romnn/ldap-manager/pkg/config"
	"github.com/rs/xid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type updateConfigurationRequest struct {
	OidcVerifyCert                   bool   `json:"oidc_verify_cert,omitempty"`
	LdapSearchPassword               string `json:"ldap_search_password,omitempty"`
	EmailIdentity                    string `json:"email_identity,omitempty"`
	SkipAuditLogDatabase             bool   `json:"skip_audit_log_database,omitempty"`
	OidcExtraRedirectParms           string `json:"oidc_extra_redirect_parms,omitempty"`
	AuthMode                         string `json:"auth_mode,omitempty"`
	SelfRegistration                 bool   `json:"self_registration,omitempty"`
	HTTPAuthProxyTokenreviewEndpoint string `json:"http_authproxy_tokenreview_endpoint,omitempty"`
	LdapSearchDN                     string `json:"ldap_search_dn,omitempty"`
	StoragePerProject                int    `json:"storage_per_project,omitempty"`
	HTTPAuthProxyVerifyCert          bool   `json:"http_authproxy_verify_cert,omitempty"`
	EmailPassword                    string `json:"email_password,omitempty"`
	LdapGroupSearchFilter            string `json:"ldap_group_search_filter,omitempty"`
	UaaClientID                      string `json:"uaa_client_id,omitempty"`
	LdapTimeout                      int    `json:"ldap_timeout,omitempty"`
	LdapBaseDN                       string `json:"ldap_base_dn,omitempty"`
	LdapFilter                       string `json:"ldap_filter,omitempty"`
	ReadOnly                         bool   `json:"read_only,omitempty"`
	RobotTokenDuration               int    `json:"robot_token_duration,omitempty"`
	OidcAutoOnboard                  bool   `json:"oidc_auto_onboard,omitempty"`
	HTTPAuthProxyServerCertificate   string `json:"http_authproxy_server_certificate,omitempty"`
	OidcName                         string `json:"oidc_name,omitempty"`
	QuotaPerProjectEnable            bool   `json:"quota_per_project_enable,omitempty"`
	LdapURL                          string `json:"ldap_url,omitempty"`
	AuditLogForwardEndpoint          string `json:"audit_log_forward_endpoint,omitempty"`
	ProjectCreationRestriction       string `json:"project_creation_restriction,omitempty"`
	UaaClientSecret                  string `json:"uaa_client_secret,omitempty"`
	LdapUID                          string `json:"ldap_uid,omitempty"`
	LdapVerifyCert                   bool   `json:"ldap_verify_cert,omitempty"`
	OidcClientID                     string `json:"oidc_client_id,omitempty"`
	LdapGroupBaseDN                  string `json:"ldap_group_base_dn,omitempty"`
	LdapGroupAttributeName           string `json:"ldap_group_attribute_name,omitempty"`
	EmailInsecure                    bool   `json:"email_insecure,omitempty"`
	LdapGroupAdminDN                 string `json:"ldap_group_admin_dn,omitempty"`
	HTTPAuthProxyAdminUsernames      string `json:"http_authproxy_admin_usernames,omitempty"`
	EmailUsername                    string `json:"email_username,omitempty"`
	HTTPAuthProxyAdminGroups         string `json:"http_authproxy_admin_groups,omitempty"`
	OidcEndpoint                     string `json:"oidc_endpoint,omitempty"`
	HTTPAuthProxyEndpoint            string `json:"http_authproxy_endpoint,omitempty"`
	OidcClientSecret                 string `json:"oidc_client_secret,omitempty"`
	OidcAdminGroup                   string `json:"oidc_admin_group,omitempty"`
	LdapScope                        int    `json:"ldap_scope,omitempty"`
	UaaEndpoint                      string `json:"uaa_endpoint,omitempty"`
	HTTPAuthProxySkipSearch          bool   `json:"http_authproxy_skip_search,omitempty"`
	LdapGroupMembershipAttribute     string `json:"ldap_group_membership_attribute,omitempty"`
	OidcScope                        string `json:"oidc_scope,omitempty"`
	TokenExpiration                  int    `json:"token_expiration,omitempty"`
	NotificationEnable               bool   `json:"notification_enable,omitempty"`
	OidcUserClaim                    string `json:"oidc_user_claim,omitempty"`
	OidcGroupsClaim                  string `json:"oidc_groups_claim,omitempty"`
	EmailFrom                        string `json:"email_from,omitempty"`
	LdapGroupSearchScope             int    `json:"ldap_group_search_scope,omitempty"`
	EmailSsl                         bool   `json:"email_ssl,omitempty"`
	EmailPort                        int    `json:"email_port,omitempty"`
	RobotNamePrefix                  string `json:"robot_name_prefix,omitempty"`
	EmailHost                        string `json:"email_host,omitempty"`
	UaaVerifyCert                    bool   `json:"uaa_verify_cert,omitempty"`
}

type auth struct {
	Username string
	Password string
}

type response struct {
	Body       string
	DebugBody  string
	Status     string
	StatusCode int
}

func newResponse(httpRes *http.Response) (*response, error) {
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	var res response
	res.Body = string(body)
	res.DebugBody = string(body)
	res.Status = httpRes.Status
	res.StatusCode = httpRes.StatusCode

	fmt.Print(res.Body)
	if len(body) > 0 {

		// interface 1
		var bodyJSON map[string]interface{}
		err = json.Unmarshal([]byte(body), &bodyJSON)
		if err == nil {
			res.DebugBody = pkg.PrettyPrint(bodyJSON)
			return &res, nil
		}

		// interface 2
		var bodyJSON2 []map[string]interface{}
		err = json.Unmarshal([]byte(body), &bodyJSON2)
		if err == nil {
			res.DebugBody = pkg.PrettyPrint(bodyJSON2)
			return &res, nil
		}
	} else {
		res.Body = "empty response"
	}

	return &res, err
}

const (
	HarborTag        = "v2.6.1"
	RedisPort        = 6379
	PostgresPort     = 5432
	PostgresPassword = "root123"
	HarborPortalPort = 8080
	HarborCorePort   = 8080
	HarborProxyPort  = 8080
)

// Test wraps a pre-configured harbor and LDAP Manager setup
type Test struct {
	LMTest                pkg.Test
	Client                http.Client
	RedisContainer        testcontainers.Container
	PostgresContainer     testcontainers.Container
	HarborCoreContainer   testcontainers.Container
	HarborPortalContainer testcontainers.Container
	HarborProxyContainer  testcontainers.Container
	NetworkName           string
	Network               testcontainers.Network
}

func readEnvFile(path string) (map[string]string, error) {
	env := make(map[string]string)
	file, err := os.Open(path)
	if err != nil {
		return env, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		pair := strings.Split(scanner.Text(), "=")
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	return env, nil
}

type TestLogConsumer struct {
	Msgs   []string
	Prefix string
}

func (g *TestLogConsumer) Accept(l testcontainers.Log) {
	fmt.Printf("%s: %s", g.Prefix, string(l.Content))
	// g.Msgs = append(g.Msgs, string(l.Content))
}

func (test *Test) StartHarborCoreContainer(ctx context.Context) error {
	port, err := nat.NewPort("", strconv.Itoa(HarborCorePort))
	if err != nil {
		return fmt.Errorf("failed to build port: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %v", err)
	}

	// postgresHost, err := test.PostgresContainer.Host(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to get postgres host: %v", err)
	// }
	// fmt.Println(postgresHost)

	// realPostgresPort, err := test.PostgresContainer.MappedPort(ctx, postgresPort)
	// if err != nil {
	// 	return fmt.Errorf("failed to get exposed postgres port: %v", err)
	// }
	// fmt.Println(realPostgresPort)

	postgresPort, _ := nat.NewPort("", strconv.Itoa(PostgresPort))
	postgresIP, err := test.PostgresContainer.ContainerIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to get postgres container ip: %v", err)
	}

	// redisPort, _ := nat.NewPort("", strconv.Itoa(RedisPort))
	redisIP, err := test.RedisContainer.ContainerIP(ctx)
	if err != nil {
		return fmt.Errorf("failed to get redis container ip: %v", err)
	}
	redisURL, err := url.Parse(fmt.Sprintf("redis://%s:%d", redisIP, RedisPort))
	fmt.Println(err)
	fmt.Println(redisURL.String())

	env, err := readEnvFile(filepath.Join(cwd, "./common/config/core/env"))
	if err != nil {
		return err
	}

	env["POSTGRES_PASSWORD"] = PostgresPassword
	// "POSTGRESQL_HOST":   postgresHost,
	env["POSTGRESQL_HOST"] = postgresIP
	// "POSTGRESQL_HOST": test.NetworkName,
	// "POSTGRESQL_HOST": "postgres",
	// "POSTGRESQL_PORT":   strconv.Itoa(realPostgresPort.Int()),
	env["POSTGRESQL_PORT"] = strconv.Itoa(postgresPort.Int())
	// "REDIS_HOST":      redisIP,
	// "REDIS_PORT":      strconv.Itoa(redisPort.Int()),
	env["_REDIS_URL_CORE"] = redisURL.String()

	pp.Print(env)

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("goharbor/harbor-core:%s", HarborTag),
			Hostname:     "core",
			Networks:     []string{test.NetworkName},
			ExposedPorts: []string{string(port)},
			Env:          env, // map[string]string{},
			Binds: []string{
				// ./:/etc/core/ca/
				// - ./:/data/
				// - ./common/config/core/certificates/:/etc/core/certificates/
				// cwd + ":/etc/core/ca/",
				// cwd + ":/data/",
				filepath.Join(cwd, "./common/config/core/certificates/") + ":/etc/core/certificates/",
				// bind mounts
				filepath.Join(cwd, "./common/config/core/app.conf") + ":/etc/core/app.conf",
				filepath.Join(cwd, "./secret/core/private_key.pem") + ":/etc/core/private_key.pem",
				filepath.Join(cwd, "./secret/keys/secretkey") + ":/etc/core/key",
				filepath.Join(cwd, "./common/config/shared/trust-certificates") + ":/harbor_cust_cert",
			},
			WaitingFor: wait.ForListeningPort(port).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	}
	test.HarborCoreContainer, err = testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}

	g := TestLogConsumer{
		Msgs:   []string{},
		Prefix: "core",
	}

	if err := test.HarborCoreContainer.StartLogProducer(ctx); err != nil {
		// do something with err
	}

	test.HarborCoreContainer.FollowOutput(&g)
	return nil
}

func (test *Test) StartPostgresContainer(ctx context.Context) error {
	port, err := nat.NewPort("", strconv.Itoa(PostgresPort))
	if err != nil {
		return fmt.Errorf("failed to build port: %v", err)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("goharbor/harbor-db:%s", HarborTag),
			Hostname:     "postgresql",
			Networks:     []string{test.NetworkName},
			ExposedPorts: []string{string(port)},
			Env: map[string]string{
				"POSTGRES_PASSWORD": PostgresPassword,
			},
			ShmSize:    1024 * 1024 * 1024, // 1 GB
			WaitingFor: wait.ForListeningPort(port).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	}
	test.PostgresContainer, err = testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return fmt.Errorf(
			"failed to start container: %v",
			err,
		)
	}
	return nil
}

func (test *Test) StartHarborProxyContainer(ctx context.Context) error {
	port, err := nat.NewPort("", strconv.Itoa(HarborProxyPort))
	if err != nil {
		return fmt.Errorf("failed to build port: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %v", err)
	}

	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 3 * time.Second,
	}, 10)

	waiter, err := test.StartWaitContainer(ctx)
	if err != nil {
		return fmt.Errorf("failed to start waiter container: %v", err)
	}
	waitForHostnames := func() error {
		// ret, reader, err := waiter.Exec(ctx, []string{"curl", "--fail", "-s", "http://core:8080"})
		// "--fail", "-s",
		ret, reader, err := waiter.Exec(ctx, []string{"curl", "-q", "http://core:8080"})
		output, err := io.ReadAll(reader)
		fmt.Printf("proxy health check return code: %d error: %v output: %s\n", ret, err, string(output))
		if err != nil {
			return err
		}
		if ret != 0 {
			return fmt.Errorf("wait returned bad exit code: %d", ret)
		}
		return nil
	}
	if err := backoff.Retry(waitForHostnames, b); err != nil {
		return err
	}

	// startContainer := func() error {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("goharbor/nginx-photon:%s", HarborTag),
			Hostname:     "proxy",
			Networks:     []string{test.NetworkName},
			ExposedPorts: []string{string(port)},
			Binds: []string{
				filepath.Join(cwd, "./common/config/nginx") + ":/etc/nginx",
				filepath.Join(cwd, "./common/config/shared/trust-certificates") + ":/harbor_cust_cert",
			},
			WaitingFor: wait.ForListeningPort(port).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	}
	test.HarborProxyContainer, err = testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}
	fmt.Println(test.HarborProxyContainer.MappedPort(ctx, port))

	g := TestLogConsumer{
		Msgs:   []string{},
		Prefix: "proxy",
	}

	if err := test.HarborProxyContainer.StartLogProducer(ctx); err != nil {
		// do something with err
	}

	test.HarborProxyContainer.FollowOutput(&g)
	return nil
	// }

	// return backoff.Retry(startContainer, b)
}

func (test *Test) StartHarborPortalContainer(ctx context.Context) error {
	port, err := nat.NewPort("", strconv.Itoa(HarborPortalPort))
	if err != nil {
		return fmt.Errorf("failed to build port: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working dir: %v", err)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("goharbor/harbor-portal:%s", HarborTag),
			Hostname:     "portal",
			Networks:     []string{test.NetworkName},
			ExposedPorts: []string{string(port)},
			Binds: []string{
				filepath.Join(cwd, "./common/config/portal/nginx.conf") + ":/etc/nginx/nginx.conf",
			},
			WaitingFor: wait.ForListeningPort(port).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	}
	test.HarborPortalContainer, err = testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}
	return nil
}

func (test *Test) StartWaitContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:    "curlimages/curl",
			Cmd:      []string{"/bin/sh", "-c", "--", "while true; do sleep 30; done;"},
			Networks: []string{test.NetworkName},
		},
		Started: true,
	}
	return testcontainers.GenericContainer(ctx, req)
}

func (test *Test) StartRedisContainer(ctx context.Context) error {
	port, err := nat.NewPort("", strconv.Itoa(RedisPort))
	if err != nil {
		return fmt.Errorf("failed to build port: %v", err)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        fmt.Sprintf("goharbor/redis-photon:%s", HarborTag),
			Hostname:     "redis",
			Networks:     []string{test.NetworkName},
			ExposedPorts: []string{string(port)},
			WaitingFor:   wait.ForListeningPort(port).WithStartupTimeout(5 * time.Minute),
		},
		Started: true,
	}
	test.RedisContainer, err = testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to start container: %v", err)
	}
	return nil
}

// CreateNetwork creates a docker network for the harbor services
func (test *Test) CreateNetwork(ctx context.Context) error {
	test.NetworkName = fmt.Sprintf("harbor-network-%s", xid.New().String())
	request := testcontainers.NetworkRequest{
		Driver:         "bridge",
		Name:           test.NetworkName,
		Attachable:     true,
		CheckDuplicate: true,
	}

	createNetwork := func() error {
		var err error
		test.Network, err = testcontainers.GenericNetwork(
			ctx,
			testcontainers.GenericNetworkRequest{
				NetworkRequest: request,
			},
		)
		return err
	}
	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 1 * time.Second,
	}, 10)
	return backoff.Retry(createNetwork, b)
}

// Start starts the containers
func (test *Test) Start(t *testing.T) *Test {
	if err := test.CreateNetwork(context.Background()); err != nil {
		t.Fatalf("failed to create docker bridge network: %v", err)
	}
	if err := test.StartPostgresContainer(context.Background()); err != nil {
		t.Fatalf("failed to start postgres container: %v", err)
	}
	if err := test.StartRedisContainer(context.Background()); err != nil {
		t.Fatalf("failed to start redis container: %v", err)
	}
	if err := test.StartHarborCoreContainer(context.Background()); err != nil {
		t.Fatalf("failed to start harbor core container: %v", err)
	}
	if err := test.StartHarborPortalContainer(context.Background()); err != nil {
		t.Fatalf("failed to start harbor portal container: %v", err)
	}
	if err := test.StartHarborProxyContainer(context.Background()); err != nil {
		t.Fatalf("failed to start harbor proxy container: %v", err)
	}

	// start OpenLDAP container
	options := pkg.ContainerOptions{
		Config:   ldapconfig.NewConfig(),
		Networks: []string{test.NetworkName},
	}
	container, err := pkg.StartOpenLDAP(context.Background(), options)
	if err != nil {
		t.Fatalf("failed to start OpenLDAP container: %v", err)
	}
	test.LMTest.Container = &container

	// create and setup the LDAP Manager service
	test.LMTest.Manager = pkg.NewLDAPManager(test.LMTest.Container.Config)
	test.LMTest.Manager.DefaultAdminUsername = "ldapadmin"
	test.LMTest.Manager.DefaultAdminPassword = "123456"
	if err := test.LMTest.Manager.Connect(); err != nil {
		t.Fatalf("failed to connect to OpenLDAP: %v", err)
	}
	return test
}

// Setup runs the setup
func (test *Test) Setup(t *testing.T) *Test {
	test.LMTest.Setup(t)
	return test
}

// Teardown stops the container
func (test *Test) Teardown() {
	for _, container := range []testcontainers.Container{
		test.RedisContainer,
		test.PostgresContainer,
		test.HarborPortalContainer,
		test.HarborCoreContainer,
		test.HarborProxyContainer,
	} {
		if container != nil {
			container.Terminate(context.Background())
		}
	}
	test.LMTest.Teardown()
	if test.Network != nil {
		test.Network.Remove(context.Background())
	}
}

func (test *Test) post(url string, body io.Reader, auth *auth) (*response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	res, err := test.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return newResponse(res)
}

func (test *Test) put(url string, body io.Reader, auth *auth) (*response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	res, err := test.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return newResponse(res)
}

func (test *Test) get(url string, auth *auth) (*response, error) {
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	res, err := test.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return newResponse(res)
}

func toJSON(value interface{}) (string, error) {
	jsonValue, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}
