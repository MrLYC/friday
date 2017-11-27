package command

import (
	"strings"
)

// ICommand : command interface
type ICommand interface {
	GetDescription() string
	SetFlags()
	Run() error
}

// TFlagStringArr :
type TFlagStringArr []string

// String :
func (i *TFlagStringArr) String() string {
	return strings.Join([]string(*i), ",")
}

// Set :
func (i *TFlagStringArr) Set(value string) error {
	*i = append(*i, value)
	return nil
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
