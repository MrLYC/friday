package config

// Sentry : sentry configuration
type Sentry struct {
	ChannelBuffer int `yaml:"channel_buffer"`
}

// Init : init Sentry
func (c *Sentry) Init() {
	c.ChannelBuffer = 10
}
