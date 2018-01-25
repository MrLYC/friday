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
	IsExpireAt(time.Time) bool
	Length() int
}

// IComplexMappingItem :
type IComplexMappingItem interface {
	IMappingItem
	Lock()
	Unlock()
	RLock()
	RUnlock()
}

// MappingItem :
type MappingItem struct {
	ExpireAt *time.Time
	Value    interface{}
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

// ComplexMappingItem :
type ComplexMappingItem struct {
	MappingItem
	RWLock sync.RWMutex
}

// Lock :
func (i *ComplexMappingItem) Lock() {
	i.RWLock.Lock()
}

// Unlock :
func (i *ComplexMappingItem) Unlock() {
	i.RWLock.Unlock()
}

// RLock :
func (i *ComplexMappingItem) RLock() {
	i.RWLock.RLock()
}

// RUnlock :
func (i *ComplexMappingItem) RUnlock() {
	i.RWLock.RUnlock()
}
