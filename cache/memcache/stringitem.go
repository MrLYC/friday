package memcache

import (
	"friday/cache"
	"strconv"
)

// TypeMappingStringItem :
const TypeMappingStringItem = "string"

// MappingStringItem :
type MappingStringItem struct {
	MappingItem
}

// Init :
func (i *MappingStringItem) Init() {
	i.SetValue("")
}

// GetString :
func (i *MappingStringItem) GetString() string {
	value := i.GetValue()
	if value == nil {
		return ""
	}
	return value.(string)
}

// Length :
func (i *MappingStringItem) Length() int {
	return len(i.GetString())
}

// Add :
func (i *MappingStringItem) Add(num float64) (float64, error) {
	var (
		value float64
		err   error
	)
	value, err = strconv.ParseFloat(i.GetString(), 64)
	if err != nil && i.Length() > 1 {
		return value, cache.ErrItemValueError
	}

	value += num
	i.SetValue(strconv.FormatFloat(value, 'G', -1, 64))

	return value, nil
}
