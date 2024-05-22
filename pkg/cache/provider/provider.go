package cache_provider

import (
	"context"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	ctx   context.Context
	cache *redis.Client
}

func NewCache(ctx context.Context) *CacheClient {
	return UseRedis(ctx)
}

func NewLFUCache(ctx context.Context, key string, capacity int64) *LFUCache {
	// todo - add lfu caches to an internal registry
	return &LFUCache{client: UseRedis(ctx), capacity: capacity, key: key}
}

var ProviderSet = wire.NewSet(NewCache, UseRedis, NewLFUCache)
