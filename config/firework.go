package config

// Firework : firework configuration
type Firework struct {
	ChannelBuffer int `yaml:"channel_buffer" validate:"min=1"`
}

// Init : init Firework
func (c *Firework) Init() {
	c.ChannelBuffer = 10
}
