package migration

import "friday/storage"

// TableNameItemTag :
var TableNameItemTag = "item_tags"

// TableNameItem :
var TableNameItem = "items"

// Migrate171112224339 :
func (c *Command) Migrate171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	// ItemTag :
	type ItemTag struct {
		Model

		Name string `gorm:"type:varchar(255)" sql:"index"`
	}
	// Item :
	type Item struct {
		Model

		Tags  []ItemTag `gorm:"many2many:item_tags;"`
		Key   string    `gorm:"type:varchar(255)" sql:"index"`
		Value string    `gorm:"type:varchar(65535)"`
		Type  string    `gorm:"type:varchar(64)" sql:"index"`
	}
	conn.AutoMigrate(ItemTag{}, Item{})
	return nil
}

// Rollback171112224339 :
func (c *Command) Rollback171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.DropTable(TableNameItemTag, TableNameItem)
	return nil
}
