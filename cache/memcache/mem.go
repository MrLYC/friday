package memcache

import (
	"friday/cache"
	"reflect"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/maps/treemap"
)

// MemCache :
type MemCache struct {
	Mappings *treemap.Map
	lock     sync.RWMutex
}

// Init :
func (c *MemCache) Init() {
	c.Mappings = treemap.NewWithStringComparator()
}

// Close :
func (c *MemCache) Close() error {
	return nil
}

// Delete :
func (c *MemCache) Delete(key string) {
	c.lock.Lock()
	c.Mappings.Remove(key)
	c.lock.Unlock()
}

// Exists :
func (c *MemCache) Exists(key string) bool {
	c.lock.RLock()
	item, err := c.GetItem(key)
	c.lock.RUnlock()
	if err != nil {
		return false
	}
	return item.IsAvailable()
}

// Scan :
func (c *MemCache) Scan(f ItemVistor) {
	c.lock.RLock()
	iter := c.Mappings.Iterator()

	for iter.Next() {
		f(iter.Key().(string), iter.Value())
	}
	c.lock.RUnlock()
}

// Clean :
func (c *MemCache) Clean() int {
	var (
		now  = time.Now()
		list = singlylinkedlist.New()
	)

	c.Scan(func(key string, value interface{}) {
		item := value.(IMappingItem)
		if item.IsExpireAt(now) {
			list.Add(key)
		}
	})

	if list.Size() == 0 {
		return 0
	}

	c.lock.Lock()
	iter := list.Iterator()
	for iter.Next() {
		c.Mappings.Remove(iter.Value())
	}
	c.lock.Unlock()

	return list.Size()
}

// Base

// SetItem :
func (c *MemCache) SetItem(key string, value IMappingItem) error {
	c.lock.Lock()
	c.Mappings.Put(key, value)
	c.lock.Unlock()
	return nil
}

// GetRaw :
func (c *MemCache) GetRaw(key string) (IMappingItem, error) {
	c.lock.RLock()
	item, ok := c.Mappings.Get(key)
	c.lock.RUnlock()
	if !ok {
		return nil, cache.ErrItemNotFound
	}
	return item.(IMappingItem), nil
}

// GetItem :
func (c *MemCache) GetItem(key string) (IMappingItem, error) {
	item, err := c.GetRaw(key)
	if err != nil {
		return nil, err
	}
	if item.IsAvailable() {
		return item, nil
	}
	return nil, cache.ErrItemNotFound
}

// GetStringItem :
func (c *MemCache) GetStringItem(key string) (*MappingStringItem, error) {
	item, err := c.GetItem(key)
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

// DeclareStringItem :
func (c *MemCache) DeclareStringItem(key string) (*MappingStringItem, error) {
	item, err := c.GetStringItem(key)
	if err == cache.ErrItemNotFound {
		item = &MappingStringItem{}
		item.Init()
		err = c.SetItem(key, item)
	}
	return item, err
}

// GetListItem :
func (c *MemCache) GetListItem(key string) (*MappingListItem, error) {
	item, err := c.GetItem(key)
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
	if err == cache.ErrItemNotFound {
		item = &MappingListItem{}
		item.Init()
		err = c.SetItem(key, item)
	}
	return item, err
}

// GetTableItem :
func (c *MemCache) GetTableItem(key string) (*MappingTableItem, error) {
	item, err := c.GetItem(key)
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

// DeclareTableItem :
func (c *MemCache) DeclareTableItem(key string) (*MappingTableItem, error) {
	item, err := c.GetTableItem(key)
	if err == cache.ErrItemNotFound {
		item = &MappingTableItem{}
		item.Init()
		err = c.SetItem(key, item)
	}
	return item, err
}

// API

// Update :
func (c *MemCache) Update(key string, f ItemVistor) error {
	item, err := c.GetItem(key)
	if err != nil {
		return err
	}

	c.lock.Lock()
	f(key, item)
	c.lock.Unlock()
	return nil
}

// TimeToLive :
func (c *MemCache) TimeToLive(key string) time.Duration {
	item, err := c.GetRaw(key)
	if err != nil {
		return time.Duration(0)
	}

	now := time.Now()
	if item.IsExpireAt(now) {
		return time.Duration(0)
	}

	expireAt := item.GetExpireAt()
	if expireAt == nil {
		return time.Hour * 1000
	}
	return now.Sub(*expireAt)
}

// Expire :
func (c *MemCache) Expire(key string, duration time.Duration) error {
	c.lock.RLock()
	item, err := c.GetItem(key)
	c.lock.RUnlock()
	if err != nil {
		return err
	}
	now := time.Now()
	if item.IsExpireAt(now) {
		return cache.ErrItemExpired
	}

	item.SetExpireAt(now.Add(duration))
	return nil
}

// Size :
func (c *MemCache) Size() int {
	c.lock.RLock()
	size := c.Mappings.Size()
	c.lock.RUnlock()
	return size
}

// TypeOf :
func (c *MemCache) TypeOf(key string) string {
	item, err := c.GetItem(key)
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

// String API

// Set :
func (c *MemCache) Set(key string, value string) error {
	item, err := c.DeclareStringItem(key)
	if err != nil {
		return err
	}

	item.Lock()
	item.SetValue(value)
	item.Unlock()
	return nil
}

// Get :
func (c *MemCache) Get(key string) (string, error) {
	item, err := c.GetStringItem(key)
	if err != nil {
		return "", err
	}

	item.RLock()
	value := item.GetString()
	item.RUnlock()
	return value, err
}

// StrLen :
func (c *MemCache) StrLen(key string) (int, error) {
	item, err := c.GetStringItem(key)
	if err != nil {
		return 0, err
	}

	item.RLock()
	value := item.Length()
	item.RUnlock()
	return value, err
}

// IncrBy :
func (c *MemCache) IncrBy(key string, num float64) (float64, error) {
	item, err := c.DeclareStringItem(key)
	if err != nil {
		return 0, err
	}

	item.Lock()
	value, err := item.Add(num)
	item.Unlock()
	return value, err
}

// DecrBy :
func (c *MemCache) DecrBy(key string, num float64) (float64, error) {
	return c.IncrBy(key, -num)
}

// Incr :
func (c *MemCache) Incr(key string) (int64, error) {
	value, err := c.IncrBy(key, 1)
	return int64(value), err
}

// Decr :
func (c *MemCache) Decr(key string) (int64, error) {
	value, err := c.IncrBy(key, -1)
	return int64(value), err
}

// List API

// LLen :
func (c *MemCache) LLen(key string) (int, error) {
	item, err := c.GetListItem(key)
	if err != nil {
		return 0, err
	}

	item.RLock()
	length := item.Length()
	item.RUnlock()
	return length, nil
}

// RPop :
func (c *MemCache) RPop(key string) (string, error) {
	item, err := c.GetListItem(key)
	if err != nil {
		return "", err
	}
	item.Lock()
	value := item.PopLastString()
	item.Unlock()
	return value, nil
}

// RPush :
func (c *MemCache) RPush(key string, value string) error {
	item, err := c.DeclareListItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	item.AppendLastString(value)
	item.Unlock()
	return nil
}

// LPop :
func (c *MemCache) LPop(key string) (string, error) {
	item, err := c.GetListItem(key)
	if err != nil {
		return "", err
	}
	item.Lock()
	value := item.PopFirstString()
	item.Unlock()
	return value, nil
}

// LPush :
func (c *MemCache) LPush(key string, value string) error {
	item, err := c.DeclareListItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	item.AppendFirstString(value)
	item.Unlock()
	return nil
}

// LIndex :
func (c *MemCache) LIndex(key string, index int) (string, error) {
	item, err := c.GetListItem(key)
	if err != nil {
		return "", err
	}
	item.RLock()
	value := item.GetString(index)
	item.RUnlock()
	return value, nil
}

// LSet :
func (c *MemCache) LSet(key string, index int, value string) error {
	item, err := c.DeclareListItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	err = item.SetAt(index, value)
	item.Unlock()
	return err
}

// LRange :
func (c *MemCache) LRange(key string, start int, stop int) ([]string, error) {
	item, err := c.GetListItem(key)
	if err != nil {
		return nil, err
	}
	item.RLock()
	value := item.GetStringByRange(start, stop)
	item.RUnlock()
	return value, nil
}

// LDel :
func (c *MemCache) LDel(key string, index int) error {
	item, err := c.GetListItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	err = item.Delete(index)
	item.Unlock()
	return err
}

// Table API

// HGet :
func (c *MemCache) HGet(key string, field string) (string, error) {
	item, err := c.GetTableItem(key)
	if err != nil {
		return "", err
	}
	item.RLock()
	value := item.GetString(field)
	item.RUnlock()
	return value, nil
}

// HSet :
func (c *MemCache) HSet(key string, field string, value string) error {
	item, err := c.DeclareTableItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	item.SetString(field, value)
	item.Unlock()
	return nil
}

// HDel :
func (c *MemCache) HDel(key string, field string) error {
	item, err := c.GetTableItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	item.Delete(field)
	item.Unlock()
	return nil
}

// HClear :
func (c *MemCache) HClear(key string) error {
	item, err := c.GetTableItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	item.Clear()
	item.Unlock()
	return nil
}

// HExists :
func (c *MemCache) HExists(key string, field string) bool {
	item, err := c.GetTableItem(key)
	if err != nil {
		return false
	}
	item.Lock()
	value := item.Exists(field)
	item.Unlock()
	return value
}

// HGetAll :
func (c *MemCache) HGetAll(key string) (map[string]string, error) {
	item, err := c.GetTableItem(key)
	if err != nil {
		return nil, err
	}
	item.Lock()
	value := item.GetAllMappings()
	item.Unlock()
	return value, nil
}

// HMGet :
func (c *MemCache) HMGet(key string, fields []string) (map[string]string, error) {
	item, err := c.GetTableItem(key)
	if err != nil {
		return nil, err
	}
	item.Lock()
	value := item.GetMappings(fields)
	item.Unlock()
	return value, nil
}

// HMSet :
func (c *MemCache) HMSet(key string, mappings map[string]string) error {
	item, err := c.GetTableItem(key)
	if err != nil {
		return err
	}
	item.Lock()
	err = item.SetMappings(mappings)
	item.Unlock()
	return err
}

// NewMemCache :
func NewMemCache() *MemCache {
	c := &MemCache{}
	c.Init()
	return c
}
