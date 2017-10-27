package config

import (
	"os"
)

// Timer : timer configuration
type Timer struct {
	CheckDuration string `yaml:"check_delta"`
}

// Init : init Timer
func (c *Timer) Init() {
	value := os.Getenv("FRIDAY_TIMER_CHECKDURATION")
	if value == "" {
		c.CheckDuration = "10s"
	} else {
		c.CheckDuration = value
	}
}
