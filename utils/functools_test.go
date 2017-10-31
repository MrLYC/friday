package utils_test

import (
	"friday/utils"
	"testing"
)

func TestCartesianProduct(t *testing.T) {
	resultMap := make(map[string][]int)
	utils.CartesianProduct([]interface{}{"a", "b", "c"}, []interface{}{1, 2}, func(v1, v2 interface{}) {
		resultMap[v1.(string)] = append(resultMap[v1.(string)], v2.(int))
	})
	if len(resultMap) != 3 {
		t.Errorf("CartesianProduct error: %v", resultMap)
	}
	if len(resultMap["a"]) != 2 {
		t.Errorf("CartesianProduct error: %v", resultMap)
	}
	if len(resultMap["b"]) != 2 {
		t.Errorf("CartesianProduct error: %v", resultMap)
	}
	if len(resultMap["c"]) != 2 {
		t.Errorf("CartesianProduct error: %v", resultMap)
	}
}

func TestCartesianProductEmpty(t *testing.T) {
	resultMap := make(map[string][]int)
	utils.CartesianProduct([]interface{}{}, []interface{}{}, func(v1, v2 interface{}) {
		resultMap[v1.(string)] = append(resultMap[v1.(string)], v2.(int))
	})
	if len(resultMap) > 0 {
		t.Errorf("CartesianProduct error: %v", resultMap)
	}
}
