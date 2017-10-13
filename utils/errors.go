package utils

import (
	"bytes"
	"container/list"
	"fmt"
	"runtime"
	"strings"
)

// RunTimeStackFrame :
type RunTimeStackFrame struct {
	Index   int
	File    string
	Line    int
	Address uintptr
	Name    string
}

// Init : load frame from runtime
func (f *RunTimeStackFrame) Init(index int, file string, line int, addr uintptr) {
	f.Address = addr
	f.File = file
	f.Line = line

	var (
		name string
		idx  int
	)

	fn := runtime.FuncForPC(addr)
	if fn != nil {
		name = fn.Name()
	} else {
		name = "<UNKNOWN>"
	}
	idx = strings.LastIndex(name, "/")
	if index != -1 {
		name = name[idx+1:]
	}
	f.Name = name
}

// RunTimeStack :
type RunTimeStack struct {
	Frames *list.List
}

// Init :
func (s *RunTimeStack) Init(skip int) {
	var (
		frame *RunTimeStackFrame
	)
	s.Frames = list.New()
	for i := skip; ; i++ {
		addr, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		frame = &RunTimeStackFrame{}
		frame.Init(i, file, line, addr)
		s.Frames.PushFront(frame)
	}
}

func (s *RunTimeStack) String() string {
	stackInfo := new(bytes.Buffer)
	for i := s.Frames.Front(); i != s.Frames.Back(); i = i.Next() {
		frame := i.Value.(*RunTimeStackFrame)
		fmt.Fprintf(stackInfo, "%v:%v %v\n", frame.File, frame.Line, frame.Name)
	}
	return stackInfo.String()
}

// TraceableError : error interface with stack info
type TraceableError struct {
	message string
	stack   string
}

// Error : return message and stack
func (e *TraceableError) Error() string {
	return fmt.Sprintf("%s\n%s", e.message, e.stack)
}

// Message : return message
func (e *TraceableError) Message() string {
	return e.message
}

// Stack : stack info
func (e *TraceableError) Stack() string {
	return e.stack
}

// Clone : clone and refresh stack info
func (e *TraceableError) Clone() *TraceableError {
	return ErrorNew(e.message, 2)
}

// Panic :
func (e *TraceableError) Panic() {
	panic(*e)
}

// ErrorNew : create a Error
func ErrorNew(message string, skip int) *TraceableError {
	stack := &RunTimeStack{}
	stack.Init(skip)
	return &TraceableError{
		message: message,
		stack:   stack.String(),
	}
}

// Errorf :
func Errorf(message string, v ...interface{}) *TraceableError {
	if v != nil {
		message = fmt.Sprintf(message, v)
	}
	return ErrorNew(message, 2)
}

// ErrorWrap : wrap error to TraceableError
func ErrorWrap(e interface{}) *TraceableError {
	var message string
	switch e := e.(type) {
	case TraceableError:
		return &e
	case *TraceableError:
		return e
	case error:
		message = e.Error()
	default:
		message = fmt.Sprintf("%v", e)
	}
	return ErrorNew(message, 2)
}

// ErrorRecoverCall : recover error
func ErrorRecoverCall(fn func(*TraceableError)) {
	err := recover()
	if err != nil {
		fn(ErrorWrap(err))
	}
}
