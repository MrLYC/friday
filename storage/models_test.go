package storage_test

import (
	"friday/config"
	"friday/storage"
	"friday/storage/migration"
	"testing"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func rebuildDB() {
	command := migration.Command{}
	config.Configuration.Read()
	command.CreateMigrationTableIfNotExists()
	command.ActionRebuild()
}

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