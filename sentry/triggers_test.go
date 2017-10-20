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

func (t *TestingTrigger) Init(s sentry.ISentry) {
	t.Name = "testing"
	t.Channel2 = make(chan *sentry.Event, 10)
	t.Sentry = s
	t.EventTemplate = &sentry.Event{
		Channel: t.Name,
	}
	if s != nil {
		t.Channel = s.DeclareChannel(t.Name)
	}
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
	trigger.SetControlChannel(make(chan *sentry.Event, 1))

	if itrigger.GetName() != trigger.Name {
		t.Errorf("name error")
	}
}

func TestBaseTriggerNewEvent(t *testing.T) {
	trigger := TestingTrigger{}
	trigger.Init(nil)
	trigger.SetControlChannel(make(chan *sentry.Event, 1))
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
		s                        = &TestingSentry{}
		trigger                  = TestingTrigger{}
		itrigger sentry.ITrigger = &trigger
	)
	s.Init([]sentry.ITrigger{itrigger}, []sentry.IHandler{})
	itrigger.SetControlChannel(make(chan *sentry.Event, 1))
	itrigger.Ready()
	channel := trigger.Channel
	event1 := trigger.NewEvent("mrlyc")
	trigger.Channel2 <- event1
	itrigger.Run()
	event2 := <-channel
	itrigger.Terminate()
	if event1.ID != event2.Name || event1.Name != event2.ID {
		t.Errorf("trigger run error")
	}
}
