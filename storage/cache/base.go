package cache

import (
	"errors"
	"time"
)

//
var (
	ErrItemNotFound  = errors.New("Item not found")
	ErrItemTypeError = errors.New("Item type error")
	ErrItemExpired   = errors.New("Item expired")
)

// CacheItemIter :
type CacheItemIter func(string, interface{})

// ICache :
type ICache interface {
	Init()
	Close() error

	KeyExpire(string, time.Duration)
	ListExpire(string, time.Duration)
	TableExpire(string, time.Duration)

	SetKey(string, string) error
	GetKey(string) (string, error)
	DelKey(string) (error, bool)

	ListPush(string, string) error
	ListPop(string) (string, error)
	ListLen(string) (error, int)
	ListGet(string, int) (string, error)
	DelList(string) (error, bool)

	TableAdd(string, string) error
	TableGet(string, string) (string, error)
	TableGetAll(string) (error, map[string]string)
	DelTable(string) (error, bool)
}

func init() {
	ErrItemNotFound = errors.New("Item not found")
	ErrItemTypeError = errors.New("Item type error")
}
