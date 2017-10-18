package sentry_test

import (
	"friday/sentry"
	"testing"
)

func TestSentryInit(t *testing.T) {
	var (
		s       = &sentry.Sentry{}
		trigger = &TestingTrigger{}
		handler = &TestingHandler{}
		ok      bool
	)
	s.Init([]sentry.ITrigger{trigger}, []sentry.IHandler{handler})
	if trigger.Sentry != s {
		t.Errorf("Trigger init error")
	}
	triggers := s.Triggers[trigger.GetName()]
	ok = false
	for _, tgr := range triggers {
		if tgr == trigger {
			ok = true
		}
	}
	if !ok {
		t.Errorf("trigger error")
	}

	if handler.Sentry != s {
		t.Errorf("Handler init error")
	}
	handlers := s.Handlers[handler.GetName()]
	ok = false
	for _, hdr := range handlers {
		if hdr == handler {
			ok = true
		}
	}
	if !ok {
		t.Errorf("handler error")
	}

	_, ok = s.Channels[sentry.ChanNameBroadcast]
	if !ok {
		t.Errorf("broadcast channel error")
	}

	_, ok = s.Channels[sentry.ChanNameInternal]
	if !ok {
		t.Errorf("internal channel error")
	}
}

func TestSentryFlow(t *testing.T) {
	var (
		s       = &sentry.Sentry{}
		trigger = &TestingTrigger{}
		handler = &TestingHandler{}
	)
	s.Init([]sentry.ITrigger{trigger}, []sentry.IHandler{handler})
	s.Ready()
	if trigger.Status != sentry.StatusControllerReady {
		t.Errorf("trigger error")
	}
	if handler.Status != sentry.StatusControllerReady {
		t.Errorf("handler error")
	}
	if s.Status != sentry.StatusControllerReady {
		t.Errorf("sentry error")
	}
	s.Terminate()
	if trigger.Status != sentry.StatusControllerTerminated {
		t.Errorf("trigger error")
	}
	if handler.Status != sentry.StatusControllerTerminated {
		t.Errorf("handler error")
	}
	if s.Status != sentry.StatusControllerTerminated {
		t.Errorf("sentry error")
	}
	s.Kill()
	if trigger.Status != sentry.StatusControllerTerminated {
		t.Errorf("trigger error")
	}
	if handler.Status != sentry.StatusControllerTerminated {
		t.Errorf("handler error")
	}
	if s.Status != sentry.StatusControllerTerminated {
		t.Errorf("sentry error")
	}
	s.Status = sentry.StatusControllerTerminating
	s.Kill()
	if s.Status != sentry.StatusControllerKilled {
		t.Errorf("sentry error")
	}
}

func TestSentryRun(t *testing.T) {
	var (
		s        = &sentry.Sentry{}
		trigger  = &TestingTrigger{}
		handler1 = &TestingHandler{}
		handler2 = &TestingHandler{
			WillPanic: true,
		}
	)
	s.Init([]sentry.ITrigger{trigger}, []sentry.IHandler{handler1, handler2})
	err := s.Run()
	if err == nil {
		t.Errorf("run error")
	}
	s.Ready()
	name := "test"
	event1 := trigger.NewEvent(name)
	event1.Payload = "123"
	trigger.Channel2 <- event1
	go func() {
		err := s.Run()
		if err != nil {
			panic(err)
		}
	}()
	event2 := <-handler1.Channel
	if event1.Name != event2.Name || event1.ID != event2.ID {
		t.Errorf("sentry run error")
	}
	if trigger.Counter != 1 {
		t.Errorf("trigger run error")
	}
	if handler1.Counter != 1 {
		t.Errorf("handler1 run error")
	}
	if event1.Channel != event2.Channel || event1.Payload != event2.Payload {
		t.Errorf("event error")
	}

	event3 := <-handler2.Channel
	if event1.Name != event3.Name || event1.ID != event3.ID {
		t.Errorf("sentry run error")
	}
	if trigger.Counter != 1 {
		t.Errorf("trigger run error")
	}
	if handler2.Counter != 1 {
		t.Errorf("handler2 run error")
	}
	if event1.Channel != event3.Channel || event1.Payload != event3.Payload {
		t.Errorf("event error")
	}
}
