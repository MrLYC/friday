package storage_test

import (
	"testing"
	"time"

	"friday/storage"
)

func TestCleanNormalItems(t *testing.T) {
	var err error
	db := storage.GetDBConnection()

	item := &storage.Item{
		Key:    "test098172387",
		Status: storage.ModelStatusNormal,
	}
	item.MakeExpireAt(time.Date(1991, time.November, 11, 0, 0, 0, 0, time.Local))
	err = db.Create(item).Error
	if err != nil {
		t.Error(err)
	}

	err = storage.CleanNormalItems(db)
	if err != nil {
		t.Error(err)
	}

	if !db.First(item).RecordNotFound() {
		t.Errorf("clean error")
	}
}

func TestCleanNormalItemTags(t *testing.T) {
	var err error
	db := storage.GetDBConnection()

	itemTag := &storage.ItemTag{
		Name:   "test9087123478",
		Status: storage.ModelStatusNormal,
	}
	itemTag.MakeExpireAt(time.Date(1991, time.November, 11, 0, 0, 0, 0, time.Local))
	err = db.Create(itemTag).Error
	if err != nil {
		t.Error(err)
	}

	err = storage.CleanNormalItemTags(db)
	if err != nil {
		t.Error(err)
	}

	if !db.First(itemTag).RecordNotFound() {
		t.Errorf("clean error")
	}
}
