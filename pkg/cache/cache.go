package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	cache_provider "core/pkg/cache/provider"
	"core/pkg/stream"

	"github.com/google/wire"
)

var (
	once  sync.Once
	cache cache_provider.RedisCache
)

func setCache() {
	if cache == nil {
		cache = Cache(context.Background())
		if cache == nil { // if cache nil after setting it, assume redis is down
			fmt.Println("---ISSUE: Can't retrieve Redis client---")
			AlertRedisConnectionError()
		}
	}
}

func AlertRedisConnectionError() {}

func Cache(ctx context.Context) cache_provider.RedisCache {
	once.Do(func() { cache = cache_provider.NewCache(ctx) })
	if cache == nil {
		log.Fatalln("!!!Can't initialize Redis client!!!")
	} else {
		fmt.Println("Redis running in ctx: ", cache.CurrentContext())
	}
	return cache
}

func GetUserSession(id string) (string, error) {
	setCache()
	return cache.Get("auth:" + id).Result()
}

func CacheRecord(recordType, cacheKey, id string, data interface{}) {
	setCache()
	go func() {
		if data == nil || cacheKey == "" || recordType == "" || id == "" {
			return
		}
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
