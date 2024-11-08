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

	test.DeepEqual(t, map[int]bool{1: true, 2: false}, map[int]bool{1: true, 2: false}) // good
	test.DeepEqual(t, nil, []int{})                                                     // bad

	s := []int{1, 2, 3}
	test.AllEqual(t, []int{1, 2, 3}, s) // good
	test.AllEqual(t, []int{3, 2, 1}, s) // bad

	var err error
	test.Nil(t, err)    // good
	test.NotNil(t, err) // bad
	err = errors.New("(O_o)")
	test.Nil(t, err)    // bad
	test.NotNil(t, err) // good

	type mytype string
	var mystring mytype = "hello, world"
	test.Contains(t, mystring, "world") // good
	test.Contains(t, mystring, "World") // bad

	test.NotContains(t, mystring, "World") // good
	test.NotContains(t, mystring, "world") // bad

	// Output:
	// want: hello; got: world
	// got: goodbye
	// reflect.DeepEqual([]int(nil), []int{}) == false
	// want: [3 2 1]; got: [1 2 3]
	// got: <nil>
	// got: (O_o)
	// "World" not in "hello, world"
	// "world" in "hello, world"
}
