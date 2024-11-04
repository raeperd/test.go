package test_test

import "github.com/raeperd/test"

func ExampleDeepEqual() {
	// mock *testing.T for example purposes
	t := test.Relaxed(&mockingT{})

	// good
	m1 := map[int]bool{1: true, 2: false}
	m2 := map[int]bool{1: true, 2: false}
	test.DeepEqual(t, m1, m2)

	// bad
	var s1 []int
	s2 := []int{}
	test.DeepEqual(t, s1, s2) // DeepEqual is picky about nil vs. len 0

	// Output:
	// reflect.DeepEqual([]int(nil), []int{}) == false
}
