package pool2

import (
	"log"
	"math/rand"
	"net"
	"sync"
	"testing"
	"time"
)

var (
	initialCap = 5
	maximumCap = 30
	network    = "tcp"
	address    = "127.0.0.1:7777"
	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
)

func init() {
	go simpleTCPServer()
	// wait until tcp server has been settled
	time.Sleep(time.Millisecond * 300)
	rand.Seed(time.Now().UTC().UnixNano())
}

func simpleTCPServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			buffer := make([]byte, 256)
			conn.Read(buffer)
		}()
	}
}

func newChannelPool() (Pool, error) {
	return NewChannelPool(initialCap, maximumCap, factory)
}

func TestChannelPool_Get_Impl(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	conn, err := p.Get()
	if err != nil {
		t.Errorf("get error : %s", err)
	}
	_, ok := conn.(*PoolConn)
	if !ok {
		t.Errorf("conn is not of type poolconn")
	}
}

func TestChannelPool_Get(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	_, err := p.Get()
	if err != nil {
		t.Errorf("get erros: %s", err)
	}

	// after one get, current capacity should be lowered by one
	if p.Len() != (initialCap - 1) {
		t.Errorf("get error. Expecting %d, got %d", initialCap-1, p.Len())
	}

	//get them all
	var wg sync.WaitGroup
	for i := 0; i < (initialCap - 1); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := p.Get()
			if err != nil {
				t.Errorf("get error : %s", err)
			}
		}()
	}
	wg.Wait()

	if p.Len() != 0 {
		t.Errorf("get error Expecting %d, got %d", initialCap-1, p.Len())
	}

	_, err = p.Get()
	if err != nil {
		t.Errorf("get error : %s", err)
	}
}

func TestChannelPool_Put(t *testing.T) {
	p, err := NewChannelPool(0, 30, factory)
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	// get from the pool
	conns := make([]net.Conn, maximumCap)
	for i := 0; i < maximumCap; i++ {
		// this get() will return PoolConn struct
		// so the next conn.close will call poolconn.close
		conn, _ := p.Get()
		conns[i] = conn
	}

	// now put them all back into the chan
	// call poolconn.close
	for _, conn := range conns {
		conn.Close()
	}

	if p.Len() != maximumCap {
		t.Errorf("Put error len. Expecting %d, got %d", 1, p.Len())
	}

	conn, _ := p.Get()
	p.Close() // close pool

	conn.Close() //try to put into a full pool

	if p.Len() != 0 {

	}
}
