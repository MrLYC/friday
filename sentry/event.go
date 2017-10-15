package sentry

import (
	"friday/config"
	"math/rand"
)

// Event : type for event
type Event struct {
	ID        string
	Channel   string
	Name      string
	Sender    string
	Type      string
	Payload   string
	RelatedTo string
}

// RefreshID : generate a new event id
func (e *Event) RefreshID() {
	configuration := config.Configuration.EventMETA
	idBuf := make([]byte, configuration.IDLength)
	rand.Read(idBuf)
	e.ID = string(idBuf)
}

// Copy : copy to a new event
func (e *Event) Copy() *Event {
	return &Event{
		ID:        e.ID,
		Channel:   e.Channel,
		Name:      e.Name,
		Sender:    e.Sender,
		Type:      e.Type,
		Payload:   e.Payload,
		RelatedTo: e.RelatedTo,
	}
}

// EventTemplate : template of event
var EventTemplate Event

// NewEvent : create an event from EventTemplate
func NewEvent() *Event {
	event := EventTemplate.Copy()
	event.RefreshID()
	return event
}
