package pool2

import (
	"net"
	"sync"
)

// poolconn is wrapper around ne.conn to modify the behavior of
// net.conn`s close method
type PoolConn struct {
	net.Conn
	mu       sync.RWMutex
	c        *channelPool
	unusable bool
}

// close puts the given connects back to the pool
// instead of closing it
func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
		return nil
	}
	return p.c.Put(p.Conn)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// newConn wraps a standard net.conn to a poolConn net.conn
func (c *channelPool) wrapConn(conn net.Conn) net.Conn {
	return &PoolConn{c: c, Conn: conn}
}
