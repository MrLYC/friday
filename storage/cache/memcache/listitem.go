package memcache

import (
	"errors"

	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

// TypeMappingListItem :
const TypeMappingListItem = "list"

//
var (
	ErrListItemValueError = errors.New("List item value error")
)

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

// GetString :
func (i *MappingListItem) GetString(index int) string {
	list := i.GetList()
	if list == nil {
		return ""
	}
	value, ok := list.Get(index)
	if !ok {
		return ""
	}
	return value.(string)
}

// Delete :
func (i *MappingListItem) Delete(index int) error {
	list := i.GetList()
	if list == nil {
		return ErrListItemValueError
	}
	list.Remove(index)
	return nil
}

// GetFirstString :
func (i *MappingListItem) GetFirstString() string {
	return i.GetString(0)
}

// PopFirstString :
func (i *MappingListItem) PopFirstString() string {
	list := i.GetList()
	if list == nil {
		return ""
	}
	value, ok := list.Get(0)
	if !ok {
		return ""
	}
	list.Remove(0)
	return value.(string)
}

// AppendFirstString :
func (i *MappingListItem) AppendFirstString(value string) {
	list := i.GetList()
	if list == nil {
		panic(ErrListItemValueError)
	}
	list.Insert(0, value)
}

// GetLastString :
func (i *MappingListItem) GetLastString() string {
	length := i.Length()
	if length <= 0 {
		return ""
	}
	return i.GetString(length - 1)
}

// AppendLastString :
func (i *MappingListItem) AppendLastString(value string) {
	list := i.GetList()
	if list == nil {
		panic(ErrListItemValueError)
	}
	list.Add(value)
}

// PopLastString :
func (i *MappingListItem) PopLastString() string {
	length := i.Length()
	if length <= 0 {
		return ""
	}
	list := i.GetList()
	value, ok := list.Get(length - 1)
	if !ok {
		return ""
	}
	list.Remove(length - 1)
	return value.(string)
}
