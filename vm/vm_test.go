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

func TestTableEachAsArray(t *testing.T) {
	var (
		v   = vm.NewVM()
		arr = make([]float64, 0, 3)
	)

	v.Execute(`value = {1, 2, 3}`)
	value := v.GetGlobal("value").(vm.Table)

	value.EachAsArray(func(index interface{}, value interface{}) {
		arr = append(arr, value.(float64))
	})

	if arr[0] != 1 || arr[1] != 2 || arr[2] != 3 {
		t.Errorf("array error: %v", arr)
	}
}

func TestTableEachAsStrMap(t *testing.T) {
	var (
		v    = vm.NewVM()
		dict = make(map[string]float64)
	)

	v.Execute(`value = {a=1, b=2, c=3}`)
	value := v.GetGlobal("value").(vm.Table)

	value.EachAsStrMap(func(key interface{}, value interface{}) {
		dict[key.(string)] = value.(float64)
	})

	if dict["a"] != 1 || dict["b"] != 2 || dict["c"] != 3 {
		t.Errorf("map error: %v", dict)
	}
}

func TestTableMix(t *testing.T) {
	var (
		v    = vm.NewVM()
		dict = make(map[string]float64)
		arr  = make([]float64, 0, 2)
	)

	v.Execute(`value = {a=1, 2, b=3, 4}`)
	value := v.GetGlobal("value").(vm.Table)

	value.EachAsStrMap(func(key interface{}, value interface{}) {
		dict[key.(string)] = value.(float64)
	})
	if dict["a"] != 1 || dict["b"] != 3 {
		t.Errorf("map error: %v", dict)
	}

	value.EachAsArray(func(index interface{}, value interface{}) {
		arr = append(arr, value.(float64))
	})
	if arr[0] != 2 || arr[1] != 4 {
		t.Errorf("array error: %v", arr)
	}
}
