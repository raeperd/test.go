package test_test

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
)

type mockingT struct {
	testing.T
	sync.RWMutex
	hasFailed bool
	cleanups  []func()
}

func (m *mockingT) setFailed(b bool) {
	m.Lock()
	defer m.Unlock()
	m.hasFailed = b
}

func (m *mockingT) failed() bool {
	m.RLock()
	defer m.RUnlock()
	return m.hasFailed
}

func (m *mockingT) Run(name string, f func(t *testing.T)) {
	m.cleanups = nil
	m.setFailed(false)
	ch := make(chan struct{})
	defer func() {
		for _, f := range m.cleanups {
			defer f()
		}
	}()
	// Use a goroutine so Fatalf can call Goexit
	go func() {
		defer close(ch)
		f(&m.T)
	}()
	<-ch
}

func (m *mockingT) Cleanup(f func()) {
	m.cleanups = append(m.cleanups, f)
}

func (*mockingT) Log(args ...any) {
	fmt.Println(args...)
}

func (*mockingT) Helper() {}

func (m *mockingT) Fatalf(format string, args ...any) {
	m.Errorf(format, args...)
	runtime.Goexit()
}

func (m *mockingT) Errorf(format string, args ...any) {
	m.setFailed(true)
	fmt.Printf(format+"\n", args...)
}

func (m *mockingT) Failed() bool { return m.failed() }
