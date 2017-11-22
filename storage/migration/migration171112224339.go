package migration

import (
	"friday/storage"
	"time"
)

// ItemTag :
type ItemTag struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Items  []Item `gorm:"many2many:item_tag_links;column:item"`
	Name   string `gorm:"type:varchar(255)" sql:"index"`
	Status int    `sql:"index"`
}

// TableName :
func (ItemTag) TableName() string {
	return "item_tags"
}

// Item :
type Item struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Tags   []ItemTag `gorm:"many2many:item_tag_links;column:tag"`
	Key    string    `gorm:"type:varchar(255)" sql:"index"`
	Value  string    `gorm:"type:varchar(65535)"`
	Type   string    `gorm:"type:varchar(64)" sql:"index"`
	Status int       `sql:"index"`
}

// TableName :
func (Item) TableName() string {
	return "items"
}

// Migrate171112224339 :
func (c *Command) Migrate171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.AutoMigrate(ItemTag{}, Item{})
	return nil
}

// Rollback171112224339 :
func (c *Command) Rollback171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.DropTable(ItemTag{}.TableName(), Item{}.TableName())
	conn.DropTable("item_tag_links")
	return nil
}
