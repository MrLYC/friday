package sentry_test

import (
	"fmt"
	"friday/sentry"
	"testing"
)

type TestingHandler struct {
	sentry.BaseHandler
	Channel   chan sentry.Event
	Counter   int
	WillPanic bool
}

// Init :
func (h *TestingHandler) Init(s *sentry.Sentry) {
	h.BaseHandler.Name = "testing"
	h.Channel = make(chan sentry.Event, 10)
	h.BaseHandler.Init(s)
}

func (h *TestingHandler) Handle(event *sentry.Event) {
	h.Counter += 1
	ev := event.Copy()
	ev.ID = event.Name
	ev.Name = event.ID
	h.Channel <- *ev
	if h.WillPanic {
		panic(fmt.Errorf(h.GetName()))
	}
}

func TestHandlerInit(t *testing.T) {
	var (
		handler                  = TestingHandler{}
		ihandler sentry.IHandler = &handler
	)
	ihandler.Init(nil)

	if ihandler.GetName() != handler.BaseHandler.Name {
		t.Errorf("name error")
	}
}

func TestHandlerHandle(t *testing.T) {
	var (
		handler                  = TestingHandler{}
		ihandler sentry.IHandler = &handler
		event1                   = &sentry.Event{
			ID:   "1",
			Name: "2",
		}
	)
	ihandler.Init(nil)
	ihandler.Ready()
	ihandler.Handle(event1)
	event2 := <-handler.Channel
	ihandler.Terminate()
	if event1.ID != event2.Name || event1.Name != event2.ID {
		t.Errorf("handle error")
	}
}
