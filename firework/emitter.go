package firework

import (
	"friday/config"
	"sync"

	"github.com/emirpasic/gods/maps/treemap"
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
	if e.Status != StatusControllerReady {
		panic(ErrEmitterNotReady)
	}
	e.chanLock.Lock()
	defer e.chanLock.Unlock()
	item, ok := e.Channels.Get(name)
	if ok {
		return item.(*ChannelItem), false
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
	channel.Handlers.Put(name, handler)
	return handler, true
}

// Off :
func (e *Emitter) Off(channelName string, name string, handler Handler) (Handler, bool) {
	channel, _ := e.declareChannelItem(channelName)

	channel.Lock.Lock()
	defer channel.Lock.Unlock()
	channel.Handlers.Remove(name)
	return handler, true
}

// Fire :
func (e *Emitter) Fire(channelName string, firework *Firework) {
	channel, _ := e.declareChannelItem(channelName)
	channel.Channel <- firework
}
