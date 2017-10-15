package sentry

import (
	"errors"
)

// Errors
var (
	ErrSenderNotReady   = errors.New("sender not ready")
	ErrReceiverNotReady = errors.New("receiver not ready")
)
