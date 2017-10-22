package firework

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
	StatusControllerInit        ControllerStatus = iota // 0
	StatusControllerReady       ControllerStatus = iota
	StatusControllerRuning      ControllerStatus = iota
	StatusControllerTerminating ControllerStatus = iota
	StatusControllerTerminated  ControllerStatus = iota
	StatusControllerKilled      ControllerStatus = iota
)

// BaseController :
type BaseController struct {
	Status ControllerStatus
}

// Ready :
func (c *BaseController) Ready() {
	c.Status = StatusControllerReady
}

// Terminate :
func (c *BaseController) Terminate() {
	c.Status = StatusControllerTerminated
}

// Kill :
func (c *BaseController) Kill() {
	c.Status = StatusControllerTerminated
}
