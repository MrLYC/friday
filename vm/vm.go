package vm

import (
	"github.com/yuin/gopher-lua"
)

type VM struct {
	LuaState *lua.LState
}

func NewVM() *VM {
	vm := &VM{}
	vm.Init()
	return vm
}

func (v *VM) Init() {
	v.LuaState = lua.NewState()
}

func (v *VM) Execute(statement string) error {
	return v.LuaState.DoString(statement)
}

func (v *VM) GetGlobal(name string) interface{} {
	value := v.LuaState.GetGlobal(name)
	switch value.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(value.(lua.LBool))
	case lua.LTNumber:
		return float64(value.(lua.LNumber))
	case lua.LTString:
		return string(value.(lua.LString))
	default:
		return value
	}
}
