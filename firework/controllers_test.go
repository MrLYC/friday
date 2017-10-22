package firework_test

import (
	"friday/firework"
	"testing"
)

type TestingController struct {
	firework.BaseController
}

func (c *TestingController) GetName() string {
	return "test"
}

func TestTestingControllerInit(t *testing.T) {
	var (
		controller                       = TestingController{}
		icontroller firework.IController = &controller
	)
	if icontroller.GetName() != "test" {
		t.Errorf("name error")
	}
	if controller.Status != firework.StatusControllerInit {
		t.Errorf("init status error")
	}
}

func TestTestingControllerFlow(t *testing.T) {
	var (
		controller                       = TestingController{}
		icontroller firework.IController = &controller
	)
	if controller.Status != firework.StatusControllerInit {
		t.Errorf("init status error")
	}

	icontroller.Ready()
	if controller.Status != firework.StatusControllerReady {
		t.Errorf("ready status error")
	}

	icontroller.Terminate()
	if controller.Status != firework.StatusControllerTerminated {
		t.Errorf("terminate status error")
	}

	icontroller.Kill()
	if controller.Status != firework.StatusControllerTerminated {
		t.Errorf("kill status error")
	}
}
