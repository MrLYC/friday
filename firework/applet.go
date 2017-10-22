package firework

// Applet :
type Applet struct {
	BaseController
	Emitter IEmitter
	Status  ControllerStatus
	Channel Chan
}

// SetEmitter : set emitter
func (a *Applet) SetEmitter(emitter IEmitter) {
	a.Emitter = emitter
}

// SetChannel : set control channel
func (a *Applet) SetChannel(channel Chan) {
	a.Channel = channel
}
