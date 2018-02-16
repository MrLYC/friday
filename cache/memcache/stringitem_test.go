package memcache_test

import (
	"friday/cache"
	"friday/cache/memcache"
	"testing"
)

func TestMappingStringItem(t *testing.T) {
	var item memcache.MappingStringItem
	if item.Length() != 0 {
		t.Errorf("length error")
	}
	if item.GetString() != "" {
		t.Errorf("value error")
	}
	item.SetValue("test")
	if item.Length() != 4 {
		t.Errorf("length error")
	}
}

func TestMappingStringItem1(t *testing.T) {
	var (
		item  memcache.MappingStringItem
		value float64
		err   error
	)

	value, err = item.Add(1.0)
	if err != nil || value != 1.0 {
		t.Errorf("value IncrBy error: %v", value)
	}

	value, err = item.Add(2.0)
	if err != nil || value != 3.0 {
		t.Errorf("value IncrBy error: %v", value)
	}

	item.SetValue("test")
	value, err = item.Add(3.0)
	if err != cache.ErrItemValueError {
		t.Errorf("value not a number")
	}
}
