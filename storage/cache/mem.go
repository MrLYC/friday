package cache

import (
	"reflect"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/treemap"
)

// IMappingItem :
type IMappingItem interface {
	IsAvailable() bool
	SetValue(interface{})
	GetValue() interface{}
	SetExpireAt(time.Time)
	IsExpireAt(time.Time) bool
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

//
const (
	TypeMappingStringItem = "string"
	TypeMappingListItem   = "list"
	TypeMappingTableItem  = "table"
)

// MappingStringItem :
type MappingStringItem struct {
	MappingItem
}

// MappingListItem :
type MappingListItem struct {
	MappingItem
	RWLock sync.RWMutex
}

// MappingTableItem :
type MappingTableItem struct {
	MappingItem
	RWLock sync.RWMutex
}

// MemCache :
type MemCache struct {
	Mappings *treemap.Map
	RWLock   sync.RWMutex
}

// Init :
func (c *MemCache) Init() {
	c.Mappings = treemap.NewWithStringComparator()
}

// Close :
func (c *MemCache) Close() error {
	return nil
}

// Remove :
func (c *MemCache) Remove(key string) {
	c.RWLock.Lock()
	c.Mappings.Remove(key)
	c.RWLock.Unlock()
}

// Exists :
func (c *MemCache) Exists(key string) bool {
	c.RWLock.RLock()
	item, err := c.Get(key)
	c.RWLock.RUnlock()
	if err != nil {
		return false
	}
	return item.IsAvailable()
}

// IterItems :
func (c *MemCache) IterItems(f CacheItemIter) {
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()
	iter := c.Mappings.Iterator()

	for iter.Next() {
		f(iter.Key().(string), iter.Value())
	}
}

// Clean :
func (c *MemCache) Clean() int {
	var (
		now  = time.Now()
		list = arraylist.New()
	)

	c.IterItems(func(key string, value interface{}) {
		item := value.(IMappingItem)
		if item.IsExpireAt(now) {
			list.Add(key)
		}
	})

	if list.Size() == 0 {
		return 0
	}

	c.RWLock.Lock()
	defer c.RWLock.Unlock()

	iter := list.Iterator()
	for iter.Next() {
		c.Mappings.Remove(iter.Value())
	}
	return list.Size()
}

// Set :
func (c *MemCache) Set(key string, value IMappingItem) error {
	c.RWLock.Lock()
	c.Mappings.Put(key, value)
	defer c.RWLock.Unlock()
	return nil
}

// GetRaw :
func (c *MemCache) GetRaw(key string) (IMappingItem, error) {
	c.RWLock.RLock()
	item, ok := c.Mappings.Get(key)
	c.RWLock.RUnlock()
	if !ok {
		return nil, ErrItemNotFound
	}
	return item.(IMappingItem), nil
}

// Get :
func (c *MemCache) Get(key string) (IMappingItem, error) {
	item, err := c.GetRaw(key)
	if err != nil {
		return nil, err
	}
	if item.IsAvailable() {
		return item, nil
	}
	return nil, ErrItemNotFound
}

// Expire :
func (c *MemCache) Expire(key string, duration time.Duration) error {
	c.RWLock.RLock()
	item, err := c.Get(key)
	c.RWLock.RUnlock()
	if err != nil {
		return err
	}
	now := time.Now()
	if item.IsExpireAt(now) {
		c.Remove(key)
		return ErrItemExpired
	}

	item.SetExpireAt(now.Add(duration))
	return nil
}

// Size :
func (c *MemCache) Size() int {
	c.RWLock.RLock()
	defer c.RWLock.RUnlock()
	return c.Mappings.Size()
}

// TypeOf :
func (c *MemCache) TypeOf(key string) string {
	item, err := c.Get(key)
	if err != nil {
		return ""
	}
	switch item.(type) {
	case *MappingStringItem:
		return TypeMappingStringItem
	case *MappingListItem:
		return TypeMappingListItem
	case *MappingTableItem:
		return TypeMappingTableItem
	default:
		typ := reflect.TypeOf(item)
		return typ.Name()
	}
}

// GetStringItem :
func (c *MemCache) GetStringItem(key string) (*MappingStringItem, error) {
	item, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	switch item.(type) {
	case *MappingStringItem:
		return item.(*MappingStringItem), nil
	default:
		return nil, ErrItemTypeError
	}
}

// SetString :
func (c *MemCache) SetString(key string, value string) error {
	item := &MappingStringItem{}
	item.Value = value
	return c.Set(key, item)
}

// GetString :
func (c *MemCache) GetString(key string) (string, error) {
	item, err := c.GetStringItem(key)
	if err != nil {
		return "", err
	}
	return item.Value.(string), nil
}

// GetListItem :
func (c *MemCache) GetListItem(key string) (*MappingListItem, error) {
	item, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	switch item.(type) {
	case *MappingListItem:
		return item.(*MappingListItem), nil
	default:
		return nil, ErrItemTypeError
	}
}

// GetTableItem :
func (c *MemCache) GetTableItem(key string) (*MappingTableItem, error) {
	item, err := c.Get(key)
	if err != nil {
		return nil, err
	}
	switch item.(type) {
	case *MappingTableItem:
		return item.(*MappingTableItem), nil
	default:
		return nil, ErrItemTypeError
	}
}

// NewMemCache :
func NewMemCache() *MemCache {
	c := &MemCache{}
	c.Init()
	return c
}
