package config

// EventMETA : event meta configuration
type EventMETA struct {
	IDLength uint `yaml:"id_length"`
}

// Init : init EventMETA
func (e *EventMETA) Init() {
	e.IDLength = 32
}
