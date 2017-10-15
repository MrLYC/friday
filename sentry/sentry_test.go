package sentry_test

import (
	"fmt"
	"friday/sentry"
	"testing"
)

type TestSenderType struct {
	sentry.BaseSender
	Name   string
	Number int
	Error  error
}

func (s *TestSenderType) GetName() string {
	return s.Name
}

func (s *TestSenderType) Start() error {
	s.Number += 1
	return s.Error
}

func TestSentryMakeChannels(t *testing.T) {
	s := sentry.Sentry{}
	s.Init([]string{"test"})
	_, ok := s.Channels["test"]
	if !ok {
		t.Errorf("channel init failded")
	}
}

func TestSentryStart(t *testing.T) {
	s := sentry.Sentry{}
	sender := &TestSenderType{}
	s.Init([]string{})
	s.AddSender(sender)
	err := s.Start()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	if sender.Number != 1 {
		t.Errorf("sender start error")
	}
}

func TestSentryStartError(t *testing.T) {
	serr := fmt.Errorf("test")
	s := sentry.Sentry{}
	sender := &TestSenderType{
		Error: serr,
	}
	s.Init([]string{})
	s.AddSender(sender)
	err := s.Start()
	if err != serr {
		t.Errorf("error: %v", err)
	}
}

func TestTestSenderType(t *testing.T) {
	var (
		sender  = "sender"
		name    = "name"
		channel = "channel"
		s       sentry.ISender
	)
	s = &TestSenderType{
		Name: sender,
	}
	s.Init(nil)

	if sender != s.GetName() {
		t.Errorf("sender name error")
	}

	event := s.NewEvent(name, channel)

	if event.Sender != "" {
		t.Errorf("event sender error")
	}
	if event.Name != name {
		t.Errorf("event name error")
	}
	if event.Channel != channel {
		t.Errorf("event channel error")
	}
}
