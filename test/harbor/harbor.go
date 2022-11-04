package harbor

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	// "net/url"
	"strings"
	// "testing"
	// "time"

	// "github.com/cenkalti/backoff/v4"
	// "github.com/go-ldap/ldap/v3"
	"github.com/romnn/ldap-manager/pkg"
	// pb "github.com/romnn/ldap-manager/pkg/grpc/gen"
)

type UpdateConfigurationRequest struct {
	OidcVerifyCert                   bool   `json:"oidc_verify_cert,omitempty"`
	LdapSearchPassword               string `json:"ldap_search_password,omitempty"`
	EmailIdentity                    string `json:"email_identity,omitempty"`
	SkipAuditLogDatabase             bool   `json:"skip_audit_log_database,omitempty"`
	OidcExtraRedirectParms           string `json:"oidc_extra_redirect_parms,omitempty"`
	AuthMode                         string `json:"auth_mode,omitempty"`
	SelfRegistration                 bool   `json:"self_registration,omitempty"`
	HttpAuthProxyTokenreviewEndpoint string `json:"http_authproxy_tokenreview_endpoint,omitempty"`
	LdapSearchDN                     string `json:"ldap_search_dn,omitempty"`
	StoragePerProject                int    `json:"storage_per_project,omitempty"`
	HttpAuthProxyVerifyCert          bool   `json:"http_authproxy_verify_cert,omitempty"`
	EmailPassword                    string `json:"email_password,omitempty"`
	LdapGroupSearchFilter            string `json:"ldap_group_search_filter,omitempty"`
	UaaClientID                      string `json:"uaa_client_id,omitempty"`
	LdapTimeout                      int    `json:"ldap_timeout,omitempty"`
	LdapBaseDN                       string `json:"ldap_base_dn,omitempty"`
	LdapFilter                       string `json:"ldap_filter,omitempty"`
	ReadOnly                         bool   `json:"read_only,omitempty"`
	RobotTokenDuration               int    `json:"robot_token_duration,omitempty"`
	OidcAutoOnboard                  bool   `json:"oidc_auto_onboard,omitempty"`
	HttpAuthProxyServerCertificate   string `json:"http_authproxy_server_certificate,omitempty"`
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
	HttpAuthProxyAdminUsernames      string `json:"http_authproxy_admin_usernames,omitempty"`
	EmailUsername                    string `json:"email_username,omitempty"`
	HttpAuthProxyAdminGroups         string `json:"http_authproxy_admin_groups,omitempty"`
	OidcEndpoint                     string `json:"oidc_endpoint,omitempty"`
	HttpAuthProxyEndpoint            string `json:"http_authproxy_endpoint,omitempty"`
	OidcClientSecret                 string `json:"oidc_client_secret,omitempty"`
	OidcAdminGroup                   string `json:"oidc_admin_group,omitempty"`
	LdapScope                        int    `json:"ldap_scope,omitempty"`
	UaaEndpoint                      string `json:"uaa_endpoint,omitempty"`
	HttpAuthProxySkipSearch          bool   `json:"http_authproxy_skip_search,omitempty"`
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

type Auth struct {
	Username string
	Password string
}

type Response struct {
	Body       string
	Status     string
	StatusCode int
}

func NewResponse(res *http.Response) (*Response, error) {
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	response.Status = res.Status
	response.StatusCode = res.StatusCode

	fmt.Print(string(body))
	if len(body) > 0 {

    // interface 1
		var bodyJSON map[string]interface{}
		err = json.Unmarshal([]byte(body), &bodyJSON)
		if err == nil {
      response.Body = pkg.PrettyPrint(bodyJSON)
			return &response, nil
		}

    // interface 2
		var bodyJSON2 []map[string]interface{}
		err = json.Unmarshal([]byte(body), &bodyJSON2)
    if err == nil {
      response.Body = pkg.PrettyPrint(bodyJSON2)
			return &response, nil
		}
    return &response, err
	} else {
		response.Body = "empty response"
	}
	return &response, nil
}

func (test *Test) post(url string, body io.Reader, auth *Auth) (*Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	response, err := test.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return NewResponse(response)
}

func (test *Test) put(url string, body io.Reader, auth *Auth) (*Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	response, err := test.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return NewResponse(response)
}

func (test *Test) get(url string, auth *Auth) (*Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if auth != nil {
		req.SetBasicAuth(auth.Username, auth.Password)
	}
	response, err := test.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return NewResponse(response)
}

func ToJson(value interface{}) (string, error) {
	jsonValue, err := json.MarshalIndent(value, "", "    ")
	if err != nil {
		return "", err
	}
	return string(jsonValue), nil
}
