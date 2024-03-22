package cache

import (
	"context"
	"errors"
	"fmt"
	redis "github.com/redis/go-redis/v9"
	//_ "github.com/golang/mock/mockgen/model"
	"time"
)

var (
	errFailToSetCache = errors.New("cache: 写入redis失败")
)

type RedisConfig struct {
	Addr string
}

type RedisCache struct {
	client redis.Cmdable
}

func (r *RedisCache) LoadAndDelete(ctx context.Context, key string) (any, error) {
	return r.client.GetDel(ctx, key).Result()
}

func NewRedisCache(client redis.Cmdable) *RedisCache {
	return &RedisCache{client: client}
}

func (r *RedisCache) Set(ctx context.Context, key string, val any, expiration time.Duration) error {
	res, err := r.client.Set(ctx, key, val, expiration).Result()
	if err != nil {
		return err
	}
	if res != "OK" {
		return fmt.Errorf("%w ,返回信息 %s", errFailToSetCache, res)
	}
	return nil
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	return err
}
