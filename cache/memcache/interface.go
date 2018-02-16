package memcache

import (
	"time"
)

// Cache :
type Cache struct {
	*MemCache
}

// Init :
func (i *Cache) Init() {
	i.MemCache.Init()
}

// Close :
func (i *Cache) Close() error {
	return i.MemCache.Close()
}

// StringSet :
func (i *Cache) StringSet(key string, value string) error {
	return i.Set(key, value)
}

// StringGet :
func (i *Cache) StringGet(key string) (string, error) {
	return i.Get(key)
}

// Delete :
func (i *Cache) Delete(key string) {
	i.MemCache.Delete(key)
}

// Expire :
func (i *Cache) Expire(key string, duration time.Duration) {
	i.MemCache.Expire(key, duration)
}

// ListPush :
func (i *Cache) ListPush(key string, value string) error {
	return i.LPush(key, value)
}

// ListPop :
func (i *Cache) ListPop(key string) (string, error) {
	return i.LPop(key)
}

// ListRPush :
func (i *Cache) ListRPush(key string, value string) error {
	return i.RPush(key, value)
}

// ListRPop :
func (i *Cache) ListRPop(key string) (string, error) {
	return i.RPop(key)
}

// ListLen :
func (i *Cache) ListLen(key string) (int, error) {
	return i.LLen(key)
}

// TableSet :
func (i *Cache) TableSet(key string, field string, value string) error {
	return i.HSet(key, field, value)
}

// TableGet :
func (i *Cache) TableGet(key string, field string) (string, error) {
	return i.HGet(key, field)
}

// TableGetAll :
func (i *Cache) TableGetAll(key string) (map[string]string, error) {
	return i.HGetAll(key)
}

// NewCache :
func NewCache() *Cache {
	return &Cache{
		MemCache: &MemCache{},
	}
}
