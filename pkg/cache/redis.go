package cache

import (
	"context"
	"microservices-template-2024/internal/conf"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	onceCacheClient sync.Once
	cacheInstance   *CacheClient
)

type RedisCache interface {
	Get(ctx context.Context, key string) redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, exp time.Duration)
}

type CacheClient struct {
	ctx   context.Context
	cache *redis.Client
}

func Cache(ctx context.Context) *CacheClient {
	onceCacheClient.Do(func() {
		cache := conf.RedisConn(ctx)
		cacheInstance = &CacheClient{ctx: ctx, cache: cache}
	})
	return cacheInstance
}

func CacheProvider(ctx context.Context) *CacheClient {
	return Cache(ctx)
}

func (cache *CacheClient) UseContext(ctx context.Context) {
	cache.ctx = ctx
}

func (cache *CacheClient) Get(ctx context.Context, key string) *redis.StringCmd {
	return cache.cache.Get(ctx, key)
}

func (cache *CacheClient) Set(ctx context.Context, key string, value interface{}, exp time.Duration) *redis.StatusCmd {
	return cache.cache.Set(ctx, key, value, exp)
}
