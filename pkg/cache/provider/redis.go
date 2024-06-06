package cache_provider

import (
	"context"
	"core/internal/conf"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	onceCacheClient sync.Once
	cacheInstance   RedisCache
)

type RedisCache interface {
	UseContext(context.Context)
	CurrentContext() context.Context
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, exp time.Duration) *redis.StatusCmd
	Del(key string) *redis.IntCmd
	GetMapField(key string, mapField string) *redis.StringCmd
	SetMap(fieldKey string, values map[string]interface{}) *redis.IntCmd
}

func UseRedis(ctx context.Context) RedisCache {
	onceCacheClient.Do(func() {
		cache := conf.RedisConn(ctx)
		cacheInstance = &CacheClient{ctx: ctx, cache: cache}
	})
	return cacheInstance
}

func ChangeRedisClientTo(r RedisCache) {
	cacheInstance = r
}

func (cache *CacheClient) UseContext(ctx context.Context) {
	cache.ctx = ctx
}

func (cache *CacheClient) CurrentContext() context.Context {
	return cache.ctx
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
