package cache

import (
	"context"
	"math/rand"
	"time"
)

type RandomExpirationCache struct {
	Cache
}

func (r *RandomExpirationCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	if expiration > 0 {
		// 加上一个0-300s的偏移量
		offset := time.Duration(rand.Intn(300)) * time.Second
		expiration = expiration + offset
		return r.Cache.Set(ctx, key, val, expiration)
	}
	return r.Cache.Set(ctx, key, val, expiration)
}
