package differ

import (
	"container/list"
	"fmt"
)

// Item :
type Item interface {
	GetID() string
	GetVersion() string
}

// DiffCallback :
type DiffCallback func(latest Item, earliest Item)

type mappingItem struct {
	version string
	item    Item
}

// Differ :
type Differ struct {
	latest   map[string]*mappingItem
	earliest map[string]*mappingItem

	added   *list.List
	updated *list.List
	deleted *list.List
}

// Init :
func (d *Differ) Init(latest *list.List, earliest *list.List) {
	d.added = list.New()
	d.updated = list.New()
	d.deleted = list.New()

	d.latest = make(map[string]*mappingItem, latest.Len())
	for i := latest.Front(); i != nil; i = i.Next() {
		item := i.Value.(Item)
		d.latest[item.GetID()] = &mappingItem{
			version: item.GetVersion(),
			item:    item,
		}
	}

	d.earliest = make(map[string]*mappingItem, earliest.Len())
	for i := earliest.Front(); i != nil; i = i.Next() {
		item := i.Value.(Item)
		id := item.GetID()
		version := item.GetVersion()
		d.earliest[id] = &mappingItem{
			version: version,
			item:    item,
		}

		litem, ok := d.latest[id]
		if !ok {
			d.deleted.PushBack(id)
		} else if litem.version != version {
			d.updated.PushBack(id)
		}
	}

	for id := range d.latest {
		_, ok := d.earliest[id]
		if !ok {
			d.added.PushBack(id)
		}
	}
}

// IterAdded :
func (d *Differ) IterAdded(callback DiffCallback) error {
	for i := d.added.Front(); i != nil; i = i.Next() {
		id := i.Value.(string)
		latest, ok := d.latest[id]
		if !ok {
			return fmt.Errorf("ErrDiffItemNotFound: %s", id)
		}
		callback(latest.item, nil)
	}
	return nil
}

// IterDeleted :
func (d *Differ) IterDeleted(callback DiffCallback) error {
	for i := d.deleted.Front(); i != nil; i = i.Next() {
		id := i.Value.(string)
		earliest, ok := d.earliest[id]
		if !ok {
			return fmt.Errorf("ErrDiffItemNotFound: %s", id)
		}
		callback(nil, earliest.item)
	}
	return nil
}

// IterUpdated :
func (d *Differ) IterUpdated(callback DiffCallback) error {
	var (
		ok       bool
		latest   *mappingItem
		earliest *mappingItem
	)
	for i := d.updated.Front(); i != nil; i = i.Next() {
		id := i.Value.(string)
		latest, ok = d.latest[id]
		if !ok {
			return fmt.Errorf("ErrDiffItemNotFound: %s", id)
		}
		earliest, ok = d.earliest[id]
		if !ok {
			return fmt.Errorf("ErrDiffItemNotFound: %s", id)
		}
		callback(latest.item, earliest.item)
	}
	return nil
}
