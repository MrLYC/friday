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
		emitter = &TestingEmitter{}
	)
	emitter.Init()
	if emitter.GetName() != emitter.Name {
		t.Errorf("name error")
	}
}
