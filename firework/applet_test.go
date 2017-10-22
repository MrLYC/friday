package firework_test

import (
	"friday/firework"
	"testing"
)

type TestingApplet struct {
	firework.Applet
}

func (a *TestingApplet) Init(name string) {
	a.SetName(name)
}

func TestAppletInit(t *testing.T) {
	var (
		applet                    = &TestingApplet{}
		trigger firework.ITrigger = applet
	)

	applet.Init("test")
	if trigger.GetName() != applet.Name {
		t.Errorf("name error")
	}
}
