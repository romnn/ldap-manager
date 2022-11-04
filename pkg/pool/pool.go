package pool

import (
	"errors"

	"github.com/go-ldap/ldap/v3"
)

var (
	// ErrClosed is returned when the pool was closed
	ErrClosed = errors.New("pool closed")
)

// ConnectionFactory is a function to create new connections.
type ConnectionFactory func() (ldap.Client, error)

// ResetFunc is a function to reset connections.
type ResetFunc func(conn ldap.Client) error

// Pool is the connection pool interface
type Pool interface {
	// Get returns a new connection from the pool. Closing the connections puts
	// it back to the Pool. Closing it when the pool is destroyed or full will
	// be counted as an error.
	Get() (*Conn, error)

	// NewConnection returns a new LDAP connection not managed by the pool.
	NewConnection() (ldap.Client, error)

	// Close closes the pool and all its connections.
	Close()

	// Len returns the current number of connections of the pool.
	Len() int
}
