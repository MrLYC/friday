package sentry

import (
	"friday/config"
)

// Init : init sentry
func Init() {
	conf := config.Configuration.EventMETA
	EventTemplate = Event{
		Type: conf.EventType,
	}
}
