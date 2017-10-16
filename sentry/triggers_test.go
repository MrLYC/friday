package sentry_test

import (
	"friday/sentry"
	"testing"
)

type TestingTrigger struct {
	sentry.BaseTrigger
}

func (t *TestingTrigger) Init(sentry *sentry.Sentry) {
	t.BaseTrigger.Name = "testing"
	t.BaseTrigger.Init(sentry)
}

func (t *TestingTrigger) Run() {
	event1 := <-t.BaseTrigger.Channel
	event2 := event1.Copy()
	event2.ID = event1.Name
	event2.Name = event1.ID
	t.BaseTrigger.Channel <- event2
}

func TestBaseTriggerInit(t *testing.T) {
	var (
		trigger                  = TestingTrigger{}
		itrigger sentry.ITrigger = &trigger
	)
	trigger.Init(nil)

	if itrigger.GetName() != trigger.Name {
		t.Errorf("name error")
	}

	if itrigger.GetChannel() != trigger.BaseTrigger.Channel {
		t.Errorf("channel error")
	}
}

func TestBaseTriggerNewEvent(t *testing.T) {
	trigger := TestingTrigger{}
	trigger.Init(nil)
	name := "test"

	event := trigger.NewEvent(name)
	if event.ID == "" {
		t.Errorf("event id is empty")
	}
	if event.Channel != trigger.GetName() {
		t.Errorf("event channel error")
	}
	if event.Name != name {
		t.Errorf("event name error")
	}
}

func TestBaseTriggerRun(t *testing.T) {
	var (
		trigger                  = TestingTrigger{}
		itrigger sentry.ITrigger = &trigger
	)

	itrigger.Init(nil)
	itrigger.Ready()
	channel := trigger.GetChannel()
	event1 := trigger.NewEvent("mrlyc")
	channel <- event1
	itrigger.Run()
	event2 := <-channel
	itrigger.Terminate()
	if event1.ID != event2.Name || event1.Name != event2.ID {
		t.Errorf("trigger run error")
	}
}
