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
	Channel      chan *Firework
	Lock         sync.Mutex
	Handlers     *treemap.Map
	HandlersLock sync.Mutex
}

// Emitter :
type Emitter struct {
	BaseController

	StrictMode bool
	RunAt      time.Time

	Channels *treemap.Map
	chanLock sync.Mutex

	Applets    *treemap.Map
	appletLock sync.Mutex
}

// Init :
func (e *Emitter) Init() {
	e.Channels = treemap.NewWithStringComparator()
	e.Applets = treemap.NewWithStringComparator()
	e.SetStatus(StatusControllerInit)
}

// AddApplet :
func (e *Emitter) AddApplet(applet IApplet) bool {
	name := applet.GetName()

	_, ok := e.Applets.Get(name)
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
func (e *Emitter) DeclareChannel(name string) chan *Firework {
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
		Channel:  make(chan *Firework, config.Configuration.Firework.ChannelBuffer),
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

	items, ok := channel.Handlers.Get(name)

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
	channel.HandlersLock.Lock()
	handlers.Add(handler)
	channel.HandlersLock.Unlock()
	return handler, true
}

// Off :
func (e *Emitter) Off(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.Lock.Lock()
	defer channel.Lock.Unlock()

	items, ok := channel.Handlers.Get(name)

	if !ok {
		return handler, false
	}
	handlers := items.(*treeset.Set)

	handlers.Remove(handler)

	return handler, true
}

// Fire :
func (e *Emitter) Fire(firework *Firework) {
	channel, _ := e.declareChannelItem(firework.Channel)
	channel.Channel <- firework
}

// Run :
func (e *Emitter) Run() {
	if e.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}
	e.BaseController.Run()

	e.RunAt = time.Now()
	logging.Infof("Emitter run at: %s", e.RunAt.String())

	channels := make([]reflect.SelectCase, e.Channels.Size())
	iter := e.Channels.Iterator()
	for i := 0; iter.Next(); i++ {
		item := iter.Value().(*ChannelItem)
		channels[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(item.Channel),
		}
	}

	for e.Status == StatusControllerRuning {
		chosen, value, ok := reflect.Select(channels)
		if !ok {
			logging.Warningf("Channel[%v] error", chosen)
			if e.StrictMode {
				break
			}
			continue
		}
		firework := value.Interface().(*Firework)
		chanItem, ok := e.Channels.Get(firework.Channel)
		if !ok {
			logging.Warningf(
				"Unknown channel %s from %s(%s)",
				firework.Channel, firework.Sender, firework.Name,
			)
			if e.StrictMode {
				break
			}
			continue
		}
		items, ok := chanItem.(*ChannelItem).Handlers.Get(firework.Name)
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

// Terminate :
func (e *Emitter) Terminate() {
	e.SetStatus(StatusControllerTerminating)
	e.Applets.Each(func(key interface{}, value interface{}) {
		applet := value.(IApplet)
		status := applet.GetStatus()
		if status != StatusControllerTerminated && status != StatusControllerTerminating {
			applet.Terminate()
		}
	})
	e.Channels.Each(func(key interface{}, value interface{}) {
		channel := value.(*ChannelItem)
		close(channel.Channel)
	})
	e.BaseController.Terminate()
}
