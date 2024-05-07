package cache_provider

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	ctx   context.Context
	cache *redis.Client
}

func NewCache(ctx context.Context) *CacheClient {
	return UseRedis(ctx)
}

func NewCacheProvider(ctx context.Context) *CacheClient {
	return NewCache(ctx)
}
