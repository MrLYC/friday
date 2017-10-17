package sentry_test

import (
	"friday/sentry"
	"testing"
)

type TestingController struct {
	sentry.BaseController
}

func (c *TestingController) GetName() string {
	return "test"
}

func TestTestingControllerInit(t *testing.T) {
	var (
		controller                     = TestingController{}
		icontroller sentry.IController = &controller
	)
	if icontroller.GetName() != "test" {
		t.Errorf("name error")
	}
	if controller.Status != sentry.StatusControllerInit {
		t.Errorf("init status error")
	}
}

func TestTestingControllerFlow(t *testing.T) {
	var (
		controller                     = TestingController{}
		icontroller sentry.IController = &controller
	)
	if controller.Status != sentry.StatusControllerInit {
		t.Errorf("init status error")
	}

	icontroller.Ready()
	if controller.Status != sentry.StatusControllerReady {
		t.Errorf("ready status error")
	}

	icontroller.Terminate()
	if controller.Status != sentry.StatusControllerTerminated {
		t.Errorf("terminate status error")
	}

	icontroller.Kill()
	if controller.Status != sentry.StatusControllerTerminated {
		t.Errorf("kill status error")
	}
}
