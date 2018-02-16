package cache_test

import (
	"friday/cache"
	"friday/cache/memcache"
	"testing"
)

func TestICache(t *testing.T) {
	cacheImplements := []cache.ICache{
		memcache.NewCache(),
	}
	for _, c := range cacheImplements {
		c.Init()
	}
}
