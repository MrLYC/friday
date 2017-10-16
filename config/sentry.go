package config

// Sentry : sentry configuration
type Sentry struct {
	ChannelBuffer int `yaml:"channel_buffer"`
}

// Init : init Sentry
func (e *Sentry) Init() {
	e.ChannelBuffer = 10
}
