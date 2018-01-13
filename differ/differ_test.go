package differ_test

import (
	"container/list"
	"fmt"
	"friday/differ"
	"testing"
)

type TestingDiffItem struct {
	ID    int
	Value int
}

func (i *TestingDiffItem) GetID() string {
	return fmt.Sprintf("%v", i.ID)
}

func (i *TestingDiffItem) GetVersion() string {
	return fmt.Sprintf("%v", i.Value)
}

func TestDiff(t *testing.T) {
	var (
		d        = &differ.Differ{}
		latest   = list.New()
		earliest = list.New()
	)

	latest.PushBack(&TestingDiffItem{
		ID:    0,
		Value: 0,
	})
	latest.PushBack(&TestingDiffItem{
		ID:    1,
		Value: 0,
	})

	earliest.PushBack(&TestingDiffItem{
		ID:    1,
		Value: 1,
	})
	earliest.PushBack(&TestingDiffItem{
		ID:    2,
		Value: 2,
	})

	d.Init(latest, earliest)

	flag := make(map[string]string)
	d.IterAdded(func(item1 differ.Item, item2 differ.Item) {
		if item2 != nil {
			t.Errorf("item2 not nil")
		}
		_, ok := flag[item1.GetID()]
		if ok {
			t.Errorf("item1 error")
		}
		flag[item1.GetID()] = "added"
	})
	d.IterDeleted(func(item1 differ.Item, item2 differ.Item) {
		if item1 != nil {
			t.Errorf("item2 not nil")
		}
		_, ok := flag[item2.GetID()]
		if ok {
			t.Errorf("item2 error")
		}
		flag[item2.GetID()] = "deleted"
	})
	d.IterUpdated(func(item1 differ.Item, item2 differ.Item) {
		if item1.GetVersion() == item2.GetVersion() {
			t.Errorf("diff error")
		}
		_, ok := flag[item1.GetID()]
		if ok {
			t.Errorf("item1 error")
		}
		flag[item1.GetID()] = "updated"
	})

	if flag["0"] != "added" {
		t.Errorf("IterAdded error")
	}
	if flag["1"] != "updated" {
		t.Errorf("IterUpdated error")
	}
	if flag["2"] != "deleted" {
		t.Errorf("IterDeleted error")
	}
}
