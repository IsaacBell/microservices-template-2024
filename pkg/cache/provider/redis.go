package cache_provider

import (
	"context"
	"fmt"
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
	GetMapField(ctx context.Context, key string, mapField string) *redis.StringCmd
	SetMap(ctx context.Context, fieldKey string, values map[string]interface{}) *redis.IntCmd
}

func UseRedis(ctx context.Context) *CacheClient {
	onceCacheClient.Do(func() {
		cache := conf.RedisConn(ctx)
		cacheInstance = &CacheClient{ctx: ctx, cache: cache}
	})
	return cacheInstance
}

func (cache *CacheClient) UseContext(ctx context.Context) {
	cache.ctx = ctx
}

func (cache *CacheClient) Get(key string) *redis.StringCmd {
	return cache.cache.Get(cache.ctx, key)
}

func (cache *CacheClient) Set(key string, value interface{}, exp time.Duration) *redis.StatusCmd {
	return cache.cache.Set(cache.ctx, key, value, exp)
}

func (cache *CacheClient) Del(key string) *redis.IntCmd {
	return cache.cache.Del(cache.ctx, key)
}

func (cache *CacheClient) GetMapField(key string, mapField string) *redis.StringCmd {
	return cache.cache.HGet(cache.ctx, key, mapField)
}

func (cache *CacheClient) SetMap(fieldKey string, values map[string]interface{}) *redis.IntCmd {
	var out *redis.IntCmd
	for k, v := range values {
		out = cache.cache.HSet(cache.ctx, fieldKey, k, v)
		if out.Err() != nil {
			fmt.Println("error caching map: ", out.Err())
		}
	}

	return out
}
