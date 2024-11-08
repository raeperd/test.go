package test_test

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/raeperd/test"
)

func Test(t *testing.T) {
	beOkay := func(callback func(tb testing.TB)) {
		t.Helper()
		var buf strings.Builder
		mt := &mockingT{w: &buf}
		callback(mt)
		if mt.Failed() {
			t.Fatal("failed too soon")
		}
		if buf.String() != "" {
			t.Fatal("wrote too much")
		}
	}
	beOkay(func(tb testing.TB) { test.Zero(tb, time.Time{}.Local()) })
	beOkay(func(tb testing.TB) { test.Zero(tb, []string(nil)) })
	beOkay(func(tb testing.TB) { test.NotZero(tb, []string{""}) })
	beOkay(func(tb testing.TB) { test.Nil(tb, nil) })
	beOkay(func(tb testing.TB) { test.NotNil(tb, errors.New("")) })
	beOkay(func(tb testing.TB) { test.True(tb, true) })
	beOkay(func(tb testing.TB) { test.False(tb, false) })
	beOkay(func(tb testing.TB) { test.Contains(tb, "hello world", "world") })
	beOkay(func(tb testing.TB) { test.NotContains(tb, "hello world", "World") })
	beBad := func(callback func(tb testing.TB)) {
		t.Helper()
		var buf strings.Builder
		mt := &mockingT{w: &buf}
		callback(mt)
		if !mt.Failed() {
			t.Fatal("did not fail")
		}
		if buf.String() == "" {
			t.Fatal("wrote too little")
		}
	}
	beBad(func(tb testing.TB) { test.AllEqual(tb, []string{}, []string{""}) })
	beBad(func(tb testing.TB) { test.NotZero(tb, time.Time{}.Local()) })
	beBad(func(tb testing.TB) { test.Zero(tb, []string{""}) })
	beBad(func(tb testing.TB) { test.NotZero(tb, []string(nil)) })
	beBad(func(tb testing.TB) { test.Nil(tb, errors.New("")) })
	beBad(func(tb testing.TB) { test.NotNil(tb, nil) })
	beBad(func(tb testing.TB) { test.True(tb, false) })
	beBad(func(tb testing.TB) { test.False(tb, true) })
	beBad(func(tb testing.TB) { test.Contains(tb, "hello world", "World") })
	beBad(func(tb testing.TB) { test.NotContains(tb, "hello world", "world") })
}

type mockingT struct {
	testing.T
	sync.RWMutex
	failed   bool
	cleanups []func()
	w        io.Writer
}

func (m *mockingT) setFailed(b bool) {
	m.Lock()
	defer m.Unlock()
	m.failed = b
}

func (m *mockingT) Failed() bool {
	m.RLock()
	defer m.RUnlock()
	return m.failed
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
	m.setFailed(true)
	if m.w != nil {
		fmt.Fprintf(m.w, format, args...)
		// Do not call runtime.Goexit here, so that caller can read the output
	} else {
		m.Errorf(format, args...)
		runtime.Goexit()
	}
}

func (m *mockingT) Errorf(format string, args ...any) {
	m.setFailed(true)
	fmt.Printf(format+"\n", args...)
}
