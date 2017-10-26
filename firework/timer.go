package firework

import (
	"friday/config"
	"friday/logging"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/emirpasic/gods/trees/binaryheap"
)

// DelayFireworkStatus :
type DelayFireworkStatus int

// DelayFireworkStatus
const (
	StatusDelayFireworkAbort   DelayFireworkStatus = iota
	StatusDelayFireworkReady   DelayFireworkStatus = iota
	StatusDelayFireworkPending DelayFireworkStatus = iota
	StatusDelayFireworkSent    DelayFireworkStatus = iota
)

// Timer constants :
const (
	TimerChannelName     = "TIMER"
	TimerFireworkDelay   = "Delay"
	TimerFireworkTimesUp = "TimesUp"
	TimerFireworkAbort   = "TimerAbort"
)

// DelayFirework :
type DelayFirework struct {
	Firework *Firework
	Time     time.Time
	Status   DelayFireworkStatus
}

// GetFirework :
func (f *DelayFirework) GetFirework() *Firework {
	return f.Firework
}

func delayFireworkComparator(i1 interface{}, i2 interface{}) int {
	t1 := i1.(*DelayFirework).Time
	t2 := i2.(*DelayFirework).Time
	if t1.After(t2) {
		return 1
	} else if t1.Equal(t2) {
		return 0
	} else {
		return -1
	}
}

// Timer :
type Timer struct {
	BaseApplet
	CheckDuration time.Duration
	ticker        *time.Ticker
	heap          *binaryheap.Heap
	queueMux      sync.Mutex
}

// Init :
func (t *Timer) Init() {
	t.heap = binaryheap.NewWith(delayFireworkComparator)
	CheckDuration, err := time.ParseDuration(config.Configuration.Timer.CheckDuration)
	if err != nil {
		panic(err)
	}
	t.CheckDuration = CheckDuration
}

// Add :
func (t *Timer) Add(firework *DelayFirework) {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	firework.Status = StatusDelayFireworkReady
	t.heap.Push(firework)
}

// Pop :
func (t *Timer) Pop() *DelayFirework {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	f, ok := t.heap.Pop()
	if !ok {
		return nil
	}
	return f.(*DelayFirework)
}

// Peek :
func (t *Timer) Peek() *DelayFirework {
	t.queueMux.Lock()
	defer t.queueMux.Unlock()
	f, ok := t.heap.Peek()
	if !ok {
		return nil
	}
	return f.(*DelayFirework)
}

// Ready :
func (t *Timer) Ready() {
	t.Emitter.On(TimerChannelName, TimerFireworkDelay, func(firework *Firework) {
		f := firework.Payload.(DelayFirework)
		t.Add(&f)
	})
	t.BaseApplet.Ready()
}

// Terminate :
func (t *Timer) Terminate() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
	t.BaseApplet.Terminate()
}

// Run :
func (t *Timer) Run() {
	if t.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}
	t.Status = StatusControllerRuning
	t.ticker = time.NewTicker(t.CheckDuration)
	for t.Status == StatusControllerRuning {
		now, ok := <-t.ticker.C
		if !ok {
			break
		}
		f := t.Peek()
		if f == nil {
			continue
		}

		abortFireworks := arraylist.New()
		activateFireworks := arraylist.New()
		lowLimitTime := now.Add(-t.CheckDuration)
		highLimitTime := now.Add(t.CheckDuration)

		for f := t.Peek(); f != nil; f = t.Peek() {
			eventTime := f.Time
			if eventTime.After(highLimitTime) {
				break
			}
			f = t.Pop()
			if eventTime.Before(lowLimitTime) || f.Status != StatusDelayFireworkReady {
				abortFireworks.Add(f)
			} else {
				activateFireworks.Add(f)
			}
		}
		go activateFireworks.Each(func(i int, df interface{}) {
			t.Emitter.Fire(df.(*DelayFirework).GetFirework())
		})
		go abortFireworks.Each(func(i int, item interface{}) {
			df := item.(*DelayFirework)
			f := df.GetFirework()
			logging.Warningf(
				"Timer abort at %v-%v: Channel[%s], Name[%s], Time[%v]",
				lowLimitTime, highLimitTime, f.Channel, f.Name, df.Time,
			)
			f.Name = TimerFireworkAbort
			t.Emitter.Fire(f)
		})
	}
}

// NewTimer :
func NewTimer() *Timer {
	timer := &Timer{}
	timer.Init()
	return timer
}
