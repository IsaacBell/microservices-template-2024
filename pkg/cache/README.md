# Cache Package

The cache package provides a simple caching mechanism using Redis. It allows you to store and retrieve values from the cache using a key-value approach.

## Usage

To use the cache package, import it in your Go code:

```go
import "microservices-template-2024/pkg/cache"
```

### Cache Initialization

The cache package uses a singleton instance of the CacheClient to manage the cache. The Cache function lazily initializes the cache client using the provided context. It ensures that only one instance of the cache client is created throughout the application.

```go
cacheClient := cache.Cache(ctx)
```

### Storing Values
To store a value in the cache, use the Set method of the CacheClient:

```go
err := cacheClient.Set("key", "value", time.Hour).Err()
if err != nil {
    // Handle error
}
```

The Set method takes three parameters:

1. key: The key associated with the value.
2. value: The value to be stored in the cache.
3. expiration: The expiration time for the cached value (optional).

### Retrieving Values

To retrieve a value from the cache, use the Get method of the CacheClient:

```go
value, err := cacheClient.Get("key").Result()
if err != nil {
    // Handle error
}
```

The Get method takes the key as a parameter and returns a StringCmd which allows you to access the cached value as a string.

### Deleting Values
To delete a value from the cache, use the Del method of the CacheClient:

```go
err := cacheClient.Del("key").Err()
if err != nil {
    // Handle error
}
```

The Del method takes the key as a parameter and removes the associated value from the cache.

### NewCache

The NewCache function in the cache_provider package creates a new instance of the CacheClient using the provided context.

```go
cacheClient := cache_provider.NewCache(ctx)
```

### CacheClient

The CacheClient struct represents the Redis client and provides methods for interacting with the cache.

```go
type CacheClient struct {
    // ...
}
```

The CacheClient struct implements the RedisCache interface, which defines the methods for cache operations.
RedisCache Interface
The RedisCache interface defines the methods for interacting with the cache using Redis.

```go
type RedisCache interface {
    Get(ctx context.Context, key string) *redis.StringCmd
    Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
    Del(ctx context.Context, key string) *redis.IntCmd
    // ...
}
```

## Dependencies

The cache package depends on the following packages:

- context: Provides context support for cache operations.
- sync: Provides synchronization primitives for lazy initialization of the cache client.
- time: Provides time-related functions for cache expiration.
- github.com/redis/go-redis/v9: Provides the Redis client for cache operations.

### Example Usage

```go
import (
    "context"
    "fmt"
    "microservices-template-2024/pkg/cache"
    "time"
)

func main() {
    ctx := context.Background()
    cacheClient := cache.Cache(ctx)

    // Store a value in the cache
    err := cacheClient.Set("my_key", "my_value", time.Hour).Err()
    if err != nil {
        fmt.Println("Error storing value in cache:", err)
        return
    }

    // Retrieve a value from the cache
    value, err := cacheClient.Get("my_key").Result()
    if err != nil {
        fmt.Println("Error retrieving value from cache:", err)
        return
    }
    fmt.Println("Cached value:", value)

    // Delete a value from the cache
    err = cacheClient.Del("my_key").Err()
    if err != nil {
        fmt.Println("Error deleting value from cache:", err)
        return
    }
    fmt.Println("Value deleted from cache")
}
This example demonstrates how to store a value in the cache, retrieve it, and then delete it using the cache package.