package pool3

import (
	"errors"
	"io"
)

var ErrClosed = errors.New("pool is closed")

// pool interface describes a pool impl
// a pool should have maximum capacity
// an ideal pool is thread safe and easy to use
type Pool interface {
	//get returns a new conn from the poll
	// closing the conn puts it back to the poll.
	// closing it when the poll is destroyed or full
	// will be counted as an error
	Get() (io.Closer, error)

	// close the pool and all its conn.
	// after close the poll is no longer usable
	Close()

	// len returns the current number of conn of the pool
	Len() int
}
