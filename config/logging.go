package config

import (
	"os"
)

// Logging : logging configuration
type Logging struct {
	Level string `yaml:"level"`
}

// Init : init Logging
func (l *Logging) Init() {

	value := os.Getenv("FRIDAY_LOGGING_LEVEL")
	if value == "" {
		l.Level = "info"
	} else {
		l.Level = value
	}
}
