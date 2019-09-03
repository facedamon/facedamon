package pool

import (
	"io"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const max = 5
const pools = 2

type db struct {
	id int32
}

func (d db) Close() error {
	log.Println("db close", d.id)
	return nil
}

var idcounter int32

func createDb() (io.Closer, error) {
	id := atomic.AddInt32(&idcounter, 1)
	return &db{id}, nil
}

func createDb2() interface{} {
	id := atomic.AddInt32(&idcounter, 1)
	return &db{id}
}

func query(q int, p *pool) {
	conn, err := p.Get()
	if err != nil {
		log.Println(err)
		return
	}
	defer p.Put(conn)
	//模拟查询
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("第%d个查询，使用的是ID为%d的数据库连接", q, conn.(*db).id)
}

func query2(q int, p *sync.Pool) {
	conn := p.Get().(*db)
	defer p.Put(conn)

	//模拟查询
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("第%d个查询，使用的是ID为%d的数据库连接", q, conn.id)
}

func TestPool_(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(max)

	p, err := p_new(createDb, pools)
	if err != nil {
		log.Println(err)
		return
	}

	for q := 0; q < max; q++ {
		go func(t int) {
			query(t, p)
			wg.Done()
		}(q)
	}

	wg.Wait()
	log.Println("开始关闭资源池")
	p.Close()
}

func TestSyncPool(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(max)

	p := &sync.Pool{New: createDb2}

	for q := 0; q < max; q++ {
		go func(t int) {
			query2(t, p)
			wg.Done()
		}(q)
	}

	wg.Wait()

}
