package sentry

// Sentry :
type Sentry struct {
	BaseController
	Triggers map[string]ITrigger
	Handlers map[string]IHandler
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
	s.Handlers = make(map[string]IHandler, len(handlers))
	for _, handler := range handlers {
		handler.Init(s)
		s.Handlers[handler.GetName()] = handler
	}
}

// Ready :
func (s *Sentry) Ready() {
	for _, trigger := range s.Triggers {
		trigger.Ready()
	}
	for _, handler := range s.Handlers {
		handler.Ready()
	}
	s.BaseController.Ready()
}

// Terminate :
func (s *Sentry) Terminate() {
	for _, trigger := range s.Triggers {
		trigger.Terminate()
	}
	for _, handler := range s.Handlers {
		handler.Terminate()
	}
	s.BaseController.Terminate()
}

// Kill :
func (s *Sentry) Kill() {
	for _, trigger := range s.Triggers {
		trigger.Kill()
	}
	for _, handler := range s.Handlers {
		handler.Kill()
	}
	if s.Status != StatusControllerTerminated {
		s.Status = StatusControllerKilled
	}
}
