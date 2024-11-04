package test_test

import (
	"fmt"
	"testing"

	"github.com/raeperd/test"
)

func ExampleDebug() {
	// mock *testing.T for example purposes
	t := &mockingT{}

	// If a test fails, the callbacks will be replayed in LIFO order
	t.Run("logging-example", func(*testing.T) {
		x := 1
		x1 := x
		test.Debug(t, func() {
			// record some debug information about x1
			fmt.Println("x1:", x1)
		})
		x = 2
		x2 := x
		test.Debug(t, func() {
			// record some debug information about x2
			fmt.Println("x2:", x2)
		})
		test.Equal(t, x, 3)
	})

	// If a test succeeds, nothing will be replayed
	t.Run("silent-example", func(*testing.T) {
		y := 1
		y1 := y
		test.Debug(t, func() {
			// record some debug information about y1
			fmt.Println("y1:", y1)
		})
		y = 2
		y2 := y
		test.Debug(t, func() {
			// record some debug information about y2
			fmt.Println("y2:", y2)
		})
		test.Unequal(t, y, 3)
	})
	// Output:
	// want: 2; got: 3
	// x2: 2
	// x1: 1
}

func ExampleDebugLog() {
	// mock *testing.T for example purposes
	t := &mockingT{}

	// If a test fails, the logs will be replayed in LIFO order
	t.Run("logging-example", func(*testing.T) {
		x := 1
		test.DebugLog(t, "x: %d", x)
		x = 2
		test.DebugLog(t, "x: %d", x)
		test.Equal(t, x, 3)
	})

	// If a test succeeds, nothing will be replayed
	t.Run("silent-example", func(*testing.T) {
		y := 1
		test.DebugLog(t, "y: %d", y)
		y = 2
		test.DebugLog(t, "y: %d", y)
		test.Unequal(t, y, 3)
	})
	// Output:
	// want: 2; got: 3
	// x: 2
	// x: 1
}
