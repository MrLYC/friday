package storage

import (
	"time"
)

// Model :
type Model struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// IsExpireAt :
func (m *Model) IsExpireAt(time time.Time) bool {
	if m.DeletedAt == nil {
		return false
	}
	return m.DeletedAt.Before(time)
}

// BeforeUpdate :
func (m *Model) BeforeUpdate() (err error) {
	m.UpdatedAt = time.Now()
	return nil
}

// Item :
type Item struct {
	Model

	Tags  []ItemTag `gorm:"many2many:item_tags;"`
	Key   string    `gorm:"type:varchar(255)" sql:"index"`
	Value string    `gorm:"type:varchar(65535)"`
	Type  string    `gorm:"type:varchar(64)" sql:"index"`
}

// ItemTag :
type ItemTag struct {
	Model

	Name string `gorm:"type:varchar(255)" sql:"index"`
}
