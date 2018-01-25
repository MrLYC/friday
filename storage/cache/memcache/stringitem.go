package memcache

// TypeMappingStringItem :
const TypeMappingStringItem = "string"

// MappingStringItem :
type MappingStringItem struct {
	MappingItem
}

// GetString :
func (i *MappingStringItem) GetString() string {
	if i.Value == nil {
		return ""
	}
	return i.Value.(string)
}

// Length :
func (i *MappingStringItem) Length() int {
	return len(i.GetString())
}
