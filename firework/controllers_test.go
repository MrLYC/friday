package firework_test

import (
	"friday/firework"
	"testing"
)

type TestingController struct {
	firework.BaseController
}

func TestTestingControllerInit(t *testing.T) {
	var (
		controller                       = TestingController{}
		icontroller firework.IController = &controller
	)
	icontroller.SetName("test")
	if icontroller.GetName() != "test" {
		t.Errorf("name error")
	}
	if icontroller.GetStatus() != controller.Status || controller.Status != firework.StatusControllerInit {
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
	if icontroller.GetStatus() != firework.StatusControllerReady {
		t.Errorf("ready status error")
	}

	icontroller.Run()
	if icontroller.GetStatus() != firework.StatusControllerRuning {
		t.Errorf("Run status error")
	}

	icontroller.Terminate()
	if icontroller.GetStatus() != firework.StatusControllerTerminated {
		t.Errorf("terminate status error")
	}

	icontroller.Kill()
	if icontroller.GetStatus() != firework.StatusControllerTerminated {
		t.Errorf("kill status error")
	}
}
