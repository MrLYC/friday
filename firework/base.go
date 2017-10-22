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

// Handler :
type Handler func(*Firework)

// IApplet :
type IApplet interface {
	IController
	SetEmitter(IEmitter)
}

// IEmitter :
type IEmitter interface {
	IController
	AddApplet(IApplet) bool
	DelApplet(IApplet) bool
	DeclareChannel(string) chan *Firework
	On(string, string, Handler) bool
	Off(string, string, Handler) bool
	Fire(string, *Firework)
}
