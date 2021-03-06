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
	TimerFireworkDelay   = "Delay"
	TimerFireworkTimesUp = "TimesUp"
	TimerFireworkAbort   = "TimerAbort"
)

// IDelayFirework :
type IDelayFirework interface {
	GetTime() time.Time
	UpdateTime() bool
	GetFirework() IFirework
	GetStatus() DelayFireworkStatus
	SetStatus(DelayFireworkStatus)
}

// DelayFirework :
type DelayFirework struct {
	IFirework
	Time   time.Time
	Status DelayFireworkStatus
}

// Copy :
func (f *DelayFirework) Copy() IFirework {
	return &DelayFirework{
		Time:      f.Time,
		Status:    f.Status,
		IFirework: f.IFirework.Copy(),
	}
}

// GetTime :
func (f *DelayFirework) GetTime() time.Time {
	return f.Time
}

// UpdateTime :
func (f *DelayFirework) UpdateTime() bool {
	return false
}

// GetFirework :
func (f *DelayFirework) GetFirework() IFirework {
	return f.IFirework
}

// GetStatus :
func (f *DelayFirework) GetStatus() DelayFireworkStatus {
	return f.Status
}

// SetStatus :
func (f *DelayFirework) SetStatus(status DelayFireworkStatus) {
	f.Status = status
}

// NewDelayFirework :
func NewDelayFirework(at time.Time, firework IFirework) *DelayFirework {
	f := &DelayFirework{}
	f.Time = at
	f.IFirework = firework
	return f
}

// DurationFirework :
type DurationFirework struct {
	*DelayFirework
	Duration time.Duration
	Times    uint
}

// Copy :
func (f *DurationFirework) Copy() IFirework {
	return &DurationFirework{
		Duration:      f.Duration,
		Times:         f.Times,
		DelayFirework: f.DelayFirework.Copy().(*DelayFirework),
	}
}

// UpdateTime :
func (f *DurationFirework) UpdateTime() bool {
	if f.Times == 1 {
		f.SetStatus(StatusDelayFireworkSent)
		return false
	} else if f.Times > 1 {
		f.Times--
	}
	// 0 for forever
	f.Time = f.Time.Add(f.Duration)
	return true
}

// NewDurationFirework :
func NewDurationFirework(duration time.Duration, times uint, firework IFirework) *DurationFirework {
	f := &DurationFirework{}
	f.DelayFirework = &DelayFirework{}
	f.Duration = duration
	f.Times = times
	f.IFirework = firework
	return f
}

func delayFireworkComparator(i1 interface{}, i2 interface{}) int {
	t1 := i1.(IDelayFirework).GetTime()
	t2 := i2.(IDelayFirework).GetTime()
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
	lock          sync.RWMutex
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
func (t *Timer) Add(firework IDelayFirework) {
	t.lock.Lock()
	t.heap.Push(firework)
	t.lock.Unlock()
	firework.SetStatus(StatusDelayFireworkReady)
}

// Pop :
func (t *Timer) Pop() IDelayFirework {
	t.lock.Lock()
	f, ok := t.heap.Pop()
	t.lock.Unlock()
	if !ok {
		return nil
	}
	return f.(IDelayFirework)
}

// Peek :
func (t *Timer) Peek() IDelayFirework {
	t.lock.RLock()
	f, ok := t.heap.Peek()
	t.lock.RUnlock()
	if !ok {
		return nil
	}
	return f.(IDelayFirework)
}

// Ready :
func (t *Timer) Ready() {
	t.Emitter.On(ChanNameTimer, TimerFireworkDelay, func(firework IFirework) {
		f := firework.GetPayload().(DelayFirework)
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
			eventTime := f.GetTime()
			if eventTime.After(highLimitTime) {
				break
			}
			f = t.Pop()
			if eventTime.Before(lowLimitTime) || f.GetStatus() != StatusDelayFireworkReady {
				abortFireworks.Add(f)
			} else {
				activateFireworks.Add(f)
			}
		}
		go activateFireworks.Each(func(i int, f interface{}) {
			df := f.(IDelayFirework)
			t.Emitter.Fire(df.GetFirework())
			if df.UpdateTime() {
				t.Add(df)
			}
		})
		go abortFireworks.Each(func(i int, f interface{}) {
			df := f.(IDelayFirework)
			rf := df.GetFirework()
			logging.Warningf(
				"Timer abort at %v-%v: Channel[%s], Name[%s], Time[%v]",
				lowLimitTime, highLimitTime, rf.GetChannel(), rf.GetName(), df.GetTime(),
			)
			rf.SetName(TimerFireworkAbort)
			t.Emitter.Fire(rf)
			if df.UpdateTime() {
				t.Add(df)
			}
		})
	}
}

// NewTimer :
func NewTimer() *Timer {
	timer := &Timer{}
	timer.Init()
	return timer
}
