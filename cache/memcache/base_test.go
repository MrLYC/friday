package memcache_test

import (
	"testing"
	"time"

	"friday/cache/memcache"
)

func TestMappingItem(t *testing.T) {
	item := &memcache.MappingItem{
		Value: "1",
	}
	now := time.Now()
	if !item.IsAvailable() {
		t.Errorf("item not available")
	}
	if item.GetValue().(string) != "1" {
		t.Errorf("item value error")
	}

	t1, _ := time.ParseDuration("10h")
	item.SetExpireAt(now.Add(t1))
	if !item.IsAvailable() {
		t.Errorf("item not available")
	}
	if item.IsExpireAt(now) {
		t.Errorf("item expires error")
	}

	t2, _ := time.ParseDuration("-10h")
	item.SetExpireAt(now.Add(t2))
	if item.IsAvailable() {
		t.Errorf("item not available")
	}
	if !item.IsExpireAt(now) {
		t.Errorf("item expires error")
	}
}
