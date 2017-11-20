package migration

import (
	"friday/storage"
	"time"
)

// ItemTag171112224339 :
type ItemTag171112224339 struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Name   string `gorm:"type:varchar(255)" sql:"index"`
	Status int    `sql:"index"`
}

// TableName :
func (ItemTag171112224339) TableName() string {
	return "item_tags"
}

// Item171112224339 :
type Item171112224339 struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Tags   []ItemTag171112224339 `gorm:"many2many:item_tags;"`
	Key    string                `gorm:"type:varchar(255)" sql:"index"`
	Value  string                `gorm:"type:varchar(65535)"`
	Type   string                `gorm:"type:varchar(64)" sql:"index"`
	Status int                   `sql:"index"`
}

// TableName :
func (Item171112224339) TableName() string {
	return "items"
}

// Migrate171112224339 :
func (c *Command) Migrate171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.AutoMigrate(ItemTag171112224339{}, Item171112224339{})
	return nil
}

// Rollback171112224339 :
func (c *Command) Rollback171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.DropTable(ItemTag171112224339{}.TableName(), Item171112224339{}.TableName())
	return nil
}
