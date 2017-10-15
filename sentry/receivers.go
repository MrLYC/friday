package sentry

import (
	"container/list"
	"friday/logging"
	"friday/utils"
)

// Context : context for handler
type Context struct {
	Event *Event
}

// ReceiverHandler : hadler for receiver
type ReceiverHandler func(*Context)

// BaseReceiver : base receiver
type BaseReceiver struct {
	Channel chan Event
	Handers list.List
	Enable  bool
	Sentry  *Sentry
}

// Init :
func (r *BaseReceiver) Init(sentry *Sentry) {
	r.Sentry = sentry
}

// Start : start receiver
func (r *BaseReceiver) Start() error {
	if r.Channel == nil {
		return ErrReceiverNotReady
	}
	r.Enable = true
	go r.Run()
	return nil
}

// AddReceiverHandler : add receiver handler
func (r *BaseReceiver) AddReceiverHandler(handler ReceiverHandler) {
	r.Handers.PushBack(handler)
}

// CallHandler : call handler
func (r *BaseReceiver) CallHandler(handler ReceiverHandler, event *Event) {
	defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
		logging.Errorf(
			"event[%v](%v) handle error: %v",
			event.ID, event.Name, err,
		)
	})
	handler(&Context{
		Event: event.Copy(),
	})
}

// CallHandlers : call handlers
func (r *BaseReceiver) CallHandlers(event *Event) {
	for i := r.Handers.Front(); i != nil; i = i.Next() {
		go r.CallHandler(i.Value.(ReceiverHandler), event)
	}
}

// Run : poll events
func (r *BaseReceiver) Run() {
	for {
		if !r.Enable {
			return
		}
		event := <-r.Channel
		r.CallHandlers(&event)
	}
}
