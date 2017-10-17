package sentry_test

import (
	"friday/sentry"
	"testing"
)

type TestingTrigger struct {
	sentry.BaseTrigger
	Channel2 chan *sentry.Event
	Counter  int
}

func (t *TestingTrigger) Init(s *sentry.Sentry) {
	t.Name = "testing"
	t.Channel2 = make(chan *sentry.Event, 10)
	t.BaseTrigger.Init(s)
}

func (t *TestingTrigger) Run() {
	event1 := <-t.Channel2
	t.Counter += 1
	event2 := event1.Copy()
	event2.ID = event1.Name
	event2.Name = event1.ID
	t.Channel <- event2
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

	if itrigger.GetChannel() != trigger.Channel {
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
	trigger.Channel2 <- event1
	itrigger.Run()
	event2 := <-channel
	itrigger.Terminate()
	if event1.ID != event2.Name || event1.Name != event2.ID {
		t.Errorf("trigger run error")
	}
}
