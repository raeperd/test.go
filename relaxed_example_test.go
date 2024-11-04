package test_test

import (
	"testing"

	"github.com/raeperd/test"
)

func ExampleRelaxed() {
	// mock *testing.T for example purposes
	t := &mockingT{}

	t.Run("dies on first error", func(*testing.T) {
		test.Equal(t, 1, 2)
		test.Equal(t, 3, 4)
	})

	t.Run("shows multiple errors", func(*testing.T) {
		relaxedT := test.Relaxed(t)
		test.Equal(relaxedT, 5, 6)
		test.Equal(relaxedT, 7, 8)
	})
	// Output:
	// want: 1; got: 2
	// want: 5; got: 6
	// want: 7; got: 8
}
