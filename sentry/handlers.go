package sentry

// IHandler : trigger interface
type IHandler interface {
	IController
	Init(ISentry)
	Handle(*Event)
}

// BaseHandler :
type BaseHandler struct {
	BaseController
	Sentry ISentry
	Name   string
}

// Init :
func (h *BaseHandler) Init(sentry ISentry) {
	h.Sentry = sentry
}

// GetName :
func (h *BaseHandler) GetName() string {
	return h.Name
}
