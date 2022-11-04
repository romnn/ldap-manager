package pool

import (
	"errors"
	"sync"

	"github.com/go-ldap/ldap/v3"
)

// channelPool implements the Pool interface using buffered channels.
type channelPool struct {
	mu       sync.Mutex
	connChan chan *Conn

	factory ConnectionFactory
	reset   ResetFunc
}

// NewChannelPool creates a new connection pool using a buffered channel
func NewChannelPool(capacity, maxCapacity uint, factory ConnectionFactory, reset ResetFunc) (Pool, error) {
	if maxCapacity == 0 {
		return nil, errors.New(
			"max pool capacity must be non zero",
		)
	}
	if capacity >= maxCapacity {
		return nil, errors.New(
			"max capacity must be greater or equal than capacity",
		)
	}
	pool := &channelPool{
		connChan: make(chan *Conn, maxCapacity),
		factory:  factory,
		reset:    reset,
	}

	// fill the pool with connections
	for i := uint(0); i < capacity; i++ {
		conn, err := factory()
		if err != nil {
			pool.Close()
			return nil, err
		}
		pool.connChan <- pool.wrapConn(conn)
	}

	return pool, nil
}

// Get returns a connection from the pool.
// If there is no connection available, a new connection is created.
func (pool *channelPool) Get() (*Conn, error) {
	if pool.connChan == nil {
		return nil, ErrClosed
	}

	select {
	case conn := <-pool.connChan:
		if conn == nil {
			return nil, ErrClosed
		}
		if conn.conn.IsClosing() {
			break
		}
		if conn.needReset {
			if err := pool.reset(conn.conn); err != nil {
				// connection is no longer useable
				conn.Close()
				break
			}
			conn.needReset = false
		}
		// connection most likely useable
		return conn, nil
	default:
	}

	// create a new connection
	newConn, err := pool.NewConnection()
	if err != nil {
		return nil, err
	}
	conn := pool.wrapConn(newConn)
	if err := pool.reset(newConn); err != nil {
		return nil, err
	}
	return conn, nil
}

// Len returns the current number of connections of the pool.
func (pool *channelPool) Len() int {
	return len(pool.connChan)
}

// Close closes the pool and all its connections
func (pool *channelPool) Close() {
	if pool.connChan == nil {
		return
	}
	// close all connections of the pool
	close(pool.connChan)
	for conn := range pool.connChan {
		conn.Close()
	}
}

// NewConnection returns a new LDAP connection
func (pool *channelPool) NewConnection() (ldap.Client, error) {
	conn, err := pool.factory()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// wrapConn wraps an LDAP client connection.
func (pool *channelPool) wrapConn(conn ldap.Client) *Conn {
	return &Conn{
		pool:      pool,
		conn:      conn,
		needReset: true,
	}
}

// put puts a connection back to the pool.
func (pool *channelPool) put(conn *Conn) {
	if conn == nil {
		return
	}

	if pool.connChan == nil {
		// pool is closed, close connection
		conn.Close()
		return
	}

	// put the connection back into the pool.
	select {
	case pool.connChan <- conn:
		return
	default:
		// pool is full, close connection
		conn.Close()
		return
	}
}
