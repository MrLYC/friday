package memcache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"friday/storage/cache"
	"friday/storage/cache/memcache"
)

func TestMemCache1(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	var err error
	var item memcache.IMappingItem

	c.Set("test", &memcache.MappingItem{
		Value: "1",
	})

	if !c.Exists("test") {
		t.Errorf("item not exists")
	}

	if c.Size() != 1 {
		t.Errorf("size error")
	}

	t1, _ := time.ParseDuration("10h")
	err = c.Expire("test", t1)
	if err != nil {
		t.Errorf("expire error")
	}
	item, err = c.Get("test")
	if err != nil {
		t.Errorf("get item error")
	}
	if item.GetValue().(string) != "1" {
		t.Errorf("get value failed")
	}
	if !item.IsAvailable() {
		t.Errorf("item not available")
	}

	t2, _ := time.ParseDuration("-10h")
	err = c.Expire("test", t2)
	if err != nil {
		t.Errorf("expire error")
	}
	if item.IsAvailable() {
		t.Errorf("item error")
	}

	item, err = c.GetRaw("test")
	if err != nil {
		t.Errorf("get item error")
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
	c := memcache.NewMemCache()
	defer c.Close()

	c.Set("test", &memcache.MappingItem{
		Value: "1",
	})

	if c.Size() != 1 {
		t.Errorf("size error")
	}
	c.Remove("test")
	if c.Size() != 0 {
		t.Errorf("size error")
	}
	if c.TypeOf("test") != "" {
		t.Errorf("type error")
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
	c := memcache.NewMemCache()
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
				c.Set(k, &memcache.MappingItem{
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

func TestTestMemCache3(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	c.Set("string", &memcache.MappingStringItem{})
	c.Set("list", &memcache.MappingListItem{})
	c.Set("table", &memcache.MappingTableItem{})

	var (
		keys    = []string{"string", "list", "table"}
		results = [][]bool{
			[]bool{true, false, false},
			[]bool{false, true, false},
			[]bool{false, false, true},
		}
	)

	for i, key := range keys {
		item, err := c.GetStringItem(key)
		result := results[i]
		if result[0] {
			if item == nil {
				t.Errorf("item[string] value error")
			}
			if err != nil {
				t.Errorf("item error: %v", err)
			}
			if c.TypeOf(key) != key {
				t.Errorf("item type error: string")
			}
		} else {
			if item != nil {
				t.Errorf("item[string] value error")
			}
			if err != cache.ErrItemTypeError {
				t.Errorf("item error: %v", err)
			}
			if !c.Exists(key) {
				t.Errorf("exists error: string")
			}
		}
	}

	for i, key := range keys {
		item, err := c.GetListItem(key)
		result := results[i]
		if result[1] {
			if item == nil {
				t.Errorf("item[list] value error")
			}
			if err != nil {
				t.Errorf("item error: %v", err)
			}
			if c.TypeOf(key) != key {
				t.Errorf("item type error: list")
			}
		} else {
			if item != nil {
				t.Errorf("item[list] value error")
			}
			if err != cache.ErrItemTypeError {
				t.Errorf("item error: %v", err)
			}
			if !c.Exists(key) {
				t.Errorf("exists error: list")
			}
		}
	}

	for i, key := range keys {
		item, err := c.GetTableItem(key)
		result := results[i]
		if result[2] {
			if item == nil {
				t.Errorf("item[table] value error")
			}
			if err != nil {
				t.Errorf("item error: %v", err)
			}
			if c.TypeOf(key) != key {
				t.Errorf("item type error: table")
			}
		} else {
			if item != nil {
				t.Errorf("item[table] value error")
			}
			if err != cache.ErrItemTypeError {
				t.Errorf("item error: %v", err)
			}
			if !c.Exists(key) {
				t.Errorf("exists error: table")
			}
		}
	}
}

func TestMemCacheUpdate(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	var readyWG sync.WaitGroup
	var wg sync.WaitGroup

	wg.Add(1)
	readyWG.Add(1)
	go func() {
		readyWG.Wait()
		value, err := c.GetString("test")
		if err != nil {
			t.Errorf("get string error: %v", err)
		}
		c.SetString("test", fmt.Sprintf("readed: %v", value))
		wg.Done()
	}()
	c.SetString("test", "")
	c.Update("test", func(key string, item interface{}) {
		readyWG.Done()
		item.(*memcache.MappingStringItem).Value = "lyc"
	})
	wg.Wait()
	value, err := c.GetString("test")
	if err != nil {
		t.Errorf("get string error: %v", err)
	}
	if value != "readed: lyc" {
		t.Errorf("update error: %v", value)
	}
}

func TestMemCacheSetGetString(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	err := c.SetString("test", "123")
	if err != nil {
		t.Errorf("set error")
	}

	value, err := c.GetString("test")
	if err != nil {
		t.Errorf("get error")
	}
	if value != "123" {
		t.Errorf("get value error")
	}

	value, err = c.GetString("nothing")
	if err != cache.ErrItemNotFound {
		t.Errorf("get error")
	}
	if value != "" {
		t.Errorf("get value error")
	}
}
