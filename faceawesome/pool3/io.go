package pool3

import (
	"io"
	"sync"
)

// poolconn is wrapper around io.closer to modify the behavior of
// io.closer`s close method
type PoolConn struct {
	io.Closer
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
		if p.Closer != nil {
			return p.Closer.Close()
		}
		return nil
	}
	return p.c.Put(p.Closer)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}

// newConn wraps a standard net.conn to a poolConn io.closer
func (c *channelPool) wrapConn(closer io.Closer) io.Closer {
	return &PoolConn{c: c, Closer: closer}
}
