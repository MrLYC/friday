package cache

import (
	"time"
)

// ICache :
type ICache interface {
	Init()
	Close() error

	KeyExpire(string, time.Duration)
	ListExpire(string, time.Duration)
	TableExpire(string, time.Duration)

	SetKey(string, string) error
	GetKey(string) (error, string)
	DelKey(string) (error, bool)

	ListPush(string, string) error
	ListPopString(string) (error, string)
	ListLen(string) (error, int)
	ListGet(string, int) (error, string)
	DelList(string) (error, bool)

	TableAdd(string, string) error
	TableGet(string, string) (error, string)
	TableGetAll(string) (error, map[string]string)
	DelTable(string) (error, bool)
}
