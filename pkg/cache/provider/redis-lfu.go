package cache_provider

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisLFUCache interface {
	RedisCache
	Flush() error
}

type LFUCache struct {
	RedisLFUCache
	client     *CacheClient
	capacity int64
	key        string // e.g. "users"
}

func (lfu *LFUCache) cacheKey() string {
	return lfu.key
}

func (lfu *LFUCache) capacityCheck() {
	capacity, _ := lfu.client.cache.ZCard(context.Background(), lfu.key+"Set").Result()
	if capacity >= lfu.capacity {
		deleteCount := capacity - lfu.capacity + 1
		items := lfu.client.cache.ZPopMin(context.Background(), lfu.key+"Set", deleteCount).Val()

		var keys []string
		for _, item := range items {
			keys = append(keys, item.Member.(string))
		}

		// todo - what should this line do?
		// func (c redis.cmdable) HDel(ctx context.Context, key string, fields ...string) *redis.IntCmd
		lfu.client.cache.HDel(context.Background(), lfu.key, keys...)
	}


}

func (lfu *LFUCache) Get(key string) (string, error) {
	val, err := lfu.client.cache.HGet(context.Background(), lfu.key, key).Result()
	if err != nil {
		return "", err
	}
	go lfu.client.recordCacheOperation(key)
	return val, err
}

func (lfu *LFUCache) Set(key string, value interface{}, exp time.Duration) (string, error) {
	val, err := lfu.client.cache.Set(lfu.client.ctx, key, value, exp).Result()
	if err != nil {
		return "", err
	}
	go lfu.client.recordCacheOperation(key)
	return val, err
}

func (lfu *LFUCache) Del(key string) error {
	err := lfu.client.cache.HDel(context.Background(), lfu.key+"Hash", key).Err()
	if err != nil {
		return err
	}

	err = lfu.client.cache.ZRem(context.Background(), lfu.key+"Set", key).Err()
	if err != nil {
		return err
	}
	return nil
}

func (lfu *LFUCache) Flush() (err error) {
	err = lfu.client.Del(lfu.key).Err()
	if err != nil {
		return
	}

	// err = lfu.client.Del(context.Background(), LFUSortedSetName).Err()
	// return

	return nil
}

func (lfu *LFUCache) GetMapField(key string, mapField string) *redis.StringCmd {
	log.Fatalln("LFUCache.GetMapField() not implemented")
	return nil
}

func (lfu *LFUCache) SetMap(fieldKey string, values map[string]interface{}) *redis.IntCmd {
	log.Fatalln("LFUCache.SetMap() not implemented")
	return nil
}
