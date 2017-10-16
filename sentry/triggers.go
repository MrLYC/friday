package sentry

import "friday/config"

// ITrigger : trigger interface
type ITrigger interface {
	Init(*Sentry)
	GetName() string
	GetChannel() chan *Event
	Run()
	Start()
}

// BaseTrigger :
type BaseTrigger struct {
	Channel       chan *Event
	Name          string
	Sentry        *Sentry
	EventTemplate *Event
}

// Init :
func (t *BaseTrigger) Init(sentry *Sentry) {
	conf := config.Configuration.Sentry
	t.Sentry = sentry
	t.EventTemplate = EventTemplate.Copy()
	t.EventTemplate.Channel = t.GetName()
	t.Channel = make(chan *Event, conf.ChannelBuffer)
}

// GetName :
func (t *BaseTrigger) GetName() string {
	return t.Name
}

// GetChannel :
func (t *BaseTrigger) GetChannel() chan *Event {
	return t.Channel
}

// Start :
func (t *BaseTrigger) Start() {

}

// NewEvent :
func (t *BaseTrigger) NewEvent(name string) *Event {
	event := t.EventTemplate.Copy()
	event.RefreshID()
	event.Name = name
	return event
}
