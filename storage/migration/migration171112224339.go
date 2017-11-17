package migration

import "friday/storage"

// TableNameItemTag :
var TableNameItemTag = "item_tags"

// TableNameItem :
var TableNameItem = "items"

// Migrate171112224339 :
func (c *Command) Migrate171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.AutoMigrate(storage.ItemTag{}, storage.Item{})
	return nil
}

// Rollback171112224339 :
func (c *Command) Rollback171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.DropTable(TableNameItemTag, TableNameItem)
	return nil
}
