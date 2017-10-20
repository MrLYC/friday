package sentry

// ITrigger : trigger interface
type ITrigger interface {
	IController
	Init(ISentry)
	SetControlChannel(chan *Event)
	Run()
}

// BaseTrigger :
type BaseTrigger struct {
	BaseController
	ControlChannel chan *Event
	Channel        chan *Event
	Name           string
	Sentry         ISentry
	EventTemplate  *Event
}

// Init :
func (t *BaseTrigger) Init(sentry ISentry) {
	name := t.GetName()
	t.Sentry = sentry
	t.EventTemplate = EventTemplate.Copy()
	t.EventTemplate.Channel = name
	t.Channel = sentry.DeclareChannel(name)
}

// SetControlChannel :
func (t *BaseTrigger) SetControlChannel(channel chan *Event) {
	t.ControlChannel = channel
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
