package sentry

// IHandler : trigger interface
type IHandler interface {
	Init(*Sentry)
	GetName() string
	Handle(Event)
	Start()
}

// BaseHandler :
type BaseHandler struct {
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

// Start :
func (h *BaseHandler) Start() {

}
