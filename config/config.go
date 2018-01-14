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

	Cache     Cache     `yaml:"cache"`
	Database  Database  `yaml:"database"`
	Event     Event     `yaml:"event"`
	Firework  Firework  `yaml:"firework"`
	Logging   Logging   `yaml:"logging"`
	Migration Migration `yaml:"migration"`
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

	c.StrictInclude = false

	c.Cache.Init()
	c.Database.Init()
	c.Event.Init()
	c.Firework.Init()
	c.Logging.Init()
	c.Migration.Init()
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
	confPath := c.ConfigurationPath

	err = c.ReadFrom(confPath)
	if err != nil {
		return err
	}

	dirPath, _ := filepath.Split(confPath)
	includes := c.Includes
	strictInclude := c.StrictInclude

	for _, p := range includes {
		if !filepath.IsAbs(p) {
			p, err = filepath.Abs(filepath.Join(dirPath, p))
			if strictInclude && err != nil {
				return err
			}
		}
		err = c.ReadFrom(p)
		if strictInclude && err != nil {
			return err
		}
	}
	c.Includes = includes
	c.StrictInclude = strictInclude
	c.ConfigurationPath = confPath

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
