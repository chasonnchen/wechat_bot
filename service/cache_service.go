package service

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

var (
	cacheServiceInstance = &CacheService{}
)

func NewCacheService() *CacheService {
	return cacheServiceInstance
}

type CacheService struct {
	Cache *gocache.Cache
}

func (c *CacheService) init() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	c.Cache = gocache.New(1*time.Minute, 5*time.Minute)
}

func (c *CacheService) Set(k string, v interface{}, d time.Duration) {
	c.Cache.Set(k, v, d)
}

func (c *CacheService) Get(k string) (interface{}, bool) {
	return c.Cache.Get(k)
}
