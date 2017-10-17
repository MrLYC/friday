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
