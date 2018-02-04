package memcache

import (
	"errors"

	"github.com/emirpasic/gods/maps/hashmap"
)

// TypeMappingTableItem :
const TypeMappingTableItem = "table"

//
var (
	ErrTableFieldValueError = errors.New("Table field value error")
)

// MappingTableItem :
type MappingTableItem struct {
	ComplexMappingItem
}

// Init :
func (i *MappingTableItem) Init() {
	i.Value = hashmap.New()
}

// GetTable :
func (i *MappingTableItem) GetTable() *hashmap.Map {
	if i.Value == nil {
		return nil
	}
	return i.Value.(*hashmap.Map)
}

// Length :
func (i *MappingTableItem) Length() int {
	table := i.GetTable()
	if table == nil {
		return 0
	}
	return table.Size()
}

// GetString :
func (i *MappingTableItem) GetString(field string) string {
	table := i.GetTable()
	if table == nil {
		return ""
	}
	value, ok := table.Get(field)
	if !ok {
		return ""
	}
	return value.(string)
}

// SetString :
func (i *MappingTableItem) SetString(field string, value string) error {
	table := i.GetTable()
	if table == nil {
		return ErrTableFieldValueError
	}
	table.Put(field, value)
	return nil
}

// Delete :
func (i *MappingTableItem) Delete(field string) error {
	table := i.GetTable()
	if table == nil {
		return ErrTableFieldValueError
	}
	table.Remove(field)
	return nil
}

// Clear :
func (i *MappingTableItem) Clear() error {
	table := i.GetTable()
	if table == nil {
		return ErrTableFieldValueError
	}
	table.Clear()
	return nil
}

// Exists :
func (i *MappingTableItem) Exists(field string) bool {
	table := i.GetTable()
	if table == nil {
		return false
	}
	_, ok := table.Get(field)
	return ok
}

// GetMappings :
func (i *MappingTableItem) GetMappings() map[string]string {
	result := make(map[string]string, i.Length())
	table := i.GetTable()
	if table != nil {
		for _, field := range table.Keys() {
			value, ok := table.Get(field)
			if ok {
				result[field.(string)] = value.(string)
			}
		}
	}
	return result
}

// SetMappings :
func (i *MappingTableItem) SetMappings(mappings map[string]string) error {
	table := i.GetTable()
	if table == nil {
		return ErrTableFieldValueError
	}
	for field, value := range mappings {
		table.Put(field, value)
	}
	return nil
}
