package command_test

import (
	"friday/command"
	"testing"
)

func TestInit(t *testing.T) {
	factory := &command.Factory{}
	factory.Init()
	if factory.HelpFlag != "h" {
		t.Errorf("HelpFlag: %v", factory.HelpFlag)
	}
	if factory.Name != "usage" {
		t.Errorf("Name: %v", factory.Name)
	}
}
