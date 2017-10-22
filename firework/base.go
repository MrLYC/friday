package firework

import (
	"crypto/rand"
	"friday/config"
)

// Firework : event type
type Firework struct {
	ID        string
	Channel   string
	Sender    string
	Name      string
	Type      string
	Payload   string
	RelatedTo string
}

// RefreshID : generate a new event id
func (e *Firework) RefreshID() {
	configuration := config.Configuration.Event
	idBuf := make([]byte, configuration.IDLength)
	rand.Read(idBuf)
	e.ID = string(idBuf)
}

// Copy : copy to a new event
func (e *Firework) Copy() *Firework {
	return &Firework{
		ID:        e.ID,
		Channel:   e.Channel,
		Name:      e.Name,
		Type:      e.Type,
		Payload:   e.Payload,
		RelatedTo: e.RelatedTo,
	}
}

// Chan :
type Chan chan *Firework

// ITrigger :
type ITrigger interface {
	IController
	SetEmitter(IEmitter)
	SetChannel(Chan)
}

// IHandler :
type IHandler interface {
	Handle(*Firework)
}

// IEmitter :
type IEmitter interface {
	IController
	On(string, IHandler) bool
	Off(string, IHandler) bool
	Fire(string, *Firework)
	AddTrigger(ITrigger)
	DelTrigger(ITrigger)
	DeclareChannel(string)
}
