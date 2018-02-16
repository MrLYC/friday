package cache

import (
	"errors"
	"time"
)

//
var (
	ErrItemNotFound   = errors.New("Item not found")
	ErrItemTypeError  = errors.New("Item type error")
	ErrItemValueError = errors.New("Item value error")
	ErrItemExpired    = errors.New("Item expired")
)

// ICache :
type ICache interface {
	Init()
	Close() error

	Expire(string, time.Duration)
	Delete(string)

	StringSet(string, string) error
	StringGet(string) (string, error)

	ListPush(string, string) error
	ListPop(string) (string, error)
	ListRPush(string, string) error
	ListRPop(string) (string, error)
	ListLen(string) (int, error)

	TableSet(string, string, string) error
	TableGet(string, string) (string, error)
	TableGetAll(string) (map[string]string, error)
}
