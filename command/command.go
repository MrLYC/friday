package command

// ICommand : command interface
type ICommand interface {
	GetDescription() string
	SetFlags()
	Run() error
}

// BaseCommand : base command
type BaseCommand struct {
	Description string
}

// GetDescription : get command description
func (c *BaseCommand) GetDescription() string {
	return c.Description
}

// SetFlags : set parsing flags
func (c *BaseCommand) SetFlags() {

}

// Run : run command
func (c *BaseCommand) Run() error {
	return nil
}
