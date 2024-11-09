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
	test.Equal(t, 1, 1)
	test.NotEqual(t, 1, 2)

	test.Zero(t, time.Time{}.Local())
	test.Zero(t, []string(nil))
	test.NotZero(t, []string{""})

	test.Nil(t, nil)
	test.NotNil(t, errors.New(""))

	test.True(t, true)
	test.False(t, false)

	test.Contains(t, "hello world", "world")
	test.Contains(t, []int{1, 2, 3, 4, 5}, 3)
	test.NotContains(t, "hello world", "World")
	test.NotContains(t, []int{1, 2, 3, 4, 5}, 6)
}

func TestNot(t *testing.T) {
	testFail := func(t *testing.T, want string, testFunc func(*mockingT)) {
		t.Helper()
		var buf strings.Builder
		m := &mockingT{w: &buf}

		testFunc(m)
		if len(want) != 0 { // if want is not empty test message
			test.Equal(t, want, buf.String())
		}
	}

	testFail(t, "want: 1; got: 2", func(m *mockingT) {
		test.Equal(m, 1, 2)
	})
	testFail(t, "got: 1", func(m *mockingT) {
		test.NotEqual(m, 1, 1)
	})
	testFail(t, "len(want): 0; len(got): 1", func(m *mockingT) {
		test.AllEqual(m, []string{}, []string{""})
	})
	testFail(t, "", func(m *mockingT) { // skip message test, local time is not predictable
		test.NotZero(m, time.Time{}.Local())
	})
	testFail(t, "got: []", func(m *mockingT) {
		test.Zero(m, []string{""})
	})
	testFail(t, "got: []", func(m *mockingT) {
		test.NotZero(m, []string(nil))
	})
	testFail(t, "got: ", func(m *mockingT) {
		test.Nil(m, errors.New(""))
	})
	testFail(t, "got: <nil>", func(m *mockingT) {
		test.NotNil(m, nil)
	})
	testFail(t, "got: false", func(m *mockingT) {
		test.True(m, false)
	})
	testFail(t, "got: true", func(m *mockingT) {
		test.False(m, true)
	})
	testFail(t, `"World" not in "hello world"`, func(m *mockingT) {
		test.Contains(m, "hello world", "World")
	})
	testFail(t, "6 not in [1 2 3 4 5]", func(m *mockingT) {
		test.Contains(m, []int{1, 2, 3, 4, 5}, 6)
	})
	testFail(t, `"world" in "hello world"`, func(m *mockingT) {
		test.NotContains(m, "hello world", "world")
	})
	testFail(t, "3 in [1 2 3 4 5]", func(m *mockingT) {
		test.NotContains(m, []int{1, 2, 3, 4, 5}, 3)
	})
}

func TestContains(t *testing.T) {
	// Case 1: String containment - when first parameter is string
	// The second parameter is automatically converted to string for comparison
	test.Contains(t, "3.141592", "3.14")
	test.Contains(t, "3.141592", 3)    // Integer is converted to string
	test.Contains(t, "3.141592", 3.14) // Float is converted to string

	// Case 2: Custom string type compatibility
	// Contains works with custom string types (~string) in any combination
	type customString string
	test.Contains(t, customString("abc"), customString("a"))
	test.Contains(t, customString("abc"), "a")
	test.Contains(t, "abc", customString("a"))

	// Case 3: Slice element containment
	// When first parameter is a slice, Contains checks if the second parameter exists as an element
	test.Contains(t, []int{1, 2, 3, 4, 5}, 3)
	test.Contains(t, []string{"apple", "banana", "orange"}, "banana")
	test.Contains(t, []float64{1.1, 2.2, 3.3}, 2.2)
	test.Contains(t, []byte{1, 2, 3}, byte(2))

	// Case 4: Custom type slice compatibility
	// Contains works with slices of any comparable type
	type customInt int
	nums := []customInt{1, 2, 3, 4}
	test.Contains(t, nums, customInt(2))
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
