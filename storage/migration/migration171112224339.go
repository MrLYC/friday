package migration

import (
	"friday/storage"
	"time"
)

// ItemTag1 :
type ItemTag1 struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Name   string `gorm:"type:varchar(255)" sql:"index"`
	Status int    `sql:"index"`
}

// TableName :
func (ItemTag1) TableName() string {
	return "item_tags"
}

// Item1 :
type Item1 struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Tags   []ItemTag1 `gorm:"many2many:item_tags;"`
	Key    string     `gorm:"type:varchar(255)" sql:"index"`
	Value  string     `gorm:"type:varchar(65535)"`
	Type   string     `gorm:"type:varchar(64)" sql:"index"`
	Status int        `sql:"index"`
}

// TableName :
func (Item1) TableName() string {
	return "items"
}

// Migrate171112224339 :
func (c *Command) Migrate171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.AutoMigrate(ItemTag1{}, Item1{})
	return nil
}

// Rollback171112224339 :
func (c *Command) Rollback171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.DropTable(ItemTag1{}.TableName(), Item1{}.TableName())
	return nil
}
