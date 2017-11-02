package vm

import (
	"github.com/Shopify/go-lua"
)

type VM struct {
	LuaState *lua.State
}

func (v *VM) Init() {
	v.LuaState = lua.NewState()
	lua.OpenLibraries(v.LuaState)
}

func (v *VM) Execute(statement string) error {
	return lua.DoString(v.LuaState, statement)
}
