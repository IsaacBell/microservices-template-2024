package cache

import (
	"context"

	cache_provider "microservices-template-2024/pkg/cache/provider"

	"github.com/google/wire"
)

func Cache(ctx context.Context) *cache_provider.CacheClient {
	return cache_provider.UseRedis(ctx)
}

func CacheProvider(ctx context.Context) *cache_provider.CacheClient {
	return cache_provider.CacheProvider(ctx)
}

var ProviderSet = wire.NewSet(CacheProvider)
