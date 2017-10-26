package firework_test

import (
	"friday/firework"
	"testing"
	"time"

	"github.com/emirpasic/gods/sets/treeset"
)

func TestTimerInit(t *testing.T) {
	timer := firework.NewTimer()
	delayFirework := &firework.DelayFirework{}
	timer.Add(delayFirework)

	if timer.Peek() != delayFirework {
		t.Errorf("peek error")
	}
	if timer.Pop() != delayFirework {
		t.Errorf("pop error")
	}

	if timer.Peek() != nil {
		t.Errorf("peek error")
	}
	if timer.Pop() != nil {
		t.Errorf("pop error")
	}
}

func TestTimerFlow(t *testing.T) {
	var ok bool

	timer := firework.NewTimer()
	timer.CheckDuration = 100 * time.Microsecond
	emitter := &TestingEmitter{
		WillRun: true,
	}
	emitter.Init()
	emitter.AddApplet(timer)
	emitter.Ready()

	item, ok := emitter.Channels.Get(firework.TimerChannelName)
	if !ok {
		t.Errorf("channel error")
	}

	chItem := item.(*firework.ChannelItem)
	handlers, ok := chItem.Handlers.Get(firework.TimerFireworkDelay)
	if !ok {
		t.Errorf("handler error")
	}

	if handlers.(*treeset.Set).Size() != 1 {
		t.Errorf("listen error")
	}

	ch := make(chan string, 1)
	f := &firework.Firework{
		Channel: "test",
		Name:    "delay",
	}
	f.RefreshID()

	delay := 4 * timer.CheckDuration

	emitter.On(f.Channel, f.Name, func(ff *firework.Firework) {
		ch <- ff.ID
	})
	emitter.On(f.Channel, firework.TimerFireworkAbort, func(ff *firework.Firework) {
		timer.Add(&firework.DelayFirework{
			Time:     time.Now().Add(delay),
			Firework: f,
		})
	})

	go emitter.Run()

	t1 := time.Now()
	timer.Add(&firework.DelayFirework{
		Time:     t1.Add(delay),
		Firework: f,
	})

	id := <-ch
	t2 := time.Now()
	if id != f.ID {
		t.Errorf("deay error")
	}
	delta := t2.Sub(t1)
	if delta < (delay-timer.CheckDuration) || delta > (delay+timer.CheckDuration) {
		t.Errorf("timer error")
	}
	emitter.Terminate()
	if timer.Status != firework.StatusControllerTerminated {
		t.Errorf("terminate error")
	}
}
