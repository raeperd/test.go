# test  
Package test is the minimalist testing helper for Go.

Forked from [earthboundkid/be](https://github.com/earthboundkid/be), Inspired by [Mat Ryer](https://github.com/matryer/is) and [Alex Edwards](https://www.alexedwards.net/blog/easy-test-assertions-with-go-generics).

## Features

- Simple and readable test assertions using generics
- Built-in helpers for common cases like `test.NilErr` and `test.In`
- Fail fast by default but easily switch to relaxed with `test.Relaxed(t)`
- Helpers for testing against golden files with the testfile subpackage
- No dependencies: just uses standard library

## Example usage

Test for simple equality using generics:

```go
// Test two unequal strings
test.Equal(t, "hello", "world")     // bad
// t.Fatal("want: hello; got: world")
// Test two equal strings
test.Equal(t, "goodbye", "goodbye") // good
// Test equal integers, etc.
test.Equal(t, 200, resp.StatusCode)
test.Equal(t, tc.wantPtr, gotPtr)

// Test for inequality
test.Unequal(t, "hello", "world")     // good
test.Unequal(t, "goodbye", "goodbye") // bad
// t.Fatal("got: goodbye")
```

Test for equality of slices:

```go
s := []int{1, 2, 3}
test.AllEqual(t, []int{1, 2, 3}, s) // good
test.AllEqual(t, []int{3, 2, 1}, s) // bad
// t.Fatal("want: [3 2 1]; got: [1 2 3]")
```

Handle errors:

```go
var err error
test.NilErr(t, err)   // good
test.Nonzero(t, err) // bad
// t.Fatal("got: <nil>")
err = errors.New("(O_o)")
test.NilErr(t, err)   // bad
// t.Fatal("got: (O_o)")
test.Nonzero(t, err) // good
```

Check substring containment:

```go
test.In(t, "world", "hello, world") // good
test.In(t, "World", "hello, world") // bad
// t.Fatal("World" not in "hello, world")
test.NotIn(t, "\x01", []byte("\a\b\x00\r\t")) // good
test.NotIn(t, "\x00", []byte("\a\b\x00\r\t")) // bad
// t.Fatal("\x00" in "\a\b\x00\r\t")
```

Test anything else:

```go
test.True(t, o.IsValid())
test.True(t, len(pages) >= 20)
```

Test using goldenfiles:

## Philosophy
Tests usually should not fail. When they do fail, the failure should be repeatable. Therefore, it doesn't make sense to spend a lot of time writing good test messages. (This is unlike error messages, which should happen fairly often, and in production, irrepeatably.) Package test is designed to simply fail a test quickly and quietly if a condition is not met with a reference to the line number of the failing test. If the reason for having the test is not immediately clear from context, you can write a comment, just like in normal code. If you do need more extensive reporting to figure out why a test is failing, use `test.DebugLog` or `test.Debug` to capture more information.

Most tests just need simple equality testing, which is handled by `test.Equal` (for comparable types), `test.AllEqual` (for slices of comparable types), and `test.DeepEqual` (which relies on `reflect.DeepEqual`). Another common test is that a string or byte slice should contain or not some substring, which is handled by `test.In` and `test.NotIn`. Rather than package be providing every possible test helper, you are encouraged to write your own advanced helpers for use with `test.True`, while package be takes away the drudgery of writing yet another simple `func nilErr(t *testing.T, err) { ... }`.

Every tool in the be module requires a `testing.TB` as its first argument. There are various [clever ways to get the testing.TB implicitly](https://dave.cheney.net/2019/12/08/dynamically-scoped-variables-in-go), but package be is designed to be simple and explicit, so it's easiest to just always pass in a testing.TB the boring way.
