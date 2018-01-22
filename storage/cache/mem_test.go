package cache_test

import (
	"fmt"
	"sync"
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

	if c.Size() != 1 {
		t.Errorf("size error")
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

	if c.Size() != 0 {
		t.Errorf("size error")
	}
}

func TestMemCache2(t *testing.T) {
	c := cache.NewMemCache()
	defer c.Close()

	c.Set("test", &cache.MappingItem{
		Value: "1",
	})

	if c.Size() != 1 {
		t.Errorf("size error")
	}
	c.Remove("test")
	if c.Size() != 0 {
		t.Errorf("size error")
	}

	item, err := c.Get("test")
	if err != cache.ErrItemNotFound {
		t.Errorf("get item error")
	}
	if item != nil {
		t.Errorf("return item error")
	}
}

func TestMemCacheConcurrency(t *testing.T) {
	c := cache.NewMemCache()
	defer c.Close()
	var (
		count = 100
		flag  = make([]int, count)
		wg    sync.WaitGroup
	)
	wg.Add(1)
	go func() {
		for index := 0; index < count; index++ {
			wg.Add(2)
			key := fmt.Sprintf("%v", index)
			go func(k string, i int) {
				c.Set(k, &cache.MappingItem{
					Value: i,
				})
				wg.Done()
			}(key, index)

			go func(k string, cnt int) {
				for cnt > 0 {
					item, err := c.Get(k)
					if err == cache.ErrItemNotFound {
						continue
					} else if err != nil {
						t.Errorf("Get error: %v", err)
						break
					}
					index := item.GetValue().(int)
					if k != fmt.Sprintf("%v", index) {
						t.Errorf("value error: %v-%v", k, index)
						break
					}
					cnt -= 1
					flag[index] += 1
				}
				c.Remove(k)
				wg.Done()
			}(key, count)
		}
		wg.Done()
	}()

	wg.Wait()

	if c.Size() != 0 {
		t.Errorf("size error")
	}
	for i, v := range flag {
		if v != count {
			t.Errorf("%v error: %v", i, v)
		}
	}
}
