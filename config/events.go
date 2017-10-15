package config

// EventMETA : event meta configuration
type EventMETA struct {
	IDLength  uint   `yaml:"id_length"`
	EventType string `yaml:"event_type"`
}

// Init : init EventMETA
func (e *EventMETA) Init() {
	e.IDLength = 32
	e.EventType = "json"
}
