package storage_test

import (
	"testing"

	"friday/storage"
)

func TestCopyWithDB(t *testing.T) {
	db1 := storage.GetDBConnection()
	db2 := db1.CopyWithDB(db1.DB)
	if db1 == db2 {
		t.Errorf("copy error")
	}
}