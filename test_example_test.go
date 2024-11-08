package test_test

import (
	"errors"

	"github.com/raeperd/test"
)

func Example() {
	// mock *testing.T for example purposes
	t := test.Relaxed(&mockingT{})

	test.Equal(t, "hello", "world")     // bad
	test.Equal(t, "goodbye", "goodbye") // good

	test.NotEqual(t, "hello", "world")     // good
	test.NotEqual(t, "goodbye", "goodbye") // bad

	s := []int{1, 2, 3}
	test.AllEqual(t, []int{1, 2, 3}, s) // good
	test.AllEqual(t, []int{3, 2, 1}, s) // bad

	var err error
	test.Nil(t, err)     // good
	test.NotZero(t, err) // bad
	err = errors.New("(O_o)")
	test.Nil(t, err)     // bad
	test.NotZero(t, err) // good

	type mytype string
	var mystring mytype = "hello, world"
	test.Contains(t, "world", mystring)                 // good
	test.Contains(t, "World", mystring)                 // bad
	test.NotContains(t, "\x01", []byte("\a\b\x00\r\t")) // good
	test.NotContains(t, "\x00", []byte("\a\b\x00\r\t")) // bad

	// Output:
	// want: hello; got: world
	// got: goodbye
	// want: [3 2 1]; got: [1 2 3]
	// got: <nil>
	// got: (O_o)
	// "World" not in "hello, world"
	// "\x00" in "\a\b\x00\r\t"
}
