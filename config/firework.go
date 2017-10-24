package config

// Firework : firework configuration
type Firework struct {
	ChannelBuffer int `yaml:"channel_buffer"`
}

// Init : init Firework
func (c *Firework) Init() {
	c.ChannelBuffer = 10
}
