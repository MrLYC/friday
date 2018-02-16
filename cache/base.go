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

	Expire(key string, duration time.Duration)
	Exists(key string) bool
	Delete(key string)

	StringSet(key string, value string) error
	StringGet(key string) (string, error)

	ListPush(key string, value string) error
	ListPop(key string) (string, error)
	ListRPush(key string, value string) error
	ListRPop(key string) (string, error)
	ListLen(key string) (int, error)

	TableSet(key string, field string, value string) error
	TableGet(key string, field string) (string, error)
	TableSetMappings(key string, mappings map[string]string) error
	TableGetAll(key string) (map[string]string, error)
}
