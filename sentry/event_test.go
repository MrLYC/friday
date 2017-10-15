package sentry_test

import (
	"friday/sentry"
	"testing"
)

func TestRefreshID(t *testing.T) {
	event := sentry.Event{}

	id := event.ID
	event.RefreshID()
	if id == event.ID {
		t.Errorf("id: %v->%v", id, event.ID)
	}

	id = event.ID
	event.RefreshID()
	if id == event.ID {
		t.Errorf("id: %v->%v", id, event.ID)
	}
}

func TestCopy(t *testing.T) {
	event1 := &sentry.Event{}
	event2 := event1.Copy()
	if event1 == event2 {
		t.Errorf("event cpoy error")
	}
}
