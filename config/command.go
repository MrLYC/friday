package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
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
	data, err := yaml.Marshal(Configuration)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))
	return nil
}
