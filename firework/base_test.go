package firework_test

import (
	"friday/firework"
	"testing"
)

func TestRefreshID(t *testing.T) {
	f := firework.Firework{}

	id := f.ID
	f.RefreshID()
	if id == f.ID {
		t.Errorf("id: %v->%v", id, f.ID)
	}

	id = f.ID
	f.RefreshID()
	if id == f.ID {
		t.Errorf("id: %v->%v", id, f.ID)
	}
}

func TestCopy(t *testing.T) {
	f1 := &firework.Firework{}
	f2 := f1.Copy()
	if f1 == f2 {
		t.Errorf("f cpoy error")
	}
}
