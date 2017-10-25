package firework

// BaseApplet :
type BaseApplet struct {
	BaseController
	Emitter IEmitter
}

// SetEmitter :
func (a *BaseApplet) SetEmitter(emitter IEmitter) {
	a.Emitter = emitter
}
