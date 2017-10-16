package sentry

import (
	"friday/config"
)

// Init : init sentry
func Init() {
	conf := config.Configuration.Event
	EventTemplate = Event{
		Type: conf.EventType,
	}
}
