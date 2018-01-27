package memcache_test

import (
	"testing"

	"friday/storage/cache/memcache"
)

func TestMappingListItem1(t *testing.T) {
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

func TestMappingListItem2(t *testing.T) {
	var item memcache.MappingListItem
	if item.GetString(0) != "" {
		t.Errorf("value error")
	}
	if item.GetFirstString() != "" {
		t.Errorf("value error")
	}
	if item.GetLastString() != "" {
		t.Errorf("value error")
	}
	if item.PopFirstString() != "" {
		t.Errorf("pop first value error")
	}
	if item.PopLastString() != "" {
		t.Errorf("pop last value error")
	}

	item.Init()

	if item.GetString(0) != "" {
		t.Errorf("value error")
	}
	if item.GetFirstString() != "" {
		t.Errorf("value error")
	}
	if item.GetLastString() != "" {
		t.Errorf("value error")
	}
	if item.PopFirstString() != "" {
		t.Errorf("pop first value error")
	}
	if item.PopLastString() != "" {
		t.Errorf("pop last value error")
	}

	item.AppendFirstString("2")
	item.AppendFirstString("1")
	item.AppendLastString("3")

	if item.GetString(0) != "1" {
		t.Errorf("value error")
	}
	if item.GetFirstString() != "1" {
		t.Errorf("value error")
	}
	if item.GetString(1) != "2" {
		t.Errorf("value error")
	}
	if item.GetString(2) != "3" {
		t.Errorf("value error")
	}
	if item.GetString(3) != "" {
		t.Errorf("value error")
	}

	if item.PopFirstString() != "1" {
		t.Errorf("pop first value error")
	}
	if item.PopLastString() != "3" {
		t.Errorf("pop last value error")
	}
}
