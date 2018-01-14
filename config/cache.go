package config

// Cache :
type Cache struct {
	CheckDuration string `yaml:"check_duration" validate:"regexp=^((\\d+(.\\d+)?)(h|m|s|ms|us|µs|ns))+$"`
}

// Init : init Cache
func (c *Cache) Init() {
	c.CheckDuration = "10s"
}
