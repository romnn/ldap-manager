package pool

import (
	"crypto/tls"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/go-ldap/ldap/v3"
	log "github.com/sirupsen/logrus"
)

// Conn implements Client to override the Close() method
type Conn struct {
	conn         ldap.Client
	pool         *channelPool
	bindUser     string
	bindPassword string
	needReset    bool
	unuseable    bool
}

// Start starts connection
func (c *Conn) Start() {
	c.conn.Start()
}

// StartTLS starts TLS connection
func (c *Conn) StartTLS(config *tls.Config) error {
	return c.conn.StartTLS(config)
}

// Close puts the connection back to the pool instead of closing it
func (c *Conn) Close() {
	if c.unuseable {
		c.conn.Close()
		return
	}
	c.pool.put(c)
}

// withRetry performs an operation and retries on temporary failures
func (c *Conn) withRetry(operation func() error) error {
	b := backoff.WithMaxRetries(&backoff.ConstantBackOff{
		Interval: 2 * time.Second,
	}, 10)

	err := backoff.Retry(func() error {
		if err := operation(); err != nil {
			connectionErr := ldap.IsErrorAnyOf(
				err,
				ldap.LDAPResultConnectError,
				ldap.ErrorNetwork,
			)
			tempErr := ldap.IsErrorAnyOf(
				err,
				ldap.LDAPResultTimeLimitExceeded,
				ldap.LDAPResultSaslBindInProgress,
				ldap.LDAPResultBusy,
				ldap.LDAPResultUnavailable,
				ldap.LDAPResultServerDown,
				ldap.LDAPResultTimeout,
				ldap.LDAPResultTooLate,
				ldap.LDAPResultSyncRefreshRequired,
			)
			if connectionErr {
				// mark connection as unuseable
				c.unuseable = true
			}
			if connectionErr || tempErr || c.conn.IsClosing() {
				if conn, err := c.pool.NewConnection(); err == nil {
					c.conn.Close()
					log.Info("swapped connection")
					c.conn = conn
					if _, err := c.SimpleBind(&ldap.SimpleBindRequest{
						Username: c.bindUser,
						Password: c.bindPassword,
						Controls: []ldap.Control{},
					}); err != nil {
						return err
					}
					c.unuseable = false
				}

				return err
			}
			return backoff.Permanent(err)
		}
		return nil
	}, b)

	if err, ok := err.(*backoff.PermanentError); ok {
		// return the underlying permanent error
		return err.Err
	}
	// number of retries exceeded for temporary error
	return err
}

// SimpleBind wraps the SimpleBind LDAP client method
func (c *Conn) SimpleBind(simpleBindRequest *ldap.SimpleBindRequest) (*ldap.SimpleBindResult, error) {
	c.needReset = true
	c.bindUser = simpleBindRequest.Username
	c.bindPassword = simpleBindRequest.Password
	var result *ldap.SimpleBindResult
	err := c.withRetry(func() error {
		var err error
		result, err = c.conn.SimpleBind(simpleBindRequest)
		return err
	})
	return result, err
}

// Bind wraps the Bind LDAP client method
func (c *Conn) Bind(username, password string) error {
	c.needReset = true
	c.bindUser = username
	c.bindPassword = password
	return c.withRetry(func() error {
		return c.conn.Bind(username, password)
	})
}

// Add wraps the Add LDAP client method
func (c *Conn) Add(addRequest *ldap.AddRequest) error {
	return c.withRetry(func() error {
		return c.conn.Add(addRequest)
	})
}

// Del wraps the Del LDAP client method
func (c *Conn) Del(delRequest *ldap.DelRequest) error {
	return c.withRetry(func() error {
		return c.conn.Del(delRequest)
	})
}

// Modify wraps the Modify LDAP client method
func (c *Conn) Modify(modifyRequest *ldap.ModifyRequest) error {
	return c.withRetry(func() error {
		return c.conn.Modify(modifyRequest)
	})
}

// ModifyDN wraps the ModifyDN LDAP client method
func (c *Conn) ModifyDN(modifyDnRequest *ldap.ModifyDNRequest) error {
	return c.withRetry(func() error {
		return c.conn.ModifyDN(modifyDnRequest)
	})
}

// Compare wraps the Compare LDAP client method
func (c *Conn) Compare(dn, attribute, value string) (bool, error) {
	var result bool
	err := c.withRetry(func() error {
		var err error
		result, err = c.conn.Compare(dn, attribute, value)
		return err
	})
	return result, err
}

// PasswordModify wraps the PasswordModify LDAP client method
func (c *Conn) PasswordModify(passwordModifyRequest *ldap.PasswordModifyRequest) (*ldap.PasswordModifyResult, error) {
	var result *ldap.PasswordModifyResult
	err := c.withRetry(func() error {
		var err error
		result, err = c.conn.PasswordModify(passwordModifyRequest)
		return err
	})
	return result, err
}

// Search wraps the Search LDAP client method
func (c *Conn) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	var result *ldap.SearchResult
	err := c.withRetry(func() error {
		var err error
		result, err = c.conn.Search(searchRequest)
		return err
	})
	return result, err
}

// SearchWithPaging wraps the SearchWithPaging LDAP client method
func (c *Conn) SearchWithPaging(searchRequest *ldap.SearchRequest, pagingSize uint32) (*ldap.SearchResult, error) {
	var result *ldap.SearchResult
	err := c.withRetry(func() error {
		var err error
		result, err = c.conn.SearchWithPaging(searchRequest, pagingSize)
		return err
	})
	return result, err
}
