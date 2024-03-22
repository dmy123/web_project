package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

var errKeyNotFound = errors.New("cache: 键不存在")

type BuildInMapCacheOption func(cache *BuildInMapCache)

type BuildInMapCache struct {
	data      map[string]*item
	mutex     sync.RWMutex
	close     chan struct{}
	onEvicted func(key string, val any)
	//onEvicted []func(key string, val any)
}

func NewBuildInMapCache(interval time.Duration, opts ...BuildInMapCacheOption) *BuildInMapCache {
	res := &BuildInMapCache{
		data:  make(map[string]*item, 100),
		close: make(chan struct{}),
		onEvicted: func(key string, val any) {

		},
	}

	for _, opt := range opts {
		opt(res)
	}

	go func() {
		ticker := time.NewTicker(interval)
		//for t := range ticker.C {
		//	res.mutex.Lock()
		//	i := 0
		//	for key, val := range res.data {
		//		if i > 1000{
		//			break
		//		}
		//		if !val.deadline.IsZero() && val.deadline.Before(t) {
		//			delete(res.data, key)
		//		}
		//		i++
		//
		//	}
		//	res.mutex.Unlock()
		//}
		select {
		case t := <-ticker.C:
			res.mutex.Lock()
			i := 0
			for key, val := range res.data {
				if i > 1000 {
					break
				}
				if !val.deadline.IsZero() && val.deadline.Before(t) {
					delete(res.data, key)
					res.delete(key)
				}
				i++

			}
			res.mutex.Unlock()
		case <-res.close:
			return
		}
	}()
	return res
}

func WithEvictedCallback(fn func(key string, val any)) BuildInMapCacheOption {
	return func(cache *BuildInMapCache) {
		cache.onEvicted = fn
	}
}

func (b *BuildInMapCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	//b.data[key] = val
	b.data[key] = &item{
		val:      val,
		deadline: time.Now().Add(expiration),
	}

	// 第10s设置key1=value1，过期时间1min
	// 30s设置key1=value2，过期时间1min
	// 封装value，判断是否与过期时间对应
	if expiration > 0 {
		time.AfterFunc(expiration, func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			//delete(b.data, key)
			val, ok := b.data[key]
			if ok && val.deadline.Before(time.Now()) {
				//delete(b.data, key)
				b.delete(key)
			}
		})
	}
	return nil
}

func (b *BuildInMapCache) set(ctx context.Context, key string, val any, expiration time.Duration) error {
	var dl time.Time
	if expiration > 0 {
		dl = time.Now().Add(expiration)
	}
	b.data[key] = &item{
		val:      val,
		deadline: dl,
	}
	return nil
}

func (b *BuildInMapCache) Get(ctx context.Context, key string) (any, error) {
	b.mutex.RLock()
	//defer b.mutex.RUnlock()
	res, exist := b.data[key]
	b.mutex.RUnlock()
	if !exist {
		return nil, fmt.Errorf("%w, key: %s", errKeyNotFound, key)
	}

	now := time.Now()
	if res.deadlineBefore(now) {
		b.mutex.Lock()
		defer b.mutex.Unlock()
		res, exist = b.data[key] // double check防止进入前新写入数据被删
		if res.deadlineBefore(now) {
			b.delete(key)
			//delete(b.data, key)
		}
		return nil, fmt.Errorf("%w, key: %s", errKeyNotFound, key)
	}
	return res.val, nil
}

func (b *BuildInMapCache) Delete(ctx context.Context, key string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	//delete(b.data, key)
	b.delete(key)
	return nil
}

func (b *BuildInMapCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	//delete(b.data, key)
	val, ok := b.data[key]
	if !ok {
		return nil, errKeyNotFound
	}
	b.delete(key)
	return val.val, nil
}

func (b *BuildInMapCache) delete(key string) {
	itm, ok := b.data[key]
	if !ok {
		return
	}
	delete(b.data, key)
	b.onEvicted(key, itm.val)
}

// 重复调用？
func (b *BuildInMapCache) Close(ctx context.Context) error {
	b.close <- struct{}{}
	return nil
}

type item struct {
	val      any
	deadline time.Time
}

func (i *item) deadlineBefore(t time.Time) bool {
	return !i.deadline.IsZero() && i.deadline.Before(t)
}
