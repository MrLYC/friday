package firework

// Applet :
type Applet struct {
	BaseController
	Emitter IEmitter
	Status  ControllerStatus
}

// SetEmitter : set emitter
func (a *Applet) SetEmitter(emitter IEmitter) {
	a.Emitter = emitter
}
