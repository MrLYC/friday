package memcache

import (
	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

// TypeMappingListItem :
const TypeMappingListItem = "list"

// MappingListItem :
type MappingListItem struct {
	ComplexMappingItem
}

// GetList :
func (i *MappingListItem) GetList() *doublylinkedlist.List {
	if i.Value == nil {
		return nil
	}
	return i.Value.(*doublylinkedlist.List)
}

// Init :
func (i *MappingListItem) Init() {
	i.Value = doublylinkedlist.New()
}

// Length :
func (i *MappingListItem) Length() int {
	list := i.GetList()
	if list == nil {
		return 0
	}
	return list.Size()
}

// GetFirstString :
func (i *MappingListItem) GetFirstString() string {
	list := i.GetList()
	if list == nil {
		return ""
	}
	value, ok := list.Get(0)
	if !ok {
		return ""
	}
	return value.(string)
}

// GetLastString :
func (i *MappingListItem) GetLastString() string {
	length := i.Length()
	if length <= 0 {
		return ""
	}
	list := i.GetList()
	if list == nil {
		return ""
	}
	value, ok := list.Get(length - 1)
	if !ok {
		return ""
	}
	return value.(string)
}
