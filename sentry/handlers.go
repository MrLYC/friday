package sentry

// IHandler : trigger interface
type IHandler interface {
	IController
	Init(*Sentry)
	Handle(Event)
}

// BaseHandler :
type BaseHandler struct {
	BaseController
	Sentry *Sentry
	Name   string
}

// Init :
func (h *BaseHandler) Init(sentry *Sentry) {
	h.Sentry = sentry
}

// GetName :
func (h *BaseHandler) GetName() string {
	return h.Name
}
