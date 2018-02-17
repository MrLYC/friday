package firework_test

import (
	"friday/firework"
	"testing"
)

type TestingEmitter struct {
	firework.Emitter
	WillRun bool
}

func (e *TestingEmitter) Init() {
	e.SetName("testing")
	e.StrictMode = true
	e.Emitter.Init()
}

func (e *TestingEmitter) Run() {
	if e.WillRun {
		e.Emitter.Run()
	} else {
		e.BaseController.Run()
	}
}

func TestEmitterFlow(t *testing.T) {
	var (
		emitter = &TestingEmitter{
			WillRun: false,
		}
		ok bool
	)
	emitter.Init()
	if emitter.GetName() != emitter.Name {
		t.Errorf("name error")
	}
	if emitter.GetStatus() != firework.StatusControllerInit {
		t.Errorf("status error")
	}

	emitter.Ready()
	if emitter.GetStatus() != firework.StatusControllerReady {
		t.Errorf("status error")
	}
	if emitter.Channels.Size() == 0 {
		t.Errorf("not default channels found")
	}

	_, ok = emitter.Channels.Get(firework.ChanBroadcast)
	if !ok {
		t.Errorf("default channel[%s] not found", firework.ChanBroadcast)
	}

	_, ok = emitter.Channels.Get(firework.ChanInternal)
	if !ok {
		t.Errorf("default channel[%s] not found", firework.ChanInternal)
	}

	emitter.Run()
	if emitter.GetStatus() != firework.StatusControllerRuning {
		t.Errorf("status error")
	}

	emitter.Terminate()
	if emitter.GetStatus() != firework.StatusControllerTerminated {
		t.Errorf("status error")
	}

	emitter.Kill()
	if emitter.GetStatus() != firework.StatusControllerTerminated {
		t.Errorf("status error")
	}
}

func TestEmitterDeclareChannelNotReady(t *testing.T) {
	var (
		emitter = &TestingEmitter{}
		name    = "test"
	)
	emitter.Init()
	defer func() {
		err := recover()
		if err != firework.ErrEmitterNotReady {
			t.Errorf("ready error: %v", err)
		}
	}()
	emitter.DeclareChannel(name)
}

func TestEmitterDeclareChannel(t *testing.T) {
	var (
		emitter = &TestingEmitter{}
		name    = "test"
	)
	emitter.Init()
	emitter.Ready()
	ch := emitter.DeclareChannel(name)
	i, _ := emitter.Channels.Get(name)
	chItem := i.(*firework.ChannelItem)
	if chItem.Channel != ch || chItem.Name != name {
		t.Errorf("channel error")
	}
}

func TestEmitterOnNotReady(t *testing.T) {
	var (
		emitter = &TestingEmitter{}
		name    = "test"
	)
	emitter.Init()
	defer func() {
		err := recover()
		if err != firework.ErrEmitterNotReady {
			t.Errorf("ready error: %v", err)
		}
	}()
	emitter.On(name, "mrlyc", func(f firework.IFirework) {
		t.Errorf("ready error")
	})
}

func TestEmitterOn(t *testing.T) {
	var (
		emitter  = &TestingEmitter{}
		name     = "test"
		handler1 firework.Handler
		ok       bool
	)
	emitter.Init()
	emitter.Ready()
	handler1, ok = emitter.On(name, "mrlyc", func(f firework.IFirework) {})
	if !ok {
		t.Errorf("on error")
	}
	_, ok = emitter.Off(name, "mrlyc", handler1)
	if !ok {
		t.Errorf("on error")
	}
}

func TestEmitterFire(t *testing.T) {
	var (
		emitter = &TestingEmitter{}
		name    = "test"
		f1      = &firework.Firework{
			Channel: name,
		}
	)
	emitter.Init()
	emitter.Ready()

	ch := emitter.DeclareChannel(name)
	emitter.Fire(f1)

	f2 := <-ch
	if f1 != f2 {
		t.Errorf("fire error")
	}
}

func TestEmitterRun(t *testing.T) {
	var (
		name    = "test"
		emitter = &TestingEmitter{
			WillRun: true,
		}
	)
	var (
		name1 = "a"
		ev1   = &firework.Firework{
			Channel: name1,
			Name:    name,
		}
		ch1      = make(chan string, 10)
		handler1 = func(f firework.IFirework) {
			ch1 <- f.GetID()
		}
	)
	var (
		name2 = "b"
		ev2   = &firework.Firework{
			Channel: name2,
			Name:    name,
		}
		ch2      = make(chan string, 10)
		handler2 = func(f firework.IFirework) {
			ch2 <- f.GetID()
		}
	)
	var (
		ch22      = make(chan string, 10)
		handler22 = func(f firework.IFirework) {
			ch22 <- f.GetChannel()
		}
	)

	emitter.Init()
	emitter.Ready()
	emitter.DeclareChannel(name1)
	emitter.DeclareChannel(name2)

	emitter.On(name1, "test", handler1)
	emitter.On(name2, "test", handler2)
	emitter.On(name2, "test", handler22)

	ev1.RefreshID()
	emitter.Fire(ev1)

	ev2.RefreshID()
	emitter.Fire(ev2)

	go emitter.Run()

	result1 := <-ch1
	if result1 != ev1.GetID() {
		t.Errorf("handler1 error")
	}

	result2 := <-ch2
	if result2 != ev2.GetID() {
		t.Errorf("handler2 error")
	}

	result22 := <-ch22
	if result22 != ev2.GetChannel() {
		t.Errorf("handler22 error")
	}

	emitter.Terminate()
}
