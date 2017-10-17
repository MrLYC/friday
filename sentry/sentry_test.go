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
	)
	s.Init([]sentry.ITrigger{trigger}, []sentry.IHandler{handler})
	if trigger.Sentry != s || s.Triggers[trigger.GetName()] != trigger {
		t.Errorf("Trigger init error")
	}
	handlers := s.Handlers[handler.GetName()]
	hdr := handlers.Front().Value.(*TestingHandler)
	if handler.Sentry != s || hdr != handler {
		t.Errorf("Handler init error")
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
		s       = &sentry.Sentry{}
		trigger = &TestingTrigger{}
		handler = &TestingHandler{}
	)
	s.Init([]sentry.ITrigger{trigger}, []sentry.IHandler{handler})
	err := s.Run()
	if err == nil {
		t.Errorf("run error")
	}
	s.Ready()
	name := "test"
	event1 := trigger.NewEvent(name)
	trigger.Channel2 <- event1
	go s.Run()
	event2 := <-handler.Channel
	if event1.Name != event2.Name || event1.ID != event2.ID {
		t.Errorf("sentry run error")
	}
	if trigger.Counter != 1 {
		t.Errorf("trigger run error")
	}
	if handler.Counter != 1 {
		t.Errorf("handler run error")
	}
}
