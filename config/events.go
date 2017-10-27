package config

// Event : event meta configuration
type Event struct {
	IDLength  uint   `yaml:"id_length" validate:"min=1"`
	EventType string `yaml:"event_type"`
}

// Init : init Event
func (e *Event) Init() {
	e.IDLength = 32
	e.EventType = "json"
}
