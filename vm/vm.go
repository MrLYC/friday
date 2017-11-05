package vm

import (
	"github.com/yuin/gopher-lua"
)

type Table map[interface{}]interface{}

type TableEachFunc func(key interface{}, value interface{})

func (t Table) EachAsArray(f TableEachFunc) {
	lenTable := len(t)
	for index := 1; index <= lenTable; index++ {
		value, ok := t[float64(index)]
		if !ok {
			break
		}
		f(index-1, value)
	}
}

func (t Table) EachAsStrMap(f TableEachFunc) {
	for key, value := range t {
		switch key.(type) {
		case string:
			f(key.(string), value)
		}
	}
}

type VM struct {
	LuaState *lua.LState
}

func NewVM() *VM {
	vm := &VM{}
	vm.Init()
	return vm
}

// Init :
func (v *VM) Init() {
	v.LuaState = lua.NewState()
}

// Execute :
func (v *VM) Execute(statement string) error {
	return v.LuaState.DoString(statement)
}

// ParseValue :
func (v *VM) ParseValue(value lua.LValue) interface{} {
	switch value.Type() {
	case lua.LTNil:
		return nil
	case lua.LTBool:
		return bool(value.(lua.LBool))
	case lua.LTNumber:
		return float64(value.(lua.LNumber))
	case lua.LTString:
		return string(value.(lua.LString))
	case lua.LTTable:
		table := Table{}
		value.(*lua.LTable).ForEach(func(lkey lua.LValue, litem lua.LValue) {
			key := v.ParseValue(lkey)
			item := v.ParseValue(litem)
			table[key] = item
		})
		return table
	default:
		return value
	}
}

// GetGlobal :
func (v *VM) GetGlobal(name string) interface{} {
	value := v.LuaState.GetGlobal(name)
	return v.ParseValue(value)
}
