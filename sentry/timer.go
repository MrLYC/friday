package sentry

import (
	"sync"
	"time"

	"github.com/emirpasic/gods/trees/binaryheap"
)

// DelayEventStatus :
type DelayEventStatus int

// DelayEventStatus
const (
	StatusDelayEventAbort   DelayEventStatus = iota
	StatusDelayEventReady   DelayEventStatus = iota
	StatusDelayEventPending DelayEventStatus = iota
	StatusDelayEventSent    DelayEventStatus = iota
)

// IDelayEvent :
type IDelayEvent interface {
	GetEvent() *Event
	GetTime() time.Time
	StatusChanged(*Timer, DelayEventStatus)
}

// DelayEvent :
type DelayEvent struct {
	Event  *Event
	Time   time.Time
	Status DelayEventStatus
}

// GetEvent :
func (e *DelayEvent) GetEvent() *Event {
	return e.Event
}

// GetTime :
func (e *DelayEvent) GetTime() time.Time {
	return e.Time
}

// StatusChanged :
func (e *DelayEvent) StatusChanged(timer *Timer, status DelayEventStatus) {
	e.Status = status
}

// Timer :
type Timer struct {
	BaseTrigger
	eventHeap *binaryheap.Heap
	queueMux  sync.Mutex
}

// Init :
func (t *Timer) Init(s *Sentry) {
	t.eventHeap = binaryheap.NewWith(t.delayEventComparator)
	t.BaseTrigger.Init(s)
}

// delayEventComparator
func (t *Timer) delayEventComparator(i1 interface{}, i2 interface{}) int {
	t1 := i1.(IDelayEvent).GetTime()
	t2 := i2.(IDelayEvent).GetTime()
	if t1.After(t2) {
		return 1
	} else if t1.Equal(t2) {
		return 0
	} else {
		return -1
	}
}

// AddEvent :
func (t *Timer) AddEvent(event IDelayEvent) {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	t.eventHeap.Push(event)
}

// PopEvent :
func (t *Timer) PopEvent() IDelayEvent {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	ev, _ := t.eventHeap.Pop()
	return ev.(IDelayEvent)
}

// PeekEvent :
func (t *Timer) PeekEvent() IDelayEvent {
	ev, _ := t.eventHeap.Peek()
	return ev.(IDelayEvent)
}
