package vm_test

import (
	"fmt"
	"friday/vm"
	"testing"
)

func TestExecute(t *testing.T) {
	v := vm.VM{}
	v.Init()
	v.Execute(`print("hello world")`)
}

func TestGetGlobal(t *testing.T) {
	v := vm.NewVM()
	items := map[string]interface{}{
		"null":  nil,
		"true":  true,
		"false": false,
		"1":     1.0,
		"2.2":   2.2,
		`"3"`:   "3",
	}
	for statement, value := range items {
		v.Execute(fmt.Sprintf("value = %s", statement))
		if v.GetGlobal("value") != value {
			t.Errorf("get global error")
		}
	}
}
