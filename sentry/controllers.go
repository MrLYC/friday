package sentry

// IController :
type IController interface {
	GetName() string
	Ready()
	Terminate()
	Kill()
}

// ControllerStatus :
type ControllerStatus int

// Controller status constant
const (
	ControllerInitStatus        ControllerStatus = iota // 0
	ControllerReadyStatus       ControllerStatus = iota
	ControllerRuningStatus      ControllerStatus = iota
	ControllerTerminatingStatus ControllerStatus = iota
	ControllerTerminatedStatus  ControllerStatus = iota
	ControllerKilledStatus      ControllerStatus = iota
)

// BaseController :
type BaseController struct {
	Status ControllerStatus
}

// Ready :
func (c *BaseController) Ready() {
	c.Status = ControllerReadyStatus
}

// Terminate :
func (c *BaseController) Terminate() {
	c.Status = ControllerTerminatedStatus
}

// Kill :
func (c *BaseController) Kill() {
	c.Status = ControllerTerminatedStatus
}
