package sentry_test

import (
	"friday/sentry"
	"testing"
	"time"
)

func TestDelayEventHeap(t *testing.T) {
	var (
		timer = sentry.Timer{}
		e1    = &sentry.DelayEvent{
			Time: time.Date(2017, 10, 18, 0, 0, 0, 0, time.UTC),
		}
		e2 = &sentry.DelayEvent{
			Time: time.Date(2017, 10, 19, 0, 0, 0, 0, time.UTC),
		}
		e3 = &sentry.DelayEvent{
			Time: time.Date(2017, 10, 20, 0, 0, 0, 0, time.UTC),
		}
		e4 = &sentry.DelayEvent{
			Time: time.Date(2017, 10, 21, 0, 0, 0, 0, time.UTC),
		}
	)
	timer.Init(nil)
	timer.AddEvent(e3)
	timer.AddEvent(e1)
	timer.AddEvent(e4)
	timer.AddEvent(e2)

	if timer.PopEvent() != e1 {
		t.Errorf("e1 error")
	}

	if timer.PopEvent() != e2 {
		t.Errorf("e2 error")
	}

	if timer.PopEvent() != e3 {
		t.Errorf("e3 error")
	}

	if timer.PopEvent() != e4 {
		t.Errorf("e4 error")
	}
}
