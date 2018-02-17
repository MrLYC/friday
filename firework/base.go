package firework

import (
	"crypto/rand"
	"friday/config"
)

//
var (
	ChanBroadcast = "*"
	ChanInternal  = "!"
)

// IFirework :
type IFirework interface {
	RefreshID()
	Copy() IFirework
	GetID() string
	GetChannel() string
	GetSender() string
	GetName() string
	GetType() string
	GetRelatedTo() string
	GetPayload() interface{}
	SetID(string)
	SetChannel(string)
	SetSender(string)
	SetName(string)
	SetType(string)
	SetRelatedTo(string)
	SetPayload(interface{})
}

// Firework : event type
type Firework struct {
	ID        string
	Channel   string
	Sender    string
	Name      string
	Type      string
	Payload   interface{}
	RelatedTo string
}

// RefreshID : generate a new event id
func (f *Firework) RefreshID() {
	configuration := config.Configuration.Event
	idBuf := make([]byte, configuration.IDLength)
	rand.Read(idBuf)
	f.ID = string(idBuf)
}

// Copy : copy to a new event
func (f *Firework) Copy() IFirework {
	return &Firework{
		ID:        f.ID,
		Channel:   f.Channel,
		Name:      f.Name,
		Type:      f.Type,
		Payload:   f.Payload,
		RelatedTo: f.RelatedTo,
	}
}

// GetID :
func (f *Firework) GetID() string {
	return f.ID
}

// GetChannel :
func (f *Firework) GetChannel() string {
	return f.Channel
}

// GetSender :
func (f *Firework) GetSender() string {
	return f.Sender
}

// GetName :
func (f *Firework) GetName() string {
	return f.Name
}

// GetType :
func (f *Firework) GetType() string {
	return f.Type
}

// GetRelatedTo :
func (f *Firework) GetRelatedTo() string {
	return f.RelatedTo
}

// GetPayload :
func (f *Firework) GetPayload() interface{} {
	return f.Payload
}

// SetID :
func (f *Firework) SetID(value string) {
	f.ID = value
}

// SetChannel :
func (f *Firework) SetChannel(value string) {
	f.Channel = value
}

// SetSender :
func (f *Firework) SetSender(value string) {
	f.Sender = value
}

// SetName :
func (f *Firework) SetName(value string) {
	f.Name = value
}

// SetType :
func (f *Firework) SetType(value string) {
	f.Type = value
}

// SetRelatedTo :
func (f *Firework) SetRelatedTo(value string) {
	f.RelatedTo = value
}

// SetPayload :
func (f *Firework) SetPayload(value interface{}) {
	f.Payload = value
}

// Handler :
type Handler func(IFirework)

// IApplet :
type IApplet interface {
	IController
	SetEmitter(IEmitter)
}

// IEmitter :
type IEmitter interface {
	IController
	AddApplet(IApplet) bool
	DeleteApplet(IApplet) bool
	DeclareChannel(string) chan IFirework
	On(string, string, Handler) (Handler, bool)
	Off(string, string, Handler) (Handler, bool)
	Fire(IFirework)
}
