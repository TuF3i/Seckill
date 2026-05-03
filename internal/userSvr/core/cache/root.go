package cache

import "github.com/redis/go-redis/v9"

type CacheReliance struct {
	rdb *redis.Client
}

type Cache struct {
	*CacheReliance
}

func NewCache(m *CacheReliance) *Cache {
	return &Cache{m}
}
