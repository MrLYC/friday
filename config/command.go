package config

import (
	"fmt"
)

// VersionCommand : version info
type VersionCommand struct{}

// GetDescription : get command description
func (c *VersionCommand) GetDescription() string {
	return "Print version infomations"
}

// SetFlags : set parsing flags
func (c *VersionCommand) SetFlags() {

}

// Run : run command
func (c *VersionCommand) Run() error {
	fmt.Printf("Version: %v\n", Version)
	fmt.Printf("Mode: %v\n", Mode)
	return nil
}
