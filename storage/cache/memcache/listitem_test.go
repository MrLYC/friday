package memcache_test

import (
	"testing"

	"friday/storage/cache/memcache"
)

func TestMappingListItem(t *testing.T) {
	var item memcache.MappingListItem

	if item.Length() != 0 {
		t.Errorf("length error")
	}

	item.Init()

	list := item.GetList()
	if list == nil {
		t.Errorf("list init failed")
	}
	if item.Length() != 0 {
		t.Errorf("length error")
	}
	if item.GetFirstString() != "" {
		t.Errorf("first value error")
	}
	if item.GetLastString() != "" {
		t.Errorf("last value error")
	}

	list.Add("1")
	if item.Length() != 1 {
		t.Errorf("length error")
	}
	list.Add("2")
	if item.Length() != 2 {
		t.Errorf("length error")
	}

	if item.GetFirstString() != "1" {
		t.Errorf("first value error")
	}
	if item.GetLastString() != "2" {
		t.Errorf("last value error")
	}
}
