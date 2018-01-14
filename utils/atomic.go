package utils

import (
	"sync"
)

// AtomicFunc :
type AtomicFunc func()

// Atomic :
type Atomic struct {
	mutex sync.Mutex
}

// Run :
func (a *Atomic) Run(f AtomicFunc) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	f()
}
