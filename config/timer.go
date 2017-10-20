package config

// Timer : timer configuration
type Timer struct {
	CheckDuration string `yaml:"check_delta"`
}

// Init : init Timer
func (c *Timer) Init() {
	c.CheckDuration = "10s"
}
