package cache

import (
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/maps/treemap"

	"friday/utils"
)

// IMappingItem :
type IMappingItem interface {
	IsAvailable() bool
	GetValue() interface{}
	SetExpireAt(time.Time)
}

// MappingItem :
type MappingItem struct {
	ExpireAt *time.Time
	Value    interface{}
}

// MappingStringItem :
type MappingStringItem struct {
	MappingItem
}

// MappingListItem :
type MappingListItem struct {
	MappingItem
	Atomic utils.Atomic
}

// MappingTableItem :
type MappingTableItem struct {
	MappingItem
	Atomic utils.Atomic
}

// IsAvailable :
func (i *MappingItem) IsAvailable() bool {
	if i.ExpireAt == nil {
		return true
	}
	if !(*i.ExpireAt).Before(time.Now()) {
		return false
	}
	return true
}

// GetValue :
func (i *MappingItem) GetValue() interface{} {
	return i.Value
}

// SetExpireAt :
func (i *MappingItem) SetExpireAt(expired time.Time) {
	i.ExpireAt = &expired
}

// MemCache :
type MemCache struct {
	Mappings            *treemap.Map
	MappingsWriteAtomic utils.Atomic
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
	c.MappingsWriteAtomic.Run(func() {
		c.Mappings.Remove(key)
	})
}

// Set :
func (c *MemCache) Set(key string, value IMappingItem) error {
	c.MappingsWriteAtomic.Run(func() {
		c.Mappings.Put(key, value)
	})
	return nil
}

// Get :
func (c *MemCache) Get(key string) (IMappingItem, error) {
	value, ok := c.Mappings.Get(key)
	if !ok {
		return nil, ErrItemNotFound
	}
	item := value.(IMappingItem)
	if item.IsAvailable() {
		return item, nil
	}
	return nil, ErrItemNotFound
}

// Expire :
func (c *MemCache) Expire(key string, duration time.Duration) error {
	item, err := c.Get(key)
	if err != nil {
		return err
	}
	item.SetExpireAt(time.Now().Add(duration))
	return nil
}

// GetString :
func (c *MemCache) GetString(key string) (*MappingStringItem, error) {
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

// GetList :
func (c *MemCache) GetList(key string) (*MappingListItem, error) {
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

// GetTable :
func (c *MemCache) GetTable(key string) (*MappingTableItem, error) {
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

// SetKey :
func (c *MemCache) SetKey(key string, value string) error {
	item := &MappingStringItem{}
	item.Value = value
	return c.Set(key, item)
}

// GetKey :
func (c *MemCache) GetKey(key string) (string, error) {
	item, err := c.GetString(key)
	if err != nil {
		return "", err
	}
	return item.Value.(string), nil
}

// DeclareList :
func (c *MemCache) DeclareList(key string) (*MappingListItem, error) {
	item, err := c.GetList(key)
	if err == ErrItemNotFound {
		err = nil
		item = &MappingListItem{}
		item.Value = arraylist.New()
		c.Set(key, item)
	} else if err != nil {
		return nil, err
	}
	return item, err
}

// ListPush :
func (c *MemCache) ListPush(key string, value string) error {
	item, err := c.DeclareList(key)
	if err != nil {
		return err
	}
	item.Atomic.Run(func() {
		list := item.Value.(*arraylist.List)
		list.Add(value)
	})
	return nil
}

// ListPop :
func (c *MemCache) ListPop(key string) (string, error) {
	item, err := c.GetList(key)
	if err != nil {
		return "", err
	}
	value := ""
	item.Atomic.Run(func() {
		list := item.Value.(*arraylist.List)
		iter := list.Iterator()
		iter.End()
		end := iter.Value()
		if iter.Last() && end != nil {
			value = end.(string)
		}
		list.Remove(iter.Index())
	})
	return value, nil
}
