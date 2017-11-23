package config

import (
	"fmt"
	"io/ioutil"
	"os"

	validator "gopkg.in/validator.v2"
	yaml "gopkg.in/yaml.v2"
)

// IConfiguration : configuration interface
type IConfiguration interface {
	Init()
}

// ConfigurationType : configuration type
type ConfigurationType struct {
	Version           string
	ConfigurationPath string `yaml:"configuration_path"`

	Database Database `yaml:"database"`
	Event    Event    `yaml:"event"`
	Logging  Logging  `yaml:"logging"`
	Firework Firework `yaml:"firework"`
	Timer    Timer    `yaml:"timer"`
}

// Init : init ConfigurationType
func (c *ConfigurationType) Init() {
	c.Version = ConfVersion

	c.ConfigurationPath = os.Getenv("FRIDAY_CONFIG_PATH")
	if c.ConfigurationPath == "" {
		c.ConfigurationPath = "friday.yaml"
	}

	c.Database.Init()
	c.Event.Init()
	c.Logging.Init()
	c.Firework.Init()
	c.Timer.Init()
}

// ReadFrom : read configuration from path
func (c *ConfigurationType) ReadFrom(path string) error {
	c.ConfigurationPath = path
	return c.Read()
}

// Read : read configuration
func (c *ConfigurationType) Read() error {
	data, err := ioutil.ReadFile(c.ConfigurationPath)
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

// Validate :
func (c *ConfigurationType) Validate() error {
	return validator.Validate(c)
}

// Configuration : global configuration
var Configuration = ConfigurationType{}

func init() {
	Configuration.Init()
}
