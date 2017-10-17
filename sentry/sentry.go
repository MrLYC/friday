package sentry

import (
	"container/list"
	"friday/logging"
	"friday/utils"
	"reflect"
	"time"
)

// Sentry :
type Sentry struct {
	BaseController
	RunAt    time.Time
	Triggers map[string]ITrigger
	Handlers map[string]*list.List
}

// GetName :
func (s *Sentry) GetName() string {
	return "Sentry"
}

// Init :
func (s *Sentry) Init(triggers []ITrigger, handlers []IHandler) {
	s.Triggers = make(map[string]ITrigger, len(triggers))
	for _, trigger := range triggers {
		trigger.Init(s)
		s.Triggers[trigger.GetName()] = trigger
	}
	s.Handlers = make(map[string]*list.List, len(handlers))
	for _, handler := range handlers {
		handler.Init(s)
		name := handler.GetName()
		handlerList, ok := s.Handlers[name]
		if !ok {
			handlerList = list.New()
			s.Handlers[name] = handlerList
		}
		handlerList.PushBack(handler)
	}
}

// Ready :
func (s *Sentry) Ready() {
	for _, trigger := range s.Triggers {
		trigger.Ready()
	}
	for _, handlers := range s.Handlers {
		for i := handlers.Front(); i != nil; i = i.Next() {
			handler := i.Value.(IHandler)
			handler.Ready()
		}
	}
	s.BaseController.Ready()
}

// Terminate :
func (s *Sentry) Terminate() {
	for _, trigger := range s.Triggers {
		trigger.Terminate()
	}
	for _, handlers := range s.Handlers {
		for i := handlers.Front(); i != nil; i = i.Next() {
			handler := i.Value.(IHandler)
			handler.Terminate()
		}
	}
	s.BaseController.Terminate()
}

// Kill :
func (s *Sentry) Kill() {
	for _, trigger := range s.Triggers {
		trigger.Kill()
	}
	for _, handlers := range s.Handlers {
		for i := handlers.Front(); i != nil; i = i.Next() {
			handler := i.Value.(IHandler)
			handler.Kill()
		}
	}
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

	triggers := s.Triggers
	cases := make([]reflect.SelectCase, len(s.Triggers))
	i := 0
	for _, trigger := range triggers {
		go trigger.Run()
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(trigger.GetChannel()),
		}
		i = i + 1
	}

	for {
		if s.Status != StatusControllerRuning {
			break
		}
		chosen, value, ok := reflect.Select(cases)
		if !ok {
			name := ""
			i := 0
			for _, trigger := range s.Triggers {
				if i == chosen {
					name = trigger.GetName()
					break
				}
			}
			logging.Warningf("Channel[%v] error", name)
			continue
		}
		event := value.Interface().(*Event)
		handlers, ok := s.Handlers[event.Channel]
		if !ok {
			continue
		}
		for i := handlers.Front(); i != nil; i = i.Next() {
			handler := i.Value.(IHandler)
			go func(ev *Event) {
				utils.ErrorRecoverCall(func(err *utils.TraceableError) {
					logging.Errorf("Handler[%v] error: %v", handler.GetName(), err)
				})
				handler.Handle(ev)
			}(event.Copy())
		}
	}
	return nil
}
