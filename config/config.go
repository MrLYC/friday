package config

import (
	"fmt"
	"io/ioutil"

	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

// IConfiguration : configuration interface
type IConfiguration interface {
	Init()
}

// ConfigurationType : configuration type
type ConfigurationType struct {
	Version  string
	Event    Event    `yaml:"event"`
	Logging  Logging  `yaml:"logging"`
	Firework Firework `yaml:"firework"`
	Timer    Timer    `yaml:"timer"`
}

// Init : init ConfigurationType
func (c *ConfigurationType) Init() {
	c.Version = ConfVersion

	c.Event.Init()
	c.Logging.Init()
	c.Firework.Init()
	c.Timer.Init()
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

func (c *ConfigurationType) Validate() error {
	return validator.Validate(c)
}

// Configuration : global configuration
var Configuration = ConfigurationType{}

func init() {
	Configuration.Init()
}
