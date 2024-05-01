package cache_provider

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	ctx   context.Context
	cache *redis.Client
}

func Cache(ctx context.Context) *CacheClient {
	return UseRedis(ctx)
}

func CacheProvider(ctx context.Context) *CacheClient {
	return Cache(ctx)
}
