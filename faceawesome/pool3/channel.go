package pool3

import (
	"errors"
	"fmt"
	"io"
	"log"
	"sync"
)

// factory is a function to create new conn
type Factory func() (io.Closer, error)

// channelPool impl the pool interface
// based on buffered channels
type channelPool struct {
	mu      sync.RWMutex
	conns   chan io.Closer
	factory Factory
}

func init() {
	log.SetFlags(log.Ldate | log.Lshortfile)
	log.SetPrefix("[init channel pool]")
}

// newchannelpool returns a new pool based on buffered channels with an initial capacity
// and maximum capacity. factory is used when initial capacity is greater than zero to fill
// the pool. a zero initialCap doesnt fill the pool until a new Get() is called.
// during a Get(), if there is no new connection available in the pool
// a new conn will be created via the factory method
func NewChannelPool(initialCap, maxCap int, fn Factory) (Pool, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}
	c := &channelPool{conns: make(chan io.Closer, maxCap), factory: fn}

	// create initial conn, if something goes wrong
	// just close the poll error out
	for i := 0; i < initialCap; i++ {
		conn, err := fn()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- conn
	}
	return c, nil
}

func (c *channelPool) getConnsAndFactory() (chan io.Closer, Factory) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	conns := c.conns
	factory := c.factory
	return conns, factory
}

// get impl the pool interface get method.
// if there is no new connection available in the pool,
// a new connection will be created via the factory method
func (c *channelPool) Get() (io.Closer, error) {
	conns, factory := c.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}

	// wrap our connections with out custom net.conn
	// that puts the connection back to the pool if its closed
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}
		return c.wrapConn(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		return c.wrapConn(conn), nil
	}
}

// put puts the connection back to the pool
// if the pool is full or closed.
// conn is simply closed, a nil conn will br rejected
func (c *channelPool) Put(conn io.Closer) error {
	if conn == nil {
		return errors.New("connections is nil. rejecting")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conns == nil {
		// pool is closed
		return conn.Close()
	}

	select {
	// put the resource back into the pool
	case c.conns <- conn:
		return nil
	//pool is full. close passed connection
	default:
		return conn.Close()
	}
}

// impl PoolConn`s closed
func (c *channelPool) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}
	// close chan
	close(conns)
	for conn := range conns {
		//close PoolConn io.Net
		conn.Close()
	}
}

func (c *channelPool) Len() int {
	conns, _ := c.getConnsAndFactory()
	return len(conns)
}
