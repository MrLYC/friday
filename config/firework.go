package config

// Firework : firework configuration
type Firework struct {
	ChannelBuffer int      `yaml:"channel_buffer" validate:"min=1"`
	Applets       []string `yaml:"applets"`
}

// Init : init Firework
func (c *Firework) Init() {
	c.ChannelBuffer = 10
}
