package utils_test

import (
	"friday/utils"
	"sync"
	"testing"
)

func TestAtomicRun(t *testing.T) {
	var (
		a     utils.Atomic
		wg    sync.WaitGroup
		value = 0
	)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go a.Run(func() {
			val := value
			value = val + 1
			wg.Done()
		})
	}
	wg.Wait()
	if value != 100 {
		t.Errorf("atomic error")
	}
}
