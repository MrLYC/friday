package storage

import (
	"time"
)

// Model :
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ExpireAt  *time.Time `sql:"index"`
}

// IsExpireAt :
func (m *Model) IsExpireAt(time time.Time) bool {
	if m.ExpireAt == nil {
		return false
	}
	return m.ExpireAt.Before(time)
}

// BeforeUpdate :
func (m *Model) BeforeUpdate() (err error) {
	m.UpdatedAt = time.Now()
	return nil
}

// TModelStatus :
type TModelStatus int

// ModelStatus
const (
	ModelStatusBusy      TModelStatus = iota
	ModelStatusNormal    TModelStatus = iota
	ModelStatusProtected TModelStatus = iota
)

// Item :
type Item struct {
	Model

	Tags   []*ItemTag   `gorm:"many2many:item_tag_links;column:tag"`
	Key    string       `gorm:"type:varchar(255)" sql:"index"`
	Value  string       `gorm:"type:varchar(65535)"`
	Type   string       `gorm:"type:varchar(64)" sql:"index"`
	Status TModelStatus `sql:"index"`
}

// ItemTagTypes
const (
	ItemTagTypeString = "STRING"
	ItemTagTypeList   = "LIST"
	ItemTagTypeTable  = "TABLE"
)

// ItemTag :
type ItemTag struct {
	Model

	Items  []*Item      `gorm:"many2many:item_tag_links;column:item"`
	Name   string       `gorm:"type:varchar(255)" sql:"index"`
	Status TModelStatus `sql:"index"`
}
