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
	emitter.On(name, "mrlyc", func(f *firework.Firework) {
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
	handler1, ok = emitter.On(name, "mrlyc", func(f *firework.Firework) {})
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
	emitter.Fire(name, f1)

	f2 := <-ch
	if f1 != f2 {
		t.Errorf("fire error")
	}
}
