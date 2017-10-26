package utils_test

import (
	"friday/utils"
	"testing"
)

type TestingContext struct {
	Value int
}

type TestingHook struct {
	utils.Hook
	Chan chan int
}

func (t *TestingHook) Fire(context utils.IHookContext) utils.HookAction {
	cxt := context.(*TestingContext)
	t.Chan <- cxt.Value
	cxt.Value += 1
	return t.Hook.Fire(context)
}

func (t *TestingContext) Copy() utils.IHookContext {
	return &TestingContext{
		Value: t.Value,
	}
}

func TestContext(t *testing.T) {
	var (
		cxt1 utils.IHookContext = &TestingContext{Value: 1}
		cxt2                    = cxt1.Copy()
	)
	if cxt1.(*TestingContext).Value != cxt2.(*TestingContext).Value {
		t.Errorf("copy error")
	}
}

func TestHook(t *testing.T) {
	var (
		hook  utils.IHook = &TestingHook{Chan: make(chan int, 1)}
		value int
		cxt   = &TestingContext{
			Value: 0,
		}
	)
	go hook.Fire(cxt)
	value = <-hook.(*TestingHook).Chan
	if value != 0 {
		t.Errorf("hook error")
	}
	if cxt.Value != 1 {
		t.Errorf("hook error")
	}

	hook.Wrap(&TestingHook{Chan: make(chan int, 1)})
	go hook.Fire(cxt)
	value = <-hook.(*TestingHook).Chan
	if value != 1 {
		t.Errorf("hook error")
	}
	value = <-hook.(*TestingHook).Next.(*TestingHook).Chan
	if value != 2 {
		t.Errorf("hook error")
	}
	if cxt.Value != 3 {
		t.Errorf("hook error")
	}
}
