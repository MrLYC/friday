package memcache_test

import (
	"testing"

	"friday/cache/memcache"
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

	if item.Delete(1) != nil {
		t.Errorf("delete error")
	}
	if item.PopFirstString() != "1" {
		t.Errorf("pop first value error")
	}
	if item.PopLastString() != "3" {
		t.Errorf("pop last value error")
	}
}

func TestGetStringByRange(t *testing.T) {
	var item memcache.MappingListItem
	item.Init()

	if item.PushStringList([]string{
		"a", "b", "c",
	}) != nil {
		t.Errorf("push values error")
	}

	var values []string

	values = item.GetStringByRange(0, 3)
	if values[0] != "a" || values[1] != "b" || values[2] != "c" || len(values) != 3 {
		t.Errorf("get values error: %s", values)
	}

	values = item.GetStringByRange(1, 2)
	if values[0] != "b" || len(values) != 1 {
		t.Errorf("get values error: %s", values)
	}

	values = item.GetStringByRange(1, 3)
	if values[0] != "b" || values[1] != "c" || len(values) != 2 {
		t.Errorf("get values error: %s", values)
	}

	values = item.GetStringByRange(-2, 10)
	if len(values) != 1 {
		t.Errorf("get values error: %s", values)
	}

	values = item.GetStringByRange(10, 20)
	if len(values) != 0 {
		t.Errorf("get values error: %s", values)
	}

	values = item.GetStringByRange(0, -1)
	if values[0] != "a" || values[1] != "b" || values[2] != "c" || len(values) != 3 {
		t.Errorf("get values error: %s", values)
	}

	values = item.GetStringByRange(-2, -1)
	if values[0] != "c" || len(values) != 1 {
		t.Errorf("get values error: %s", values)
	}
}
