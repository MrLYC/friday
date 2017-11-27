package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

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
	Debug             bool
	ConfigurationPath string `yaml:"configuration_path"`

	StrictInclude bool     `yaml:"strict_include"`
	Includes      []string `yaml:"includes,omitempty"`

	Database  Database  `yaml:"database"`
	Migration Migration `yaml:"migration"`
	Event     Event     `yaml:"event"`
	Logging   Logging   `yaml:"logging"`
	Firework  Firework  `yaml:"firework"`
	Timer     Timer     `yaml:"timer"`
}

// Init : init ConfigurationType
func (c *ConfigurationType) Init() {
	c.Version = ConfVersion
	c.Debug = Mode == "debug"

	c.ConfigurationPath = os.Getenv("FRIDAY_CONFIG_PATH")
	if c.ConfigurationPath == "" {
		c.ConfigurationPath = "friday.yaml"
	}

	c.StrictInclude = true

	c.Database.Init()
	c.Migration.Init()
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
	return nil
}

// Read : read configuration
func (c *ConfigurationType) Read() error {
	var err error

	err = c.ReadFrom(c.ConfigurationPath)
	if err != nil {
		return err
	}

	dirPath, _ := filepath.Split(c.ConfigurationPath)

	for _, p := range c.Includes {
		if !filepath.IsAbs(p) {
			p, err = filepath.Abs(filepath.Join(dirPath, p))
			if c.StrictInclude && err != nil {
				return err
			}
		}
		err = c.ReadFrom(p)
		if c.StrictInclude && err != nil {
			return err
		}
	}

	if c.Version != ConfVersion {
		return fmt.Errorf("Unknown configuration version: %v", c.Version)
	}
	return nil
}

// Validate :
func (c *ConfigurationType) Validate() error {
	return validator.Validate(c)
}

// Dumps :
func (c *ConfigurationType) Dumps() (string, error) {
	data, err := yaml.Marshal(Configuration)
	return string(data), err
}

// Configuration : global configuration
var Configuration = ConfigurationType{}

func init() {
	Configuration.Init()
}
