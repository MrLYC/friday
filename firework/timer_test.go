package firework_test

import (
	"friday/firework"
	"os"
	"testing"
	"time"
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
	value := os.Getenv("FRIDAY_TIMER_CHECKDURATION")
	if value == "" {
		return
	} else {
		timer.CheckDuration, _ = time.ParseDuration(value)
	}
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

	if handlers.(*firework.HandlerItem).Size() != 1 {
		t.Errorf("listen error")
	}

	ch := make(chan string, 1)
	f := &firework.Firework{
		Channel: "test",
		Name:    "delay",
	}
	f.RefreshID()

	delay := 4 * timer.CheckDuration

	emitter.On(f.GetChannel(), f.GetName(), func(ff firework.IFirework) {
		ch <- ff.GetID()
	})
	emitter.On(f.GetChannel(), firework.TimerFireworkAbort, func(ff firework.IFirework) {
		timer.Add(&firework.DelayFirework{
			Time:      time.Now().Add(delay),
			IFirework: f,
		})
	})

	go emitter.Run()

	t1 := time.Now()
	timer.Add(&firework.DelayFirework{
		Time:      t1.Add(delay),
		IFirework: f,
	})

	id := <-ch
	t2 := time.Now()
	if id != f.GetID() {
		t.Errorf("delay error")
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

func TestDurationFirework(t *testing.T) {
	timer := firework.NewTimer()
	value := os.Getenv("FRIDAY_TIMER_CHECKDURATION")
	if value == "" {
		return
	} else {
		timer.CheckDuration, _ = time.ParseDuration(value)
	}
	emitter := &TestingEmitter{
		WillRun: true,
	}
	emitter.Init()
	emitter.AddApplet(timer)
	emitter.Ready()

	ch := make(chan string, 1)
	f := &firework.Firework{
		Channel: "test",
		Name:    "delay",
	}
	f.RefreshID()

	delay := 4 * timer.CheckDuration

	emitter.On(f.GetChannel(), f.GetName(), func(ff firework.IFirework) {
		ch <- ff.GetID()
	})

	go emitter.Run()

	t1 := time.Now()
	df := &firework.DurationFirework{
		Duration: delay,
		Times:    2,
		DelayFirework: &firework.DelayFirework{
			Time:      t1.Add(delay),
			IFirework: f,
		},
	}
	timer.Add(df)

	id1 := <-ch
	t2 := time.Now()
	id2 := <-ch
	t3 := time.Now()

	if df.Times != 1 || df.Status != firework.StatusDelayFireworkSent {
		t.Errorf("firework error")
	}

	if id1 != f.GetID() || id2 != f.GetID() {
		t.Errorf("delay error")
	}

	delta1 := t2.Sub(t1)
	if delta1 < (delay-timer.CheckDuration) || delta1 > (delay+timer.CheckDuration) {
		t.Errorf("timer error")
	}
	delta2 := t3.Sub(t2)
	if delta2 < (delay-timer.CheckDuration) || delta2 > (delay+timer.CheckDuration) {
		t.Errorf("timer error")
	}

	emitter.Terminate()
	if timer.Status != firework.StatusControllerTerminated {
		t.Errorf("terminate error")
	}
}
