package cache

import (
	"context"
	"sync"

	cache_provider "microservices-template-2024/pkg/cache/provider"

	"github.com/google/wire"
)

var (
	once  sync.Once
	cache *cache_provider.CacheClient = nil
)

func Cache() *cache_provider.CacheClient {
	once.Do(func() { cache = cache_provider.NewCache(context.Background()) })
	return cache
}

func LFUCache(cacheKey string, capacity int64) *cache_provider.LFUCache {
	return cache_provider.NewLFUCache(context.Background(), cacheKey, capacity)
}

var ProviderSet = wire.NewSet(Cache, LFUCache)
