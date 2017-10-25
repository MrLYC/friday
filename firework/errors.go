package firework

import "errors"

// Errors
var (
	ErrEmitterNotReady = errors.New("emitter not ready")
	ErrAppletNotReady  = errors.New("applet not ready")
)
