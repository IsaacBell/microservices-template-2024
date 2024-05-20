package cache

import (
	"context"
	"sync"

	cache_provider "microservices-template-2024/pkg/cache/provider"

	"github.com/google/wire"
)

var (
	once sync.Once
	cache *cache_provider.CacheClient = nil
)

func Cache(ctx context.Context) *cache_provider.CacheClient {
	once.Do(func () { cache = cache_provider.NewCache(ctx) })
	return cache
}

var ProviderSet = wire.NewSet(Cache)
