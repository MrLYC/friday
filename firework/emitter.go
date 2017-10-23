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
	Name     string
	Channel  chan *Firework
	Handlers *treemap.Map
	Lock     sync.Mutex
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

	e.appletLock.Lock()
	defer e.appletLock.Unlock()

	_, ok := e.Applets.Get(name)
	if ok {
		return false
	}
	e.Applets.Put(name, applet)
	return true
}

// DeleteApplet :
func (e *Emitter) DeleteApplet(applet IApplet) bool {
	e.appletLock.Lock()
	defer e.appletLock.Unlock()
	e.Applets.Remove(applet)
	return true
}

// DeclareChannel :
func (e *Emitter) DeclareChannel(name string) chan *Firework {
	chanItem, _ := e.declareChannelItem(name)
	return chanItem.Channel
}

func (e *Emitter) declareChannelItem(name string) (*ChannelItem, bool) {
	e.chanLock.Lock()
	defer e.chanLock.Unlock()
	item, ok := e.Channels.Get(name)
	if ok {
		return item.(*ChannelItem), false
	}

	if e.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}

	chanItem := &ChannelItem{
		Name:     name,
		Channel:  make(chan *Firework, config.Configuration.Sentry.ChannelBuffer),
		Handlers: treemap.NewWithStringComparator(),
	}
	e.Channels.Put(name, chanItem)
	return chanItem, true
}

// On :
func (e *Emitter) On(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.Lock.Lock()
	defer channel.Lock.Unlock()

	items, ok := channel.Handlers.Get(name)

	var handlers *treeset.Set
	if !ok {
		handlers = treeset.NewWith(func(a interface{}, b interface{}) int {
			va := reflect.ValueOf(a)
			vb := reflect.ValueOf(b)
			return int(va.Pointer() - vb.Pointer())
		})
		channel.Handlers.Put(name, handlers)
	} else {
		handlers = items.(*treeset.Set)
	}
	handlers.Add(handler)
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
func (e *Emitter) Fire(channelName string, firework *Firework) {
	channel, _ := e.declareChannelItem(channelName)
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
