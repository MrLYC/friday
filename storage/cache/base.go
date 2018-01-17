package cache

import (
	"errors"
	"time"
)

//
var (
	ErrItemNotFound  error
	ErrItemTypeError error
)

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
	ListPopString(string) (string, error)
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
