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

	Items  []*Item              `gorm:"many2many:item_tag_x_refs;column:item"`
	Name   string               `gorm:"type:varchar(255)" sql:"index"`
	Status storage.TModelStatus `sql:"index"`
}

// Item :
type Item struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`

	Tags   []*ItemTag           `gorm:"many2many:item_tag_x_refs;column:tag"`
	Key    string               `gorm:"type:varchar(255)" sql:"index"`
	Value  string               `gorm:"type:text"`
	Type   string               `gorm:"type:varchar(64)" sql:"index"`
	Status storage.TModelStatus `sql:"index"`
}

// ItemTagXRef :
type ItemTagXRef struct {
	ItemTagID int `gorm:"item_tag_id;primary_key"`
	ItemID    int `gorm:"item_id;primary_key"`
}

// Migrate171112224339 :
func (c *Command) Migrate171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.AutoMigrate(ItemTag{}, Item{})
	return nil
}

// Rollback171112224339 :
func (c *Command) Rollback171112224339(migration *Migration, conn *storage.DatabaseConnection) error {
	conn.DropTable(ItemTagXRef{}, ItemTag{}, Item{})
	return nil
}
