package cache

import (
	"context"
	"fmt"
	"time"
)

type BloomFilterCache struct {
	ReadThroughCache
}

func NewBloomFilterCache(cache Cache, bf BloomFilter, loadFunc func(ctx context.Context, key string) (any, error),
	expiration time.Duration) *BloomFilterCache {
	return &BloomFilterCache{
		ReadThroughCache: ReadThroughCache{
			Cache: cache,
			LoadFunc: func(ctx context.Context, key string) (any, error) {
				if !bf.HasKey(ctx, key) {
					return nil, errKeyNotFound
				}
				return loadFunc(ctx, key)
			},
			Expiration: expiration,
		},
	}
}

type BloomFilterCacheV1 struct {
	ReadThroughCache
	bf BloomFilter
}

func (r *BloomFilterCacheV1) Get(ctx context.Context, key string) (any, error) {
	val, err := r.Cache.Get(ctx, key)
	if err == errKeyNotFound && r.bf.HasKey(ctx, key) {
		val, err = r.LoadFunc(ctx, key)
		if err == nil {
			er := r.Cache.Set(ctx, key, val, r.Expiration)
			if er != nil {
				return val, fmt.Errorf("%s，原因：%s", errFailToSetCache, er.Error())
			}
		}
	}
	return val, err
}

type BloomFilter interface {
	HasKey(ctx context.Context, key string) bool
}
