package cache_test

import (
	"testing"
	"time"

	"friday/storage/cache"
)

func TestMappingItem(t *testing.T) {
	item := &cache.MappingItem{
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

func TestMemCache1(t *testing.T) {
	c := cache.NewMemCache()
	defer c.Close()

	c.Set("test", &cache.MappingItem{
		Value: "1",
	})

	if !c.Exists("test") {
		t.Errorf("item not exists")
	}

	item, err := c.Get("test")
	if err != nil {
		t.Errorf("get item error")
	}
	if item.GetValue().(string) != "1" {
		t.Errorf("get value failed")
	}
	if !item.IsAvailable() {
		t.Errorf("item not available")
	}

	tt, _ := time.ParseDuration("-10h")
	if !item.IsAvailable() {
		t.Errorf("item error")
	}
	err = c.Expire("test", tt)
	if err != nil {
		t.Errorf("expire error")
	}
	if item.IsAvailable() {
		t.Errorf("item error")
	}

	item, err = c.Get("test")
	if err != cache.ErrItemNotFound {
		t.Errorf("get item error")
	}
	if item != nil {
		t.Errorf("return item error")
	}

	if c.Clean() != 1 {
		t.Errorf("clean error")
	}
}

func TestMemCache2(t *testing.T) {
	c := cache.NewMemCache()
	defer c.Close()

	c.Set("test", &cache.MappingItem{
		Value: "1",
	})

	c.Remove("test")

	item, err := c.Get("test")
	if err != cache.ErrItemNotFound {
		t.Errorf("get item error")
	}
	if item != nil {
		t.Errorf("return item error")
	}
}
