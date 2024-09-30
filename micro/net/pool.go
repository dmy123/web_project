package net

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Pool struct {
	idleConns   chan *idleConn // 空闲连接
	reqConn     []connReq      // 请求队列
	maxCnts     int
	cnts        int
	maxIdleTime time.Duration
	factory     func() (net.Conn, error)
	lock        sync.Mutex
}

type idleConn struct {
	conn           net.Conn
	lastActiveTime time.Time
}

type connReq struct {
	connChan chan net.Conn
}

func NewPool(initCnt, maxIdleCnt, maxCnts int, maxIdleTime time.Duration, factory func() (net.Conn, error)) (*Pool, error) {
	if initCnt > maxIdleCnt {
		return nil, fmt.Errorf("init count %d exceeds max idle count %d", initCnt, maxIdleCnt)
	}
	idlesConns := make(chan *idleConn, maxIdleCnt)
	for i := 0; i < initCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idlesConns <- &idleConn{conn: conn, lastActiveTime: time.Now()}
	}
	return &Pool{
		idleConns:   idlesConns,
		maxCnts:     maxCnts,
		cnts:        0,
		maxIdleTime: maxIdleTime,
		factory:     factory,
	}, nil
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	for {
		select {
		case ic := <-p.idleConns:
			if ic.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = ic.conn.Close()
				continue
			}
			return ic.conn, nil
		default:
			p.lock.Lock()
			if p.cnts >= p.maxCnts {
				req := connReq{connChan: make(chan net.Conn, 1)}
				p.reqConn = append(p.reqConn, req)
				p.lock.Unlock()
				select {
				case <-ctx.Done():
					go func() {
						c := <-req.connChan
						_ = p.Put(ctx, c)
					}()
					return nil, ctx.Err()
				case ic := <-req.connChan:
					return ic, nil
				}
			}
			c, err := p.factory()
			if err != nil {
				return nil, err
			}
			p.cnts++
			p.lock.Unlock()
			return c, nil
		}
	}
}

func (p *Pool) Put(ctx context.Context, c net.Conn) error {
	p.lock.Lock()
	if len(p.reqConn) > 0 {
		req := p.reqConn[0]
		p.reqConn = p.reqConn[1:]
		p.lock.Unlock()
		req.connChan <- c
		return nil
	}
	defer p.lock.Unlock()
	ic := &idleConn{conn: c, lastActiveTime: time.Now()}
	select {
	case p.idleConns <- ic:
	default:
		_ = c.Close()
		p.cnts--

	}
	return nil
}
