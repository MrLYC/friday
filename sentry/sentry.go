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
