package config

import (
	"os"
)

// Timer : timer configuration
type Timer struct {
	CheckDuration string `yaml:"check_duration" validate:"regexp=^((\\d+(.\\d+)?)(h|m|s|ms|us|µs|ns))+$"`
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
