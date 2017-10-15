package sentry_test

import (
	"friday/sentry"
	"testing"
)

func TestBaseReceiverEmptyChannel(t *testing.T) {
	r := sentry.BaseReceiver{}
	if r.Channel != nil {
		t.Errorf("channel not nil")
	}
	err := r.Start()
	if err == nil {
		t.Errorf("reciver start with empty channel")
	}
}
