package net

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
)

type Pool struct {
	idlesConns  chan *idleConn
	reqQueue    []connReq
	maxCnt      int
	cnt         int
	maxIdleTime time.Duration
	factory     func() (net.Conn, error)
	lock        sync.Mutex
}

type idleConn struct {
	c              net.Conn
	lastActiveTime time.Time
}

type connReq struct {
	connChan chan net.Conn
}

func NewPool(initCnt int, maxIdleCnt int, maxCnt int, maxIdleTime time.Duration, factory func() (net.Conn, error)) (*Pool, error) {
	if initCnt > maxIdleCnt {
		return nil, fmt.Errorf("init count %d exceeds max idle count %d", initCnt, maxIdleCnt)
	}
	idlesConns := make(chan *idleConn, maxIdleCnt)
	for i := 0; i < initCnt; i++ {
		conn, err := factory()
		if err != nil {
			return nil, err
		}
		idlesConns <- &idleConn{c: conn, lastActiveTime: time.Now()}
	}
	res := &Pool{
		idlesConns:  idlesConns,
		maxCnt:      maxCnt,
		cnt:         0,
		maxIdleTime: maxIdleTime,
		factory:     factory,
	}
	return res, nil
}

func (p *Pool) Get(ctx context.Context) (net.Conn, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:

	}

	for {
		select {
		case ic := <-p.idlesConns:
			// todo 判断是否超时
			if ic.lastActiveTime.Add(p.maxIdleTime).Before(time.Now()) {
				_ = ic.c.Close()
				continue
			}
			return ic.c, nil
		default:
			p.lock.Lock()
			// todo 阻塞
			if p.cnt >= p.maxCnt {
				req := connReq{connChan: make(chan net.Conn, 1)}
				p.reqQueue = append(p.reqQueue, req)
				p.lock.Unlock()
				select {
				case <-ctx.Done():
					go func() {
						// 处理一：删除
						// 处理二：转发
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
			p.cnt++
			p.lock.Unlock()
			return c, nil
		}
	}
}

func (p *Pool) Put(ctx context.Context, c net.Conn) error {
	p.lock.Lock()
	if len(p.reqQueue) > 0 {
		req := p.reqQueue[0]
		p.reqQueue = p.reqQueue[1:]
		p.lock.Unlock()
		req.connChan <- c
		return nil
	}
	defer p.lock.Unlock()
	ic := &idleConn{
		c:              c,
		lastActiveTime: time.Now(),
	}
	select {
	case p.idlesConns <- ic:
	default:
		_ = c.Close()
		p.cnt--
	}

	return nil
}
