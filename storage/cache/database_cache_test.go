package cache_test

import (
	"testing"

	"friday/storage/cache"

	"friday/config"
	"friday/logging"
	"friday/storage"
	"friday/storage/migration"
)

func setUp() {
	config.Configuration.Init()
	config.Configuration.Read()
	logging.Init()

	command := migration.Command{}
	command.CreateMigrationTableIfNotExists()
	command.ActionRebuild()
}

func TestGetItemByKeyAndTag(t *testing.T) {
	setUp()
	key := "test1980237234"
	value := "test0981234879"
	c := &cache.DBCache{}
	c.Init()
	c.Conn.Create(&storage.Item{
		Key: key, Value: value,
		Tags: []*storage.ItemTag{
			c.TagTypeString,
		},
	})

	item, err := c.GetItemByKeyAndTag(key, c.TagTypeString.Name)
	if item == nil {
		t.Errorf("item not found")
	}
	if err != nil {
		t.Error(err)
	}
	if item.Value != value {
		t.Errorf("value not equal")
	}
	if len(item.Tags) != 1 {
		t.Errorf("tags not loaded")
	}
}

func TestMakeItemExpired(t *testing.T) {
	setUp()
	item := &storage.Item{
		Key: "test2345567856",
	}
	c := &cache.DBCache{}
	c.Init()
	c.Conn.Create(item)
	c.MakeItemExpired(item.Key, item.CreatedAt)
	c.Conn.Find(item, "key = ?", item.Key)
	if !item.CreatedAt.Equal(*(item.ExpireAt)) {
		t.Errorf("item not expired: %s", *(item.ExpireAt))
	}
}

func TestMakeItemTagExpired(t *testing.T) {
	setUp()
	itemTag := &storage.ItemTag{
		Name: "test6789235234",
	}
	c := &cache.DBCache{}
	c.Init()
	c.Conn.Create(itemTag)
	c.MakeItemTagExpired(itemTag.Name, itemTag.CreatedAt)
	c.Conn.Find(itemTag, "name = ?", itemTag.Name)
	if !itemTag.CreatedAt.Equal(*(itemTag.ExpireAt)) {
		t.Errorf("item tag not expired: %s", *(itemTag.ExpireAt))
	}
}

func TestGetKey(t *testing.T) {
	setUp()
	c := &cache.DBCache{}
	c.Init()
	item := &storage.Item{
		Key: "test8912876345", Value: "mrlyc",
		Tags: []*storage.ItemTag{
			c.TagTypeString,
		},
	}
	c.Conn.Create(item)
	val, err := c.GetKey(item.Key)
	if err != nil {
		t.Error(err)
	}
	if item.Value != val {
		t.Errorf("item not found: %s", val)
	}
}

func TestSetKey(t *testing.T) {
	setUp()
	c := &cache.DBCache{}
	c.Init()
	var err error
	key := "test0987123445"

	err = c.SetKey(key, "1")
	if err != nil {
		t.Error(err)
	}

	err = c.SetKey(key, "2")
	if err != nil {
		t.Error(err)
	}

	item1 := &storage.Item{}
	item2 := &storage.Item{}
	db := c.Conn.Where("key = ?", key)

	db.First(item1)
	if item1.ExpireAt == nil || item1.Value != "1" {
		t.Errorf("item error: %v", item1)
	}

	db.Last(item2)
	if item2.ExpireAt != nil || item2.Value != "2" {
		t.Errorf("item error: %v", item2)
	}
}
