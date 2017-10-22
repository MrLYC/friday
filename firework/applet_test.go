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
		applet                   = &TestingApplet{}
		iapplet firework.IApplet = applet
	)

	applet.Init("test")
	if iapplet.GetName() != applet.Name {
		t.Errorf("name error")
	}
	if iapplet.GetStatus() != firework.StatusControllerInit {
		t.Errorf("status error")
	}
}
