package pool3

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"testing"
)

var (
	initialCap = 5
	maximumCap = 30

	factory = func() (io.Closer, error) {
		return sql.Open("mysql",
			"damon:damon@tcp(127.0.0.1:3306)/damon?charset=utf8")
	}
)

func newChannelPool() (Pool, error) {
	return NewChannelPool(initialCap, maximumCap, factory)
}

func TestChannelPool_Get(t *testing.T) {
	p, _ := newChannelPool()
	defer p.Close()

	conn, err := p.Get()
	if err != nil {
		t.Errorf("get error : %s", err)
	}
	c, ok := conn.(*PoolConn)
	if !ok {
		t.Errorf("conn is not of type poolconn")
	}

	_, ok = c.Closer.(*sql.DB)
	if !ok {
		t.Errorf("conn is not of type sql.db")
	}
}
