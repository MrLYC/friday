package config

// Logging : logging configuration
type Logging struct {
	Level string `yaml:"level"`
}

// Init : init Logging
func (l *Logging) Init() {
	l.Level = "info"
}
