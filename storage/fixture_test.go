package storage_test

import (
	"os"
	"testing"

	"friday/config"
	"friday/logging"
	"friday/storage/migration"
)

func TestMain(m *testing.M) {
	config.Configuration.Init()
	config.Configuration.Read()
	logging.Init()

	command := migration.Command{}
	command.CreateMigrationTableIfNotExists()
	command.ActionRebuild()

	os.Exit(m.Run())
}
