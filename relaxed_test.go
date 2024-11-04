package test_test

import (
	"testing"

	"github.com/raeperd/test"
)

func runTest(test func(*testing.T)) bool {
	return testing.RunTests(func(pat, str string) (bool, error) {
		return true, nil
	}, []testing.InternalTest{{"test", test}})
}

func TestRelaxed(t *testing.T) {
	finished := false
	test.False(t, runTest(func(t *testing.T) {
		rt := test.Relaxed(t)
		rt.FailNow()
		rt.Fatal("boom!")
		rt.Fatalf("msg: %v", "boom!")
		finished = true
	}))
	test.True(t, finished)
}
