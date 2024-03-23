package cache

import (
	"context"
	"fmt"
	"golang.org/x/sync/singleflight"
	"time"
)

type SingleflightCacheV1 struct {
	ReadThroughCache
}

func NewSingleflightCacheV1(cache Cache, loadFunc func(ctx context.Context, key string) (any, error),
	expiration time.Duration) *SingleflightCacheV1 {
	g := &singleflight.Group{}
	return &SingleflightCacheV1{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				val, err, _ := g.Do(key, func() (interface{}, error) {
					return loadFunc(ctx, key)
				})
				return val, err
			},
			Expiration: expiration,
		},
	}
}

type SingleflightCacheV2 struct {
	ReadThroughCache
	g *singleflight.Group
}

func (r *SingleflightCacheV2) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound {
		val, err, _ = r.g.Do(key, func() (interface{}, error) {
			val, err := r.LoadFunc(ctx, key)
			if err == nil {
				er := r.Cache.Set(ctx, key, val, r.Expiration)
				if er != nil {
					return val, fmt.Errorf("%s，原因：%s", errFailToSetCache, er.Error())
				}
			}
			return val, err
		})

	}
	return val, err
}
