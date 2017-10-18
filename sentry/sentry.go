package sentry

import (
	"friday/config"
	"friday/logging"
	"friday/utils"
	"reflect"
	"time"
)

// Channel names
var (
	ChanNameBroadcast = "*"
	ChanNameInternal  = "-"
)

// Sentry :
type Sentry struct {
	BaseController
	RunAt    time.Time
	Channels map[string]chan *Event
	Triggers map[string][]ITrigger
	Handlers map[string][]IHandler
}

// GetName :
func (s *Sentry) GetName() string {
	return "Sentry"
}

// Init :
func (s *Sentry) Init(triggers []ITrigger, handlers []IHandler) {
	s.Channels = make(map[string]chan *Event)
	s.DeclareChannel(ChanNameBroadcast)
	s.DeclareChannel(ChanNameInternal)

	s.Triggers = make(map[string][]ITrigger)
	for _, trigger := range triggers {
		s.AddTrigger(trigger)
	}

	s.Handlers = make(map[string][]IHandler)
	for _, handler := range handlers {
		s.AddHandler(handler)
	}
}

// DeclareChannel :
func (s *Sentry) DeclareChannel(name string) chan *Event {
	channel, ok := s.Channels[name]
	if !ok {
		channel = make(chan *Event, config.Configuration.Sentry.ChannelBuffer)
		s.Channels[name] = channel
	}
	return channel
}

// AddTrigger :
func (s *Sentry) AddTrigger(trigger ITrigger) {
	trigger.Init(s)
	name := trigger.GetName()
	channel := s.DeclareChannel(name)
	trigger.SetChannel(channel)

	triggers, ok := s.Triggers[name]
	if !ok {
		triggers = make([]ITrigger, 0, 1)
	}
	s.Triggers[name] = append(triggers, trigger)
}

// AddHandler :
func (s *Sentry) AddHandler(handler IHandler) {
	handler.Init(s)
	name := handler.GetName()
	handlers, ok := s.Handlers[name]
	if !ok {
		handlers = make([]IHandler, 0, 1)
	}
	s.Handlers[name] = append(handlers, handler)
}

// EachTrigger :
func (s *Sentry) EachTrigger(f func(ITrigger)) {
	for _, triggers := range s.Triggers {
		for _, trigger := range triggers {
			f(trigger)
		}
	}
}

// EachHandler :
func (s *Sentry) EachHandler(f func(IHandler)) {
	for _, handlers := range s.Handlers {
		for _, handler := range handlers {
			f(handler)
		}
	}
}

// Ready :
func (s *Sentry) Ready() {
	s.EachTrigger(func(trigger ITrigger) {
		trigger.Ready()
	})
	s.EachHandler(func(handler IHandler) {
		handler.Ready()
	})
	s.BaseController.Ready()
}

// Terminate :
func (s *Sentry) Terminate() {
	s.EachTrigger(func(trigger ITrigger) {
		trigger.Terminate()
	})
	s.EachHandler(func(handler IHandler) {
		handler.Terminate()
	})
	s.BaseController.Terminate()
}

// Kill :
func (s *Sentry) Kill() {
	s.EachTrigger(func(trigger ITrigger) {
		trigger.Kill()
	})
	s.EachHandler(func(handler IHandler) {
		handler.Kill()
	})
	if s.Status != StatusControllerTerminated {
		s.Status = StatusControllerKilled
	}
}

// Run :
func (s *Sentry) Run() error {
	if s.Status != StatusControllerReady {
		return ErrNotReady
	}
	s.Status = StatusControllerRuning
	s.RunAt = time.Now()
	logging.Infof("Sentry run at: %s", s.RunAt.String())

	s.EachTrigger(func(trigger ITrigger) {
		go func() {
			defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
				logging.Errorf("Trigger[%v] error: %s", trigger.GetName(), err)
			})
			trigger.Run()
		}()
	})

	channels := s.Channels
	cases := make([]reflect.SelectCase, 0, len(channels))
	for _, channel := range channels {
		cases = append(cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(channel),
		})
	}

	for {
		if s.Status != StatusControllerRuning {
			break
		}
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			logging.Warningf("Channel[%v] error", chosen)
			continue
		}
		event := value.Interface().(*Event)
		handlers, ok := s.Handlers[event.Channel]
		if !ok {
			continue
		}
		for _, handler := range handlers {
			go func(handler IHandler, ev *Event) {
				defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
					logging.Errorf("Handler[%v] error: %s", handler.GetName(), err)
				})
				handler.Handle(ev)
			}(handler, event.Copy())
		}
	}
	return nil
}
