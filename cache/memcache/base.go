package memcache

import (
	"sync"
	"time"
)

// ItemVistor :
type ItemVistor func(string, interface{})

// IMappingItem :
type IMappingItem interface {
	IsAvailable() bool
	SetValue(interface{})
	GetValue() interface{}
	SetExpireAt(time.Time)
	GetExpireAt() *time.Time
	IsExpireAt(time.Time) bool
	Length() int
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

// MappingItem :
type MappingItem struct {
	ExpireAt *time.Time
	Value    interface{}
	lock     sync.RWMutex
}

// Init :
func (i *MappingItem) Init() {
}

// IsAvailable :
func (i *MappingItem) IsAvailable() bool {
	if i.ExpireAt == nil {
		return true
	}
	return !i.IsExpireAt(time.Now())
}

// SetValue :
func (i *MappingItem) SetValue(value interface{}) {
	i.Value = value
}

// GetValue :
func (i *MappingItem) GetValue() interface{} {
	return i.Value
}

// SetExpireAt :
func (i *MappingItem) SetExpireAt(expired time.Time) {
	i.ExpireAt = &expired
}

// GetExpireAt :
func (i *MappingItem) GetExpireAt() *time.Time {
	return i.ExpireAt
}

// IsExpireAt :
func (i *MappingItem) IsExpireAt(t time.Time) bool {
	if i.ExpireAt == nil {
		return false
	}
	return t.After(*i.ExpireAt)
}

// Length :
func (i *MappingItem) Length() int {
	return 0
}

// Lock :
func (i *MappingItem) Lock() {
	i.lock.Lock()
}

// Unlock :
func (i *MappingItem) Unlock() {
	i.lock.Unlock()
}

// RLock :
func (i *MappingItem) RLock() {
	i.lock.RLock()
}

// RUnlock :
func (i *MappingItem) RUnlock() {
	i.lock.RUnlock()
}
