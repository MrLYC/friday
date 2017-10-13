package utils_test

import (
	"errors"
	"fmt"
	"friday/utils"
	"strings"
	"testing"
)

func TestErrorNew(t *testing.T) {
	var (
		err     = utils.ErrorNew("mrlyc", 1)
		stack   = err.Stack()
		message = err.Message()
		info    = err.Error()
	)
	if strings.Index(message, "mrlyc") == -1 {
		t.Errorf("Message not found: %v", "mrlyc")
	}
	if strings.Index(info, "mrlyc") == -1 {
		t.Errorf("Info not found: %v", "mrlyc")
	}
	if strings.Index(stack, "TestErrorNew") == -1 {
		t.Errorf("Stack not found: %v", "TestErrorNew")
	}
	if strings.Index(info, "TestErrorNew") == -1 {
		t.Errorf("Info not found: %v", "TestErrorNew")
	}
}

func TestErrorf(t *testing.T) {
	var (
		name1 = "yakov1"
		name2 = "yakov2"
	)
	if strings.Index(utils.Errorf(name1).Message(), name1) == -1 {
		t.Errorf("Message not found: %v", name1)
	}
	if strings.Index(utils.Errorf("This is %v", name2).Message(), name2) == -1 {
		t.Errorf("Message not found: %v", name2)
	}
}

func TestErrorfInGoroutine(t *testing.T) {
	var (
		ch  = make(chan *utils.TraceableError)
		err *utils.TraceableError
	)
	go func() {
		e := utils.Errorf("mrlyc")
		ch <- e
	}()

	err = <-ch
	close(ch)
	if strings.Index(err.Stack(), "TestErrorfInGoroutine.func1") == -1 {
		t.Error(err)
	}
}

func TestErrorWrap(t *testing.T) {
	var (
		str  = "string"
		err  = errors.New("errors")
		terr = utils.Errorf("tracable")
		e    *utils.TraceableError
	)
	e = utils.ErrorWrap(str)
	if e.Message() != str || strings.Index(e.Stack(), "TestErrorWrap") == -1 {
		t.Error(e)
	}
	e = utils.ErrorWrap(err)
	if e.Message() != err.Error() || strings.Index(e.Stack(), "TestErrorWrap") == -1 {
		t.Error(e)
	}
	e = utils.ErrorWrap(terr)
	if *e != *terr {
		t.Error(e)
	}
	e = utils.ErrorWrap(*terr)
	if *e != *terr {
		t.Error(e)
	}
}

func TestTraceableErrorClone(t *testing.T) {
	var (
		err1 = utils.Errorf("mrlyc")
		err2 = err1.Clone()
	)
	if err1.Message() != err2.Message() || err1.Stack() == err2.Stack() {
		t.Errorf("error clone failed")
	}
}

func TestTraceableErrorPanic(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			err := utils.ErrorWrap((e))
			if strings.Index(err.Stack(), "TestTraceableErrorPanic") == -1 {
				t.Error(err)
			}
			if strings.Index(err.Message(), "mrlyc") == -1 {
				t.Error(err)
			}
		}
	}()
	utils.Errorf("mrlyc").Panic()
}

func TestErrorRecoverCall1(t *testing.T) {
	defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
		if strings.Index(err.Stack(), "TestErrorRecoverCall1") == -1 {
			t.Error(err)
		}
		if strings.Index(err.Message(), "mrlyc") == -1 {
			t.Error(err)
		}
	})
	utils.Errorf("mrlyc").Panic()
}

func TestErrorRecoverCall2(t *testing.T) {
	defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
		if strings.Index(err.Stack(), "TestErrorRecoverCall2") == -1 {
			t.Error(err)
		}
		if strings.Index(err.Message(), "mrlyc") == -1 {
			t.Error(err)
		}
	})
	panic("mrlyc")
}

func TestErrorRecoverCall3(t *testing.T) {
	defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
		if strings.Index(err.Stack(), "TestErrorRecoverCall3") == -1 {
			t.Error(err)
		}
		if strings.Index(err.Message(), "mrlyc") == -1 {
			t.Error(err)
		}
	})
	panic(fmt.Errorf("mrlyc"))
}

func TestErrorRecoverCall4(t *testing.T) {
	defer utils.ErrorRecoverCall(func(err *utils.TraceableError) {
		t.Errorf("will not happend")
	})
}
