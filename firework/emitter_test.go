package firework_test

import (
	"friday/firework"
	"testing"
)

type TestingEmitter struct {
	firework.Emitter
}

func (e *TestingEmitter) Init() {
	e.SetName("testing")
}

func TestEmitterInit(t *testing.T) {
	var (
		emitter                    = &TestingEmitter{}
		iemiiter firework.IEmitter = emitter
	)
	emitter.Init()
	if iemiiter.GetName() != emitter.Name {
		t.Errorf("name error")
	}
}
