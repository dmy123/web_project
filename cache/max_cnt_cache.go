package cache

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

var errOverCapacity = errors.New("cache: 超过容量限制")

type MaxCntCache struct {
	*BuildInMapCache
	cnt    int32
	maxCnt int32
}

func NewMaxCntCache(c *BuildInMapCache, maxCnt int32) *MaxCntCache {
	res := &MaxCntCache{
		BuildInMapCache: c,
		maxCnt:          maxCnt,
	}
	origin := c.onEvicted
	res.onEvicted = func(key string, val any) {
		atomic.AddInt32(&res.cnt, -1)
		if origin != nil {
			origin(key, val)
		}
	}
	return res
}

func (c *MaxCntCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	////法一：key已存在时，计数不准
	//cnt := atomic.AddInt32(&c.cnt, 1)
	//if cnt > c.maxCnt {
	//	atomic.AddInt32(&c.cnt, -1)
	//	return errOverCapacity
	//}
	//return c.BuildInMapCache.Set(ctx, key, val, expiration)

	////法二：可能会两个人都走到！ok从而cnt加两次
	//c.mutex.Lock()
	//_, ok := c.data[key]
	//if !ok {
	//	c.cnt++
	//}
	//if c.cnt > c.maxCnt {
	//	c.mutex.Unlock()
	//	return errOverCapacity
	//}
	//c.mutex.Unlock()
	//return c.BuildInMapCache.Set(ctx, key, val, expiration)

	// 法三：
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.data[key]
	if !ok {
		if c.cnt+1 > c.maxCnt {
			// 可在这里设计复杂淘汰策略
			return errOverCapacity
		}
		c.cnt++
	}

	return c.BuildInMapCache.set(ctx, key, val, expiration)
}
