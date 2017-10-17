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
	if handler.Sentry != s || s.Handlers[handler.GetName()] != handler {
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
