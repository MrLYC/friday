package storage_test

import (
	"friday/storage"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestModelIsExpireAt(t *testing.T) {
	model := storage.Model{}

	time0 := time.Date(2017, 11, 10, 21, 46, 0, 0, time.Local)
	time1 := time.Date(2017, 11, 10, 21, 46, 1, 0, time.Local)
	time2 := time.Date(2017, 11, 10, 21, 46, 2, 0, time.Local)

	if model.IsExpireAt(time1) {
		t.Errorf("expire error")
	}

	model.ExpireAt = &time1
	if model.IsExpireAt(time0) {
		t.Errorf("expire error")
	}
	if model.IsExpireAt(time1) {
		t.Errorf("expire error")
	}
	if !model.IsExpireAt(time2) {
		t.Errorf("expire error")
	}
}

func TestModelItem(t *testing.T) {
	tagName := "test_72873489712"
	itemKey := "test_79871238076"
	itemValue := "0"
	item := &storage.Item{
		Key:   itemKey,
		Value: itemValue,
		Tags: []*storage.ItemTag{
			&storage.ItemTag{
				Name: tagName,
			},
			&storage.ItemTag{
				Name: "test_0982345234",
			},
		},
	}

	conn := storage.GetDBConnection()
	if err := conn.Create(item).Error; err != nil {
		t.Error(err)
	}

	if err := conn.Create(&storage.Item{
		Key:   itemKey,
		Value: "8761234876",
	}).Error; err != nil {
		t.Error(err)
	}

	if err := conn.Create(&storage.Item{
		Key:   itemKey,
		Value: "0871234786",
		Tags: []*storage.ItemTag{
			&storage.ItemTag{
				Name: tagName,
			},
		},
	}).Error; err != nil {
		t.Error(err)
	}

	queryItemTag := &storage.ItemTag{}

	if err := conn.Where(
		"name = ?", tagName,
	).Preload(
		"Items", "key = ?", itemKey,
	).First(queryItemTag).Error; err != nil {
		t.Error(err)
	}

	if len(queryItemTag.Items) != 1 {
		t.Errorf("query error: %v", queryItemTag.Items)
	}

	queryItem := queryItemTag.Items[0]
	if queryItem.Value != item.Value {
		t.Errorf("value error: %s", queryItem.Value)
	}

	if err := conn.Model(queryItem).Related(&(queryItem.Tags), "Tags").Error; err != nil {
		t.Error(err)
	}
	if len(queryItem.Tags) != 2 {
		t.Errorf("tag query failed")
	}
}
