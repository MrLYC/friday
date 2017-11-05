package vm_test

import (
	"friday/vm"
	"testing"
)

func TestExecute(t *testing.T) {
	vm := vm.VM{}
	vm.Init()
	vm.Execute(`print("hello world")`)
}
