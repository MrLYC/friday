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
