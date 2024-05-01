package cache

import (
	cache_provider "microservices-template-2024/pkg/cache/provider"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(cache_provider.CacheProvider)
