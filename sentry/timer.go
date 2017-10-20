package sentry

import (
	"friday/config"
	"friday/logging"
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
	GetStatus() DelayEventStatus
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

// GetStatus :
func (e *DelayEvent) GetStatus() DelayEventStatus {
	return e.Status
}

// Timer :
type Timer struct {
	BaseTrigger
	CheckDuration time.Duration
	ticker        *time.Ticker
	eventHeap     *binaryheap.Heap
	queueMux      sync.Mutex
}

// Init :
func (t *Timer) Init(s *Sentry) {
	t.eventHeap = binaryheap.NewWith(t.delayEventComparator)
	CheckDuration, err := time.ParseDuration(config.Configuration.Timer.CheckDuration)
	if err != nil {
		panic(err)
	}
	t.CheckDuration = CheckDuration
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

// Handle :
func (t *Timer) Handle(event *Event) {
	if event.Name == EventBroadcastNameQuit {
		t.Status = StatusControllerTerminating
	}
}

func (t *Timer) handleAbortedEvents(lowerLimitTime time.Time) {
	for delayEvent := t.PeekEvent(); delayEvent != nil; delayEvent = t.PeekEvent() {
		eventTime := delayEvent.GetTime()
		if eventTime.After(lowerLimitTime) {
			break
		}
		delayEvent = t.PopEvent()
		event := delayEvent.GetEvent()
		logging.Warningf(
			"Timer abort delay event: id=%v, name=%v, status: %v",
			event.ID, event.Name, delayEvent.GetStatus(),
		)
		delayEvent.StatusChanged(t, StatusDelayEventAbort)
	}
}

func (t *Timer) handleActivateEvents(highLimitTime time.Time) {
	for delayEvent := t.PeekEvent(); delayEvent != nil; delayEvent = t.PeekEvent() {
		eventTime := delayEvent.GetTime()
		if highLimitTime.Before(eventTime) {
			break
		}
		delayEvent = t.PopEvent()
		if delayEvent.GetStatus() != StatusDelayEventReady {
			event := delayEvent.GetEvent()
			logging.Warningf(
				"Timer abort delay event: id=%v, name=%v, status: %v",
				event.ID, event.Name, delayEvent.GetStatus(),
			)
			delayEvent.StatusChanged(t, StatusDelayEventAbort)
			continue
		}
		delayEvent.StatusChanged(t, StatusDelayEventPending)
		t.Channel <- delayEvent.GetEvent()
		delayEvent.StatusChanged(t, StatusDelayEventSent)
	}
}

// Run :
func (t *Timer) Run() {
	if t.Status != StatusControllerReady {
		panic(ErrNotReady)
	}
	t.Status = StatusControllerRuning
	broadcastChan := t.Sentry.DeclareChannel(ChanNameBroadcast)
	t.ticker = time.NewTicker(t.CheckDuration)

	for t.Status == StatusControllerRuning {
		select {
		case now := <-t.ticker.C:
			delayEvent := t.PeekEvent()
			if delayEvent == nil {
				continue
			}
			t.handleAbortedEvents(now)
			t.handleActivateEvents(now.Add(t.CheckDuration))
		case broadcastEv := <-broadcastChan:
			t.Handle(broadcastEv)
		}
	}
	t.Status = StatusControllerTerminated
}

// Terminate :
func (t *Timer) Terminate() {
	t.Status = StatusControllerTerminating
	t.ticker.Stop()
}
