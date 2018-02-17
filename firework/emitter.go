package firework

import (
	"friday/config"
	"friday/logging"
	"reflect"
	"sync"
	"time"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
)

// ChannelItem :
type ChannelItem struct {
	Name         string
	Channel      chan IFirework
	Lock         sync.Mutex
	Handlers     *treemap.Map
	handlersLock sync.RWMutex
}

// Emitter :
type Emitter struct {
	BaseController

	StrictMode bool
	RunAt      time.Time

	Channels *treemap.Map
	chanLock sync.RWMutex

	Applets    *treemap.Map
	appletLock sync.RWMutex
}

// Init :
func (e *Emitter) Init() {
	e.Channels = treemap.NewWithStringComparator()
	e.Applets = treemap.NewWithStringComparator()
	e.BaseController.Init()
}

// AddApplet :
func (e *Emitter) AddApplet(applet IApplet) bool {
	name := applet.GetName()

	e.appletLock.RLock()
	_, ok := e.Applets.Get(name)
	e.appletLock.RUnlock()
	if ok {
		return false
	}
	e.appletLock.Lock()
	e.Applets.Put(name, applet)
	e.appletLock.Unlock()
	applet.SetEmitter(e)
	return true
}

// DeleteApplet :
func (e *Emitter) DeleteApplet(applet IApplet) bool {
	e.appletLock.Lock()
	e.Applets.Remove(applet)
	e.appletLock.Unlock()
	return true
}

// DeclareChannel :
func (e *Emitter) DeclareChannel(name string) chan IFirework {
	chanItem, _ := e.declareChannelItem(name)
	return chanItem.Channel
}

func (e *Emitter) declareChannelItem(name string) (*ChannelItem, bool) {
	item, ok := e.Channels.Get(name)
	if ok {
		return item.(*ChannelItem), false
	}
	if e.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}

	chanItem := &ChannelItem{
		Name:     name,
		Channel:  make(chan IFirework, config.Configuration.Firework.ChannelBuffer),
		Handlers: treemap.NewWithStringComparator(),
	}
	e.chanLock.Lock()
	e.Channels.Put(name, chanItem)
	e.chanLock.Unlock()
	return chanItem, true
}

// On :
func (e *Emitter) On(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.handlersLock.RLock()
	items, ok := channel.Handlers.Get(name)
	channel.handlersLock.RUnlock()

	var handlers *treeset.Set
	if !ok {
		handlers = treeset.NewWith(func(a interface{}, b interface{}) int {
			va := reflect.ValueOf(a)
			vb := reflect.ValueOf(b)
			return int(va.Pointer() - vb.Pointer())
		})
		channel.Lock.Lock()
		channel.Handlers.Put(name, handlers)
		channel.Lock.Unlock()
	} else {
		handlers = items.(*treeset.Set)
	}
	channel.handlersLock.Lock()
	handlers.Add(handler)
	channel.handlersLock.Unlock()
	return handler, true
}

// Off :
func (e *Emitter) Off(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.Lock.Lock()
	defer channel.Lock.Unlock()

	channel.handlersLock.RLock()
	items, ok := channel.Handlers.Get(name)
	channel.handlersLock.RUnlock()

	if !ok {
		return handler, false
	}
	handlers := items.(*treeset.Set)

	handlers.Remove(handler)

	return handler, true
}

// Fire :
func (e *Emitter) Fire(firework IFirework) {
	if firework.GetID() == "" {
		firework.RefreshID()
	}
	channel, _ := e.declareChannelItem(firework.GetChannel())
	channel.Channel <- firework
}

// FireAt :
func (e *Emitter) FireAt(at time.Time, firework IFirework) {
	e.Fire(NewDelayFirework(at, &Firework{
		Channel: TimerChannelName,
		Sender:  firework.GetSender(),
		Name:    TimerFireworkDelay,
		Payload: firework,
	}))
}

// FireDelay :
func (e *Emitter) FireDelay(duration time.Duration, firework IFirework) {
	e.FireDelayN(duration, 1, firework)
}

// FireDelayN :
func (e *Emitter) FireDelayN(duration time.Duration, times uint, firework IFirework) {
	e.Fire(NewDurationFirework(duration, times, &Firework{
		Channel: TimerChannelName,
		Sender:  firework.GetSender(),
		Name:    TimerFireworkDelay,
		Payload: firework,
	}))
}

// FireCron :
func (e *Emitter) FireCron(rule string, firework IFirework) {
	e.Fire(NewCronFirework(rule, time.Now(), &Firework{
		Channel: TimerChannelName,
		Sender:  firework.GetSender(),
		Name:    TimerFireworkDelay,
		Payload: firework,
	}))
}

// Run :
func (e *Emitter) Run() {
	if e.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}
	e.BaseController.Run()

	e.RunAt = time.Now()
	logging.Infof("Emitter run at: %s", e.RunAt.String())

	e.appletLock.RLock()
	iterApplet := e.Applets.Iterator()
	for iterApplet.Next() {
		applet := iterApplet.Value().(IApplet)
		go applet.Run()
	}
	e.appletLock.RUnlock()

	e.chanLock.RLock()
	channels := make([]reflect.SelectCase, e.Channels.Size())
	iterChannel := e.Channels.Iterator()
	for i := 0; iterChannel.Next(); i++ {
		item := iterChannel.Value().(*ChannelItem)
		channels[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(item.Channel),
		}
	}
	e.chanLock.RUnlock()

	for e.Status == StatusControllerRuning {
		chosen, value, ok := reflect.Select(channels)
		if !ok {
			logging.Warningf("Channel[%v] error", chosen)
			if e.StrictMode {
				break
			}
			continue
		}
		firework := value.Interface().(IFirework)
		channel := firework.GetChannel()
		e.chanLock.RLock()
		chanItem, ok := e.Channels.Get(channel)
		e.chanLock.RUnlock()
		if !ok {
			logging.Warningf(
				"Unknown channel %s from %s(%s)",
				channel, firework.GetSender(), firework.GetName(),
			)
			if e.StrictMode {
				break
			}
			continue
		}

		items, ok := chanItem.(*ChannelItem).Handlers.Get(firework.GetName())
		if !ok {
			continue
		}
		handlers := items.(*treeset.Set)
		handlers.Each(func(index int, value interface{}) {
			f := firework.Copy()
			go value.(Handler)(f)
		})
	}
}

// Ready :
func (e *Emitter) Ready() {
	e.BaseController.Ready()
	e.appletLock.RLock()
	iter := e.Applets.Iterator()
	for iter.Next() {
		applet := iter.Value().(IApplet)
		applet.Ready()
	}
	e.appletLock.RUnlock()
}

// Terminate :
func (e *Emitter) Terminate() {
	e.SetStatus(StatusControllerTerminating)
	e.appletLock.Lock()
	e.Applets.Each(func(key interface{}, value interface{}) {
		applet := value.(IApplet)
		status := applet.GetStatus()
		if status != StatusControllerTerminated && status != StatusControllerTerminating {
			applet.Terminate()
		}
	})
	e.appletLock.Unlock()
	e.chanLock.Lock()
	e.Channels.Each(func(key interface{}, value interface{}) {
		channel := value.(*ChannelItem)
		close(channel.Channel)
	})
	e.chanLock.Unlock()
	e.BaseController.Terminate()
}
