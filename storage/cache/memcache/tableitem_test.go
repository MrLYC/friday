package memcache_test

import (
	"testing"

	"friday/storage/cache/memcache"
)

func TestMappingTableItem1(t *testing.T) {
	var item memcache.MappingTableItem
	item.Init()

	if item.SetString("a", "1") != nil {
		t.Errorf("set string failed")
	}
	if item.SetString("b", "2") != nil {
		t.Errorf("set string failed")
	}
	if item.Length() != 2 {
		t.Errorf("mapping length error")
	}

	if !item.Exists("a") {
		t.Errorf("exists error")
	}
	if !item.Exists("b") {
		t.Errorf("exists error")
	}
	if item.Exists("c") {
		t.Errorf("exists error")
	}

	if item.GetString("a") != "1" {
		t.Errorf("get string failed")
	}
	if item.GetString("b") != "2" {
		t.Errorf("get string failed")
	}
	if item.GetString("c") != "" {
		t.Errorf("get string failed")
	}

	if item.Delete("b") != nil {
		t.Errorf("delete failed")
	}
	if item.Length() != 1 {
		t.Errorf("mapping length error")
	}
	if item.Clear() != nil {
		t.Errorf("clear failed")
	}
	if item.Length() != 0 {
		t.Errorf("mapping length error")
	}
}

func TestMappingTableItem2(t *testing.T) {
	var item memcache.MappingTableItem
	item.Init()

	if item.SetMappings(map[string]string{
		"a": "1",
		"b": "2",
	}) != nil {
		t.Errorf("set mappings failed")
	}

	values := item.GetAllMappings()
	if item.Length() != len(values) {
		t.Errorf("mapping length error")
	}
	if values["a"] != "1" || values["b"] != "2" {
		t.Errorf("get mappings error")
	}

	if item.SetMappings(map[string]string{
		"b": "0",
		"c": "3",
	}) != nil {
		t.Errorf("set mappings failed")
	}

	values = item.GetAllMappings()
	if item.Length() != len(values) {
		t.Errorf("mapping length error")
	}
	if values["a"] != "1" || values["b"] != "0" || values["c"] != "3" {
		t.Errorf("get mappings error")
	}

	values = item.GetMappings([]string{"a", "b", "x"})
	if len(values) != 3 {
		t.Errorf("mapping length error")
	}
	if values["a"] != "1" || values["b"] != "0" || values["x"] != "" {
		t.Errorf("get mappings error")
	}
}
