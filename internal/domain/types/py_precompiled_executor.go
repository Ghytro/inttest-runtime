package types

import (
	"inttest-runtime/pkg/embedded"
	"sync"
)

type PyPrecompiledExecutor struct {
	execLock sync.Mutex
	pyCtx    *embedded.PyRuntime
	funcs    map[string]*embedded.PyCallable
}

func NewPyPrecompiledExecutor(interpreter *embedded.PyRuntime) *PyPrecompiledExecutor {
	return &PyPrecompiledExecutor{
		execLock: sync.Mutex{},
		pyCtx:    interpreter,
		funcs:    map[string]*embedded.PyCallable{},
	}
}

func (s *PyPrecompiledExecutor) AddFunc(impl embedded.CodeSnippet) (*embedded.PyCallable, error) {
	callable, err := s.pyCtx.NewCallable(impl)
	if err != nil {
		return nil, err
	}
	s.funcs[impl.String()] = callable
	return callable, nil
}

func (s *PyPrecompiledExecutor) ExecFunc(impl embedded.CodeSnippet, args ...embedded.PyValue) (result embedded.PyValue, err error) {
	s.execLock.Lock()
	defer s.execLock.Unlock()

	f, ok := s.funcs[impl.String()]
	if !ok {
		f, err = s.AddFunc(impl)
		if err != nil {
			return embedded.PyValue{}, err
		}
	}

	return f.Call(args...)
}
