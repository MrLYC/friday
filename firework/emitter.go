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

// HandlerItem :
type HandlerItem struct {
	*treeset.Set
	lock sync.RWMutex
}

// SafeEach :
func (i *HandlerItem) SafeEach(f func(index int, value interface{})) {
	i.lock.RLock()
	items := make([]interface{}, i.Size())
	iter := i.Iterator()
	for iter.Next() {
		items[iter.Index()] = iter.Value()
	}
	i.lock.RUnlock()

	for index, value := range items {
		f(index, value)
	}
}

// ChannelItem :
type ChannelItem struct {
	Name     string
	Channel  chan IFirework
	lock     sync.RWMutex
	Handlers *treemap.Map
}

// SafeEachHandlers :
func (i *ChannelItem) SafeEachHandlers(name string, f func(index int, value interface{})) {
	i.lock.RLock()
	items, ok := i.Handlers.Get(name)
	i.lock.RUnlock()
	if !ok {
		return
	}

	handlers := items.(*HandlerItem)
	handlers.SafeEach(f)
}

// Emitter :
type Emitter struct {
	BaseController

	StrictMode bool
	RunAt      time.Time

	Channels    *treemap.Map
	channelLock sync.RWMutex

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
	e.channelLock.Lock()
	e.Channels.Put(name, chanItem)
	e.channelLock.Unlock()
	return chanItem, true
}

// On :
func (e *Emitter) On(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.lock.RLock()
	items, ok := channel.Handlers.Get(name)
	channel.lock.RUnlock()

	var handlers *HandlerItem
	if !ok {
		handlers = &HandlerItem{
			Set: treeset.NewWith(func(a interface{}, b interface{}) int {
				va := reflect.ValueOf(a)
				vb := reflect.ValueOf(b)
				return int(va.Pointer() - vb.Pointer())
			}),
		}
		channel.lock.Lock()
		channel.Handlers.Put(name, handlers)
		channel.lock.Unlock()
	} else {
		handlers = items.(*HandlerItem)
	}

	handlers.lock.Lock()
	handlers.Add(handler)
	handlers.lock.Unlock()

	return handler, true
}

// Off :
func (e *Emitter) Off(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.lock.RLock()
	items, ok := channel.Handlers.Get(name)
	channel.lock.RUnlock()

	if !ok {
		return handler, false
	}
	handlers := items.(*HandlerItem)

	handlers.lock.Lock()
	handlers.Remove(handler)
	handlers.lock.Unlock()

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

// runApplets :
func (e *Emitter) runApplets() {
	e.appletLock.RLock()
	iterApplet := e.Applets.Iterator()
	for iterApplet.Next() {
		applet := iterApplet.Value().(IApplet)
		go applet.Run()
	}
	e.appletLock.RUnlock()
}

// getChanSelectCases :
func (e *Emitter) getChanSelectCases() []reflect.SelectCase {
	e.channelLock.RLock()
	channels := make([]reflect.SelectCase, e.Channels.Size())
	iterChannel := e.Channels.Iterator()
	for i := 0; iterChannel.Next(); i++ {
		item := iterChannel.Value().(*ChannelItem)
		channels[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(item.Channel),
		}
	}
	e.channelLock.RUnlock()
	return channels
}

// Run :
func (e *Emitter) Run() {
	if e.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}
	e.BaseController.Run()

	e.RunAt = time.Now()
	logging.Infof("Emitter run at: %s", e.RunAt.String())

	e.runApplets()

	channels := e.getChanSelectCases()

	for e.Status == StatusControllerRuning {
		chosen, selectcase, ok := reflect.Select(channels)
		if !ok {
			logging.Warningf("Channel[%v] error", chosen)
			if e.StrictMode {
				break
			}
			continue
		}

		ifirework := selectcase.Interface()
		if ifirework == nil {
			continue
		}

		firework := ifirework.(IFirework)
		chName := firework.GetChannel()

		e.channelLock.RLock()
		ichanItem, ok := e.Channels.Get(chName)
		e.channelLock.RUnlock()
		if !ok {
			logging.Warningf(
				"Unknown channel %s from %s(%s)",
				chName, firework.GetSender(), firework.GetName(),
			)
			if e.StrictMode {
				break
			}
			continue
		}

		chanItem := ichanItem.(*ChannelItem)
		chanItem.SafeEachHandlers(firework.GetName(), func(index int, value interface{}) {
			f := firework.Copy()
			go value.(Handler)(f)
		})
	}
}

func (e *Emitter) setDefaultChannels() {
	e.DeclareChannel(ChanBroadcast)
	e.DeclareChannel(ChanInternal)
}

// Ready :
func (e *Emitter) Ready() {
	e.BaseController.Ready()

	e.setDefaultChannels()

	e.appletLock.RLock()
	iter := e.Applets.Iterator()
	for iter.Next() {
		applet := iter.Value().(IApplet)
		applet.Ready()
	}
	e.appletLock.RUnlock()
}

// Kill :
func (e *Emitter) Kill() {
	e.SetStatus(StatusControllerTerminating)
	e.appletLock.Lock()
	e.Applets.Each(func(key interface{}, value interface{}) {
		applet := value.(IApplet)
		status := applet.GetStatus()
		if status != StatusControllerTerminated {
			applet.Kill()
		}
	})
	e.appletLock.Unlock()

	e.channelLock.Lock()
	e.Channels.Each(func(key interface{}, value interface{}) {
		channel := value.(*ChannelItem)
		close(channel.Channel)
	})
	e.channelLock.Unlock()
	e.BaseController.Kill()
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

	e.channelLock.Lock()
	e.Channels.Each(func(key interface{}, value interface{}) {
		channel := value.(*ChannelItem)
		channel.Channel <- nil
	})
	e.channelLock.Unlock()
	e.BaseController.Terminate()
}
