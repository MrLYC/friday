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
	fmt.Printf("BuildTag: %v\n", BuildTag)
	return nil
}

// ConfigurationCommand : configuration info
type ConfigurationCommand struct{}

// GetDescription : get command description
func (c *ConfigurationCommand) GetDescription() string {
	return "Print configuration"
}

// SetFlags : set parsing flags
func (c *ConfigurationCommand) SetFlags() {

}

// Run : run command
func (c *ConfigurationCommand) Run() error {
	data, err := Configuration.Dumps()
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
	return nil
}
