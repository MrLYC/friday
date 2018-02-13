package memcache

import (
	"github.com/emirpasic/gods/lists/doublylinkedlist"
)

// TypeMappingListItem :
const TypeMappingListItem = "list"

// MappingListItem :
type MappingListItem struct {
	MappingItem
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

// GetStringByRange :
func (i *MappingListItem) GetStringByRange(index int, count int) []string {
	items := make([]string, 0, count)
	list := i.GetList()
	if list == nil {
		return items
	}

	if index < 0 {
		return items
	}

	for idx := 0; idx < count; idx++ {
		value, ok := list.Get(idx + index)
		if !ok {
			break
		}
		items = append(items, value.(string))
	}
	return items
}

// PushStringList :
func (i *MappingListItem) PushStringList(values []string) error {
	list := i.GetList()
	if list == nil {
		return ErrItemValueError
	}
	for _, value := range values {
		list.Add(value)
	}
	return nil
}

// Delete :
func (i *MappingListItem) Delete(index int) error {
	list := i.GetList()
	if list == nil {
		return ErrItemValueError
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
		panic(ErrItemValueError)
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
		panic(ErrItemValueError)
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
