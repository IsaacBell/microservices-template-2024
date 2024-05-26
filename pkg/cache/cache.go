package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	cache_provider "microservices-template-2024/pkg/cache/provider"
	"microservices-template-2024/pkg/stream"

	"github.com/google/wire"
)

var (
	once  sync.Once
	cache *cache_provider.CacheClient = nil
)

func Cache(ctx context.Context) *cache_provider.CacheClient {
	once.Do(func() { cache = cache_provider.NewCache(ctx) })
	return cache
}

func CacheRecord(recordType, cacheKey, id string, data interface{}) {
	go func() {
		cacheData, _ := json.Marshal(data)
		cache.Set(cacheKey, cacheData, time.Hour*2)
		stream.ProduceKafkaMessage("channel/"+recordType, "Cached "+recordType+": "+id)
	}()
}

func Delete(ctx context.Context, cacheKey, recordType string) {
	go func() {
		err := Cache(ctx).Del(cacheKey).Err()
		if err != nil {
			fmt.Printf("Failed to delete cache entry %s: %v\n", cacheKey, err)
		}
		stream.ProduceKafkaMessage("channel/"+recordType, "Deleted from cache: "+cacheKey)
	}()
}

var ProviderSet = wire.NewSet(Cache)
