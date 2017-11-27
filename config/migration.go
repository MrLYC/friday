package config

// Migration : megration configuration
type Migration struct {
	LogMode bool `yaml:"log_mode"`
}

// Init : init Migration
func (e *Migration) Init() {
	e.LogMode = true
}
