package memcache_test

import (
	"testing"

	"friday/storage/cache/memcache"
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
