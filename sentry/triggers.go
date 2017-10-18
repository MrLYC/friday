package sentry

// ITrigger : trigger interface
type ITrigger interface {
	IController
	Init(*Sentry)
	SetChannel(chan *Event)
	Run()
}

// BaseTrigger :
type BaseTrigger struct {
	BaseController
	Channel       chan *Event
	Name          string
	Sentry        *Sentry
	EventTemplate *Event
}

// Init :
func (t *BaseTrigger) Init(sentry *Sentry) {
	t.Sentry = sentry
	t.EventTemplate = EventTemplate.Copy()
	t.EventTemplate.Channel = t.GetName()
}

// SetChannel
func (t *BaseTrigger) SetChannel(channel chan *Event) {
	t.Channel = channel
}

// GetName :
func (t *BaseTrigger) GetName() string {
	return t.Name
}

// NewEvent :
func (t *BaseTrigger) NewEvent(name string) *Event {
	event := t.EventTemplate.Copy()
	event.RefreshID()
	event.Name = name
	return event
}
