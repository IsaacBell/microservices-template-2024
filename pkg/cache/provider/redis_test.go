package cache_provider_test

import (
	"context"
	"fmt"
	cache_provider "microservices-template-2024/pkg/cache/provider"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisCache struct {
	mock.Mock
	cache_provider.RedisCache
}

func (m *MockRedisCache) Get(key string) *redis.StringCmd {
	fmt.Println("key", key)
	return redis.NewStringCmd(context.Background(), key)
}

func (m *MockRedisCache) Set(key string, value interface{}, exp time.Duration) *redis.StatusCmd {
	args := m.Called(key, value, exp)
	return args.Get(0).(*redis.StatusCmd)
}

func (m *MockRedisCache) Del(key string) *redis.IntCmd {
	args := m.Called(key)
	return args.Get(0).(*redis.IntCmd)
}

func (m *MockRedisCache) GetMapField(key string, mapField string) *redis.StringCmd {
	return redis.NewStringCmd(m.CurrentContext(), key)
}

func (m *MockRedisCache) SetMap(fieldKey string, values map[string]interface{}) *redis.IntCmd {
	args := m.Called(fieldKey, values)
	return args.Get(0).(*redis.IntCmd)
}

func testUseRedis(t *testing.T) (MockRedisCache, cache_provider.RedisCache) {
	mockRedis := new(MockRedisCache)
	cacheClient := cache_provider.NewCache(context.Background())
	cache_provider.ChangeRedisClientTo(mockRedis)
	cacheClient = cache_provider.NewCache(context.Background())
	assert.NotNil(t, mockRedis)
	assert.NotNil(t, cacheClient)
	fmt.Println("client: ", cacheClient)

	return *mockRedis, cacheClient
}

func TestCacheClient_Get(t *testing.T) {
	key := "test_key"
	expected := redis.NewStringCmd(context.Background(), key)

	mockRedis, cacheClient := testUseRedis(t)

	result := cacheClient.Get(key)
	assert.Equal(t, expected, result)
	mockRedis.AssertExpectations(t)
}

func TestCacheClient_Set(t *testing.T) {
	mockRedis := new(MockRedisCache)
	cacheClient := cache_provider.NewCache(context.Background())
	cache_provider.ChangeRedisClientTo(mockRedis)
	cacheClient = cache_provider.NewCache(context.Background())
	fmt.Println("client: ", cacheClient)

	key := "test_key"
	value := "test_value"
	exp := time.Duration(60) * time.Second
	mockRedis.On("Set", key, value, exp).Return(redis.NewStatusResult("", nil))

	result := cacheClient.Set(key, value, exp)
	fmt.Println("res.nil?: ", result, " - ", result == nil)
	assert.NotNil(t, result)
	mockRedis.AssertExpectations(t)
}
