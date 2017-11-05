package vm

import (
	"github.com/yuin/gopher-lua"
)

type VM struct {
	LuaState *lua.LState
}

func (v *VM) Init() {
	v.LuaState = lua.NewState()
}

func (v *VM) Execute(statement string) error {
	return v.LuaState.DoString(statement)
}
