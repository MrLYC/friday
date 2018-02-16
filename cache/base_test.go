package cache_test

import (
	"friday/cache"
	"friday/cache/memcache"
	"testing"
	"time"
)

func testICacheString(name string, icache cache.ICache, t *testing.T) {
	var (
		err   error
		value string
		key   = "tom"
	)

	if icache.Exists(key) {
		t.Errorf("icache[%s] Exists error: %v", name, key)
	}

	err = icache.StringSet(key, "1")
	if err != nil {
		t.Error("icache[%s] StringSet error: %v", name, err)
	}

	value, err = icache.StringGet(key)
	if value != "1" || err != nil {
		t.Errorf("icache[%s] StringGet error: %v, %v", name, value, err)
	}

	icache.Delete(key)

	value, err = icache.StringGet(key)
	if value != "" || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] StringGet error: %v, %v", name, value, err)
	}
}

func testICacheList(name string, icache cache.ICache, t *testing.T) {
	var (
		err    error
		value  string
		length int
		key    = "bill"
	)

	if icache.Exists(key) {
		t.Errorf("icache[%s] Exists error: %v", name, key)
	}

	value, err = icache.ListPop(key)
	if value != "" || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] ListPop error: %v, %v", name, value, err)
	}

	value, err = icache.ListRPop(key)
	if value != "" || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] ListRPop error: %v, %v", name, value, err)
	}

	length, err = icache.ListLen(key)
	if length != 0 || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] ListLen error: %v, %v", name, length, err)
	}

	icache.ListPush(key, "")

	err = icache.ListPush(key, "1")
	if err != nil {
		t.Errorf("icache[%s] ListPush error: %v", name, err)
	}

	err = icache.ListRPush(key, "2")
	if err != nil {
		t.Errorf("icache[%s] ListPush error: %v", name, err)
	}

	length, err = icache.ListLen(key)
	if length != 3 || err != nil {
		t.Errorf("icache[%s] ListLen error: %v, %v", name, length, err)
	}

	value, err = icache.ListPop(key)
	if value != "1" || err != nil {
		t.Errorf("icache[%s] ListPop error: %v, %v", name, value, err)
	}

	value, err = icache.ListRPop(key)
	if value != "2" || err != nil {
		t.Errorf("icache[%s] ListRPop error: %v, %v", name, value, err)
	}

	icache.Delete(key)

	length, err = icache.ListLen(key)
	if length != 0 || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] ListLen error: %v, %v", name, length, err)
	}
}

func testICacheTable(name string, icache cache.ICache, t *testing.T) {
	var (
		err    error
		value  string
		values map[string]string
		key    = "joe"
	)

	if icache.Exists(key) {
		t.Errorf("icache[%s] Exists error: %v", name, key)
	}

	values, err = icache.TableGetAll(key)
	if len(values) != 0 || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] TableGetAll error: %v, %v", name, values, err)
	}

	err = icache.TableSet(key, "a", "1")
	if err != nil {
		t.Errorf("icache[%s] TableSet error: %v", name, err)
	}

	value, err = icache.TableGet(key, "a")
	if err != nil {
		t.Errorf("icache[%s] TableGet error: %v, %v", name, value, err)
	}

	err = icache.TableSetMappings(key, map[string]string{
		"b": "2",
		"c": "3",
	})
	if err != nil {
		t.Errorf("icache[%s] TableSetMappings error: %v", name, err)
	}

	values, err = icache.TableGetAll(key)
	if len(values) != 3 || err != nil {
		t.Errorf("icache[%s] TableGetAll error: %v, %v", name, values, err)
	}
	if values["a"] != "1" || values["b"] != "2" || values["c"] != "3" {
		t.Errorf("icache[%s] TableGetAll error: %v, %v", name, values, err)
	}

	icache.Delete(key)

	values, err = icache.TableGetAll(key)
	if len(values) != 0 || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] TableGetAll error: %v, %v", name, values, err)
	}
}

func testICacheExpire(name string, icache cache.ICache, t *testing.T) {
	var (
		err   error
		value string
		key   string
	)

	duration, _ := time.ParseDuration("-10h")

	key = "jack"
	err = icache.StringSet(key, "1")
	if err != nil {
		t.Error("icache[%s] StringSet error: %v", name, err)
	}
	icache.Expire(key, duration)
	if icache.Exists(key) {
		t.Errorf("icache[%s] Exists error: %v", name, key)
	}

	value, err = icache.StringGet(key)
	if value != "" || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] StringGet error: %v, %v", name, value, err)
	}

	key = "tony"
	err = icache.ListPush(key, "1")
	if err != nil {
		t.Error("icache[%s] ListPush error: %v", name, err)
	}

	icache.Expire(key, duration)
	if icache.Exists(key) {
		t.Errorf("icache[%s] Exists error: %v", name, key)
	}

	value, err = icache.ListPop(key)
	if value != "" || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] ListPop error: %v, %v", name, value, err)
	}

	key = "mary"
	err = icache.TableSet(key, "a", "1")
	if err != nil {
		t.Error("icache[%s] TableSet error: %v", name, err)
	}

	icache.Expire(key, duration)
	if icache.Exists(key) {
		t.Errorf("icache[%s] Exists error: %v", name, key)
	}

	value, err = icache.TableGet(key, "a")
	if value != "" || err != cache.ErrItemNotFound {
		t.Errorf("icache[%s] TableGet error: %v, %v", name, value, err)
	}
}

func testICacheOverwrite(name string, icache cache.ICache, t *testing.T) {
	type setter func(key string) error
	var (
		err     error
		setters = map[string]setter{
			"string": func(key string) error {
				return icache.StringSet(key, key)
			},
			"list": func(key string) error {
				return icache.ListPush(key, key)
			},
			"table": func(key string) error {
				return icache.TableSet(key, key, key)
			},
		}
	)
	for k1, f1 := range setters {
		err = f1(k1)
		if err != nil {
			t.Errorf("icache[%s] set error: %v, %v", name, k1, err)
		}
		for k2, f2 := range setters {
			err = f2(k1)
			if k1 != k2 && err != cache.ErrItemTypeError {
				t.Errorf("icache[%s] overwrite[%s] error: %v, %v", name, k1, k2, err)
			} else if k1 == k2 && err != nil {
				t.Errorf("icache[%s] rewrite[%s] error: %v, %v", name, k1, k2, err)
			}
		}
		icache.Delete(k1)
	}
}

func testICacheExpiredOverwrite(name string, icache cache.ICache, t *testing.T) {
	type setter func(key string) error
	var (
		err     error
		setters = map[string]setter{
			"string": func(key string) error {
				return icache.StringSet(key, key)
			},
			"list": func(key string) error {
				return icache.ListPush(key, key)
			},
			"table": func(key string) error {
				return icache.TableSet(key, key, key)
			},
		}
	)

	duration, _ := time.ParseDuration("-10h")

	for k1, f1 := range setters {
		for k2, f2 := range setters {
			err = f1(k1)
			if err != nil {
				t.Errorf("icache[%s] set error: %v, %v", name, k1, err)
			}

			icache.Expire(k1, duration)
			if icache.Exists(k1) {
				t.Errorf("icache[%s] Exists error: %v", name, k1)
			}

			err = f2(k1)
			if k1 != k2 && err != nil {
				t.Errorf("icache[%s] overwrite[%s] error: %v, %v", name, k1, k2, err)
			} else if k1 == k2 && err != nil {
				t.Errorf("icache[%s] rewrite[%s] error: %v, %v", name, k1, k2, err)
			}

			icache.Delete(k1)
		}
	}
}

func TestICache(t *testing.T) {
	var (
		cacheImplements = map[string]cache.ICache{
			"mem": memcache.NewCache(),
		}
		err error
	)
	for name, icache := range cacheImplements {
		icache.Init()

		testICacheString(name, icache, t)
		testICacheList(name, icache, t)
		testICacheTable(name, icache, t)
		testICacheExpire(name, icache, t)
		testICacheOverwrite(name, icache, t)
		testICacheExpiredOverwrite(name, icache, t)

		err = icache.Close()
		if err != nil {
			t.Errorf("icache[%s] Close error: %v", name, err)
		}
	}
}
