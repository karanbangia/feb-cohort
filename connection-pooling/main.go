package main

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//test sql connection pooling
	benchmarkNonPool()
	//benchmarkPool()
}

func benchmarkNonPool() {
	startTime := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			username := "dbeaver"
			password := "dbeaver"
			host := "127.0.0.1"
			port := "3306"

			// Connect without specifying a database
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)
			db, err := sql.Open("mysql", dsn)
			if err != nil {
				panic(err)
			}
			_, err = db.Exec("SELECT SLEEP(0.1)")
			if err != nil {
				panic(err)
			}
			db.Close()
		}()
	}
	wg.Wait()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	//elapsed time in seconds
	fmt.Println("Elapsed time: ", elapsedTime.Seconds(), "s")
}

func benchmarkPool() {
	startTime := time.Now()
	wg := sync.WaitGroup{}
	connPool := NewConnPool(1000)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			db, err := connPool.GetConn(5 * time.Second)
			if err != nil {
				panic(err)
			}
			_, err = db.Exec("SELECT SLEEP(0.1)")
			if err != nil {
				panic(err)
			}
			connPool.ReleaseConn(db)
		}()
	}
	wg.Wait()
	connPool.Close()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	//elapsed time in seconds
	fmt.Println("Elapsed time: ", elapsedTime.Seconds(), "s")
}

type ConnPool struct {
	conn     chan *sql.DB
	mu       *sync.Mutex
	maxConn  int
	isClosed bool
}

func NewConnPool(maxConn int) *ConnPool {
	poolChan := make(chan *sql.DB, maxConn)
	for i := 0; i < maxConn; i++ {
		//create new connection
		//add to conn
		username := "dbeaver"
		password := "dbeaver"
		host := "127.0.0.1"
		port := "3306"

		// Connect without specifying a database
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", username, password, host, port)
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		poolChan <- conn
	}
	return &ConnPool{
		conn:     poolChan,
		mu:       &sync.Mutex{},
		maxConn:  maxConn,
		isClosed: false,
	}
}

func (p *ConnPool) GetConn(timeout time.Duration) (*sql.DB, error) {
	select {
	case conn := <-p.conn:
		return conn, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout: no connection avaialble")
	}
}

func (p *ConnPool) ReleaseConn(conn *sql.DB) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isClosed {
		conn.Close()
		return
	}
	p.conn <- conn
}

func (p *ConnPool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.isClosed {
		return
	}
	close(p.conn)
	for conn := range p.conn { // range only happens when channel is closed
		conn.Close()
	}
	p.isClosed = true
}
