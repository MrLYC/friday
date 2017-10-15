package sentry

import (
	"errors"
)

// Errors
var (
	ErrUnknownChannel   = errors.New("unknown channel")
	ErrSenderNotReady   = errors.New("sender not ready")
	ErrReceiverNotReady = errors.New("receiver not ready")
)
