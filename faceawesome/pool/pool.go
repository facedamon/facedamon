package pool

import (
	"errors"
	"io"
	"log"
	"sync"
)

var ErrPoolClosed = errors.New("the pool had been closed")

type pool struct {
	m       sync.Mutex
	res     chan io.Closer
	factory func() (io.Closer, error)
	closed  bool
}

func p_new(fn func() (io.Closer, error), n uint) (*pool, error) {
	if n < 0 {
		return nil, errors.New("pool size less than zero")
	}
	return &pool{factory: fn, res: make(chan io.Closer, n)}, nil
}

func (p *pool) Get() (io.Closer, error) {
	select {
	case r, ok := <-p.res:
		if !ok {
			return nil, ErrPoolClosed
		}
		log.Println("共享资源")
		return r, nil
	default:
		log.Println("生成新资源")
		return p.factory()
	}
}

func (p *pool) Close() {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		return
	}

	// not write to chan
	p.closed = true
	close(p.res)

	// release the io
	for r := range p.res {
		r.Close()
	}
}

func (p *pool) Put(r io.Closer) {
	p.m.Lock()
	defer p.m.Unlock()

	if p.closed {
		r.Close()
		return
	}

	select {
	case p.res <- r:
		log.Println("将资源回收至资源池")
	default:
		log.Println("资源池满了，释放该资源")
		r.Close()
	}
}
