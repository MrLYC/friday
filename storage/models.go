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
func (m *Model) IsExpireAt(t time.Time) bool {
	if m.ExpireAt == nil {
		return false
	}
	return m.ExpireAt.Before(t)
}

// MakeExpireAt :
func (m *Model) MakeExpireAt(t time.Time) {
	m.ExpireAt = &t
}

// BeforeUpdate :
func (m *Model) BeforeUpdate() error {
	m.UpdatedAt = time.Now()
	return nil
}

// TModelStatus :
type TModelStatus int

// ModelStatus
const (
	ModelStatusNormal    TModelStatus = iota
	ModelStatusBusy      TModelStatus = iota
	ModelStatusProtected TModelStatus = iota
)

// Item :
type Item struct {
	Model

	Tags   []*ItemTag   `gorm:"many2many:item_tag_x_refs;column:tag"`
	Key    string       `gorm:"type:varchar(255)" sql:"index"`
	Value  string       `gorm:"type:text"`
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

	Items  []*Item      `gorm:"many2many:item_tag_x_refs;column:item"`
	Name   string       `gorm:"type:varchar(255)" sql:"index"`
	Status TModelStatus `sql:"index"`
}

// ItemTagXRef :
type ItemTagXRef struct {
	ItemTagID int `gorm:"item_tag_id;primary_key"`
	ItemID    int `gorm:"item_id;primary_key"`
}
