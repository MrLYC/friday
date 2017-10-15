package sentry

import (
	"container/list"
	"friday/config"
)

// ISender : sender interface
type ISender interface {
	Init(*Sentry)
	GetName() string
	NewEvent(name string, channel string) *Event
	Start() error
}

// BaseSender : base sender type
type BaseSender struct {
	Sentry        *Sentry
	EventTemplate *Event
}

// Init :
func (s *BaseSender) Init(sentry *Sentry) {
	s.Sentry = sentry
	s.EventTemplate = EventTemplate.Copy()
}

// NewEvent :
func (s *BaseSender) NewEvent(name string, channel string) *Event {
	event := s.EventTemplate.Copy()
	event.RefreshID()
	event.Name = name
	event.Channel = channel
	return event
}

// IReceiver : receiver interface
type IReceiver interface {
	Init(*Sentry)
	GetName() string
	Start() error
	Handle(*Event)
}

// Sentry :
type Sentry struct {
	Channels  map[string]chan Event
	Senders   list.List
	Receivers list.List
}

// Init :
func (s *Sentry) Init(channels []string) {
	conf := config.Configuration.Sentry
	s.Channels = make(map[string]chan Event, len(channels))
	for ch := range channels {
		s.Channels[channels[ch]] = make(chan Event, conf.ChannelBuffer)
	}
}

// AddSender : add sender
func (s *Sentry) AddSender(sender ISender) {
	sender.Init(s)
	s.Senders.PushBack(sender)
}

// AddReceiver : add receiver
func (s *Sentry) AddReceiver(receiver IReceiver) {
	receiver.Init(s)
	s.Receivers.PushBack(receiver)
}

// Start : start sentry
func (s *Sentry) Start() error {
	var err error
	for i := s.Receivers.Front(); i != nil; i = i.Next() {
		receiver := i.Value.(IReceiver)
		err = receiver.Start()
		if err != nil {
			return err
		}
	}
	for i := s.Senders.Front(); i != nil; i = i.Next() {
		sender := i.Value.(ISender)
		err = sender.Start()
		if err != nil {
			return err
		}
	}
	return nil
}
