package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// IConfiguration : configuration interface
type IConfiguration interface {
	Init()
}

// ConfigurationType : configuration type
type ConfigurationType struct {
	Version   string
	EventMETA EventMETA `yaml:"event_meta"`
}

// Init : init ConfigurationType
func (c *ConfigurationType) Init() {
	c.Version = ConfVersion

	c.EventMETA.Init()
}

// ReadFrom : read configuration from path
func (c *ConfigurationType) ReadFrom(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, c)
	if err != nil {
		return err
	}
	if c.Version != ConfVersion {
		panic(fmt.Errorf("Unknown configuration version: %v", c.Version))
	}
	return nil
}

// Configuration : global configuration
var Configuration = ConfigurationType{}
