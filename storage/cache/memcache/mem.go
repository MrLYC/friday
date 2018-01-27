package memcache

import (
	"reflect"
	"sync"
	"time"

	"friday/storage/cache"

	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/maps/treemap"
)

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
func (c *MemCache) IterItems(f ItemVistor) {
	c.RWLock.RLock()
	iter := c.Mappings.Iterator()

	for iter.Next() {
		f(iter.Key().(string), iter.Value())
	}
	c.RWLock.RUnlock()
}

// Clean :
func (c *MemCache) Clean() int {
	var (
		now  = time.Now()
		list = singlylinkedlist.New()
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
	iter := list.Iterator()
	for iter.Next() {
		c.Mappings.Remove(iter.Value())
	}
	c.RWLock.Unlock()

	return list.Size()
}

// Set :
func (c *MemCache) Set(key string, value IMappingItem) error {
	c.RWLock.Lock()
	c.Mappings.Put(key, value)
	c.RWLock.Unlock()
	return nil
}

// GetRaw :
func (c *MemCache) GetRaw(key string) (IMappingItem, error) {
	c.RWLock.RLock()
	item, ok := c.Mappings.Get(key)
	c.RWLock.RUnlock()
	if !ok {
		return nil, cache.ErrItemNotFound
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
	return nil, cache.ErrItemNotFound
}

// Update :
func (c *MemCache) Update(key string, f ItemVistor) error {
	item, err := c.Get(key)
	if err != nil {
		return err
	}

	c.RWLock.Lock()
	f(key, item)
	c.RWLock.Unlock()
	return nil
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
		return cache.ErrItemExpired
	}

	item.SetExpireAt(now.Add(duration))
	return nil
}

// Size :
func (c *MemCache) Size() int {
	c.RWLock.RLock()
	size := c.Mappings.Size()
	c.RWLock.RUnlock()
	return size
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
		return nil, cache.ErrItemTypeError
	}
}

// SetString :
func (c *MemCache) SetString(key string, value string) error {
	item := &MappingStringItem{}
	item.SetValue(value)
	return c.Set(key, item)
}

// GetString :
func (c *MemCache) GetString(key string) (string, error) {
	item, err := c.GetStringItem(key)
	if err != nil {
		return "", err
	}
	return item.GetString(), nil
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
		return nil, cache.ErrItemTypeError
	}
}

// DeclareListItem :
func (c *MemCache) DeclareListItem(key string) (*MappingListItem, error) {
	item, err := c.GetListItem(key)
	if err != nil {
		item = &MappingListItem{}
		item.Init()
		err = c.Set(key, item)
	}
	return item, err
}

// GetListLength :
func (c *MemCache) GetListLength(key string) (int, error) {
	item, err := c.DeclareListItem(key)
	if err == nil {
		return 0, err
	}

	item.RLock()
	length := item.Length()
	item.RUnlock()
	return length, nil
}

// PopListString :
func (c *MemCache) PopListString(key string) (string, error) {
	item, err := c.DeclareListItem(key)
	if err == nil {
		return "", err
	}
	item.RLock()
	value := item.PopLastString()
	item.RUnlock()
	return value, nil
}

// AppendListString :
func (c *MemCache) AppendListString(key string, value string) error {
	item, err := c.DeclareListItem(key)
	if err == nil {
		return err
	}
	item.RLock()
	item.AppendLastString(value)
	item.RUnlock()
	return nil
}

// LPopListString :
func (c *MemCache) LPopListString(key string) (string, error) {
	item, err := c.DeclareListItem(key)
	if err == nil {
		return "", err
	}
	item.Lock()
	value := item.PopFirstString()
	item.Unlock()
	return value, nil
}

// LAppendString :
func (c *MemCache) LAppendString(key string, value string) error {
	item, err := c.DeclareListItem(key)
	if err == nil {
		return err
	}
	item.Lock()
	item.AppendFirstString(value)
	item.Unlock()
	return nil
}

// GetListString :
func (c *MemCache) GetListString(key string, index int) (string, error) {
	item, err := c.DeclareListItem(key)
	if err == nil {
		return "", err
	}
	item.RLock()
	value := item.GetString(index)
	item.RUnlock()
	return value, nil
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
		return nil, cache.ErrItemTypeError
	}
}

// NewMemCache :
func NewMemCache() *MemCache {
	c := &MemCache{}
	c.Init()
	return c
}
