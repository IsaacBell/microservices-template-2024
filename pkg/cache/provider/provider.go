package cache_provider

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type CacheClient struct {
	RedisCache
	ctx   context.Context
	cache *redis.Client
}

func NewCache(ctx context.Context) RedisCache {
	return UseRedis(ctx)
}
