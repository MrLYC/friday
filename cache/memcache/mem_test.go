package memcache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"friday/cache"
	"friday/cache/memcache"
)

func TestMemCache1(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	var err error
	var item memcache.IMappingItem

	c.SetItem("test", &memcache.MappingItem{
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
	item, err = c.GetItem("test")
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

	item, err = c.GetItem("test")
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

	c.SetItem("test", &memcache.MappingItem{
		Value: "1",
	})

	if c.Size() != 1 {
		t.Errorf("size error")
	}
	c.Delete("test")
	if c.Size() != 0 {
		t.Errorf("size error")
	}
	if c.TypeOf("test") != "" {
		t.Errorf("type error")
	}

	item, err := c.GetItem("test")
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
				c.SetItem(k, &memcache.MappingItem{
					Value: i,
				})
				wg.Done()
			}(key, index)

			go func(k string, cnt int) {
				for cnt > 0 {
					item, err := c.GetItem(k)
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
				c.Delete(k)
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

	c.SetItem("string", &memcache.MappingStringItem{})
	c.SetItem("list", &memcache.MappingListItem{})
	c.SetItem("table", &memcache.MappingTableItem{})

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
		value, err := c.Get("test")
		if err != nil {
			t.Errorf("get string error: %v", err)
		}
		c.Set("test", fmt.Sprintf("readed: %v", value))
		wg.Done()
	}()
	c.Set("test", "")
	c.Update("test", func(key string, item interface{}) {
		readyWG.Done()
		item.(*memcache.MappingStringItem).Value = "lyc"
	})
	wg.Wait()
	value, err := c.Get("test")
	if err != nil {
		t.Errorf("get string error: %v", err)
	}
	if value != "readed: lyc" {
		t.Errorf("update error: %v", value)
	}
}

func TestMemCacheStringItem(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	err := c.Set("test", "123")
	if err != nil {
		t.Errorf("set error")
	}

	value, err := c.Get("test")
	if err != nil {
		t.Errorf("get error")
	}
	if value != "123" {
		t.Errorf("get value error")
	}

	value, err = c.Get("nothing")
	if err != cache.ErrItemNotFound {
		t.Errorf("get error")
	}
	if value != "" {
		t.Errorf("get value error")
	}
}

func TestMemCacheStringItemAsNumber(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	var (
		value1 float64
		value2 int64
		err    error
	)

	value2, err = c.Incr("n")
	if value2 != 1 || err != nil {
		t.Errorf("value Incr error: %v", value2)
	}

	value1, err = c.IncrBy("n", 2.0)
	if value1 != 3.0 || err != nil {
		t.Errorf("value IncrBy error: %v", value1)
	}

	value2, err = c.Decr("n")
	if value2 != 2.0 || err != nil {
		t.Errorf("value Decr error: %v", value2)
	}

	value1, err = c.DecrBy("n", 2.0)
	if value1 != 0.0 || err != nil {
		t.Errorf("value DecrBy error: %v", value1)
	}
}

func TestMemCacheListItem1(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	length, err := c.LLen("list")
	if err != cache.ErrItemNotFound {
		t.Errorf("get length error")
	}
	if length != 0 {
		t.Errorf("get length error")
	}

	c.LPush("list", "1")
	c.RPush("list", "2")

	val, err := c.LIndex("list", 0)
	if err != nil {
		t.Errorf("get index error")
	}
	if val != "1" {
		t.Errorf("get index error")
	}

	val, err = c.LIndex("list", 1)
	if err != nil {
		t.Errorf("get index error")
	}
	if val != "2" {
		t.Errorf("get index error")
	}

	length, err = c.LLen("list")
	if err != nil {
		t.Errorf("get length error")
	}
	if length != 2 {
		t.Errorf("get length error")
	}

	val, err = c.LPop("list")
	if err != nil {
		t.Errorf("lpop error")
	}
	if val != "1" {
		t.Errorf("lpop error")
	}

	val, err = c.RPop("list")
	if err != nil {
		t.Errorf("rpop error")
	}
	if val != "2" {
		t.Errorf("rpop error")
	}

	length, err = c.LLen("list")
	if err != nil {
		t.Errorf("get length error")
	}
	if length != 0 {
		t.Errorf("get length error")
	}

	c.Delete("list")
	length, err = c.LLen("list")
	if err != cache.ErrItemNotFound {
		t.Errorf("get length error")
	}
	if length != 0 {
		t.Errorf("get length error")
	}
}

func TestMemCacheListItem2(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	err := c.LSet("list", 0, "0")
	if err != nil {
		t.Errorf("LSet error")
	}
	err = c.LSet("list", 2, "2")
	if err != nil {
		t.Errorf("LSet error")
	}
	values, err := c.LRange("list", 0, -1)
	if err != nil || len(values) != 1 || values[0] != "0" {
		t.Errorf("LSet error： %v", values)
	}

	err = c.LSet("list", 1, "1")
	if err != nil {
		t.Errorf("LSet error")
	}

	values, err = c.LRange("list", 0, -1)
	if err != nil || len(values) != 2 || values[0] != "0" || values[1] != "1" {
		t.Errorf("LSet error： %v", values)
	}
}

func TestMemCacheTableItem(t *testing.T) {
	c := memcache.NewMemCache()
	defer c.Close()

	var (
		value    string
		mappings map[string]string
		err      error
	)

	err = c.HSet("table", "a", "1")
	if err != nil {
		t.Errorf("table HSet value error")
	}

	value, err = c.HGet("table", "a")
	if err != nil {
		t.Errorf("table HGet value error")
	}
	if value != "1" {
		t.Errorf("table HGet value error")
	}

	err = c.HMSet("table", map[string]string{
		"b": "2",
		"c": "3",
	})
	if err != nil {
		t.Errorf("table HMSet value error")
	}

	mappings, err = c.HMGet("table", []string{"a", "b"})
	if err != nil || len(mappings) != 2 || mappings["a"] != "1" || mappings["b"] != "2" {
		t.Errorf("table HMGet value error: %v", mappings)
	}

	if !c.HExists("table", "c") {
		t.Errorf("HExists error")
	}
	err = c.HDel("table", "c")
	if err != nil {
		t.Errorf("HDel error")
	}
	if c.HExists("table", "c") {
		t.Errorf("HExists error")
	}

	c.HClear("table")
	if c.HExists("table", "a") {
		t.Errorf("HClear error")
	}
	if c.HExists("table", "b") {
		t.Errorf("HClear error")
	}

	if !c.Exists("table") {
		t.Errorf("HClear error")
	}
}
