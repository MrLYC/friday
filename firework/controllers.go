package firework

// IController :
type IController interface {
	SetName(string)
	SetStatus(ControllerStatus)
	GetName() string
	GetStatus() ControllerStatus
	Init()
	Ready()
	Run()
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
	Name   string
	Status ControllerStatus
}

// Init :
func (c *BaseController) Init() {
	c.SetStatus(StatusControllerInit)
}

// SetName :
func (c *BaseController) SetName(name string) {
	c.Name = name
}

// SetStatus :
func (c *BaseController) SetStatus(status ControllerStatus) {
	c.Status = status
}

// GetName :
func (c *BaseController) GetName() string {
	return c.Name
}

// GetStatus :
func (c *BaseController) GetStatus() ControllerStatus {
	return c.Status
}

// Ready :
func (c *BaseController) Ready() {
	c.Status = StatusControllerReady
}

// Run :
func (c *BaseController) Run() {
	c.Status = StatusControllerRuning
}

// Terminate :
func (c *BaseController) Terminate() {
	c.Status = StatusControllerTerminated
}

// Kill :
func (c *BaseController) Kill() {
	c.Status = StatusControllerTerminated
}
