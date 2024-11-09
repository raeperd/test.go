package test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Equal calls t.Fatalf if want != got.
func Equal[T comparable](t testing.TB, want, got T) {
	t.Helper()
	if want != got {
		t.Fatalf("want: %v; got: %v", want, got)
	}
}

// NotEqual calls t.Fatalf if got == bad.
func NotEqual[T comparable](t testing.TB, bad, got T) {
	t.Helper()
	if got == bad {
		t.Fatalf("got: %v", got)
	}
}

// DeepEqual calls t.Fatalf if want and got are different according to reflect.DeepEqual.
func DeepEqual[T any](t testing.TB, want, got T) {
	t.Helper()
	// Pass as pointers to get around the nil-interface problem
	if !reflect.DeepEqual(&want, &got) {
		t.Fatalf("reflect.DeepEqual(%#v, %#v) == false", want, got)
	}
}

// AllEqual calls t.Fatalf if want != got.
func AllEqual[T comparable](t testing.TB, want, got []T) {
	t.Helper()
	if len(want) != len(got) {
		t.Fatalf("len(want): %d; len(got): %v", len(want), len(got))
		return
	}
	for i := range want {
		if want[i] != got[i] {
			t.Fatalf("want: %v; got: %v", want, got)
			return
		}
	}
}

// Zero calls t.Fatalf if value != the zero value for T.
func Zero[T any](t testing.TB, value T) {
	t.Helper()
	if !isZero(value) {
		t.Fatalf("got: %v", value)
	}
}

// NotZero calls t.Fatalf if value == the zero value for T.
func NotZero[T any](t testing.TB, value T) {
	t.Helper()
	if isZero(value) {
		t.Fatalf("got: %v", value)
	}
}

func isZero[T any](v T) bool {
	switch m := any(v).(type) {
	case interface{ IsZero() bool }:
		return m.IsZero()
	}

	switch rv := reflect.ValueOf(&v).Elem(); rv.Kind() {
	case reflect.Map, reflect.Slice:
		return rv.Len() == 0
	default:
		return rv.IsZero()
	}
}

// Nil calls t.Fatalf if v is not nil.
func Nil(t testing.TB, v any) {
	t.Helper()
	if v != nil {
		t.Fatalf("got: %v", v)
	}
}

// NotNil calls t.Fatalf if v is nil.
func NotNil(t testing.TB, v any) {
	t.Helper()
	if v == nil {
		t.Fatalf("got: %v", v)
	}
}

// True calls t.Fatalf if value is not true.
func True(t testing.TB, value bool) {
	t.Helper()
	if !value {
		t.Fatalf("got: false")
	}
}

// False calls t.Fatalf if value is not false.
func False(t testing.TB, value bool) {
	t.Helper()
	if value {
		t.Fatalf("got: true")
	}
}

// Contains calls t.Fatalf if needle is not contained in haystack.
// H can be either a string type (including custom string types) or a slice of comparable type N.
// When H is a string type and N is any type, fmt.Sprintf is used to convert N to string for comparison.
// When H is a slice, N is same type as the slice elements for direct comparison.
func Contains[H interface{ ~string | []N }, N comparable](t testing.TB, haystack H, needle N) {
	t.Helper()
	msg, found := containsElement(haystack, needle)
	if !found {
		t.Fatal(msg)
	}
}

// NotContains calls t.Fatalf if needle is contained in haystack.
// For type of H, N see [Contains]
func NotContains[H interface{ ~string | []N }, N comparable](t testing.TB, haystack H, needle N) {
	t.Helper()
	msg, found := containsElement(haystack, needle)
	if found {
		t.Fatal(msg)
	}
}

func containsElement[H interface{ ~string | []N }, N comparable](haystack H, needle N) (string, bool) {
	switch h := any(haystack).(type) {
	case string:
		n := fmt.Sprintf("%v", needle)
		if strings.Contains(h, n) {
			return fmt.Sprintf("%q in %q", n, h), true
		}
		return fmt.Sprintf("%q not in %q", n, h), false
	case []N:
		// TODO: refactor this using slices.Contains
		for _, v := range h {
			if v == needle {
				return fmt.Sprintf("%v in %v", needle, haystack), true
			}
		}
		return fmt.Sprintf("%v not in %v", needle, haystack), false
	default: // h is custom string type
		hs := fmt.Sprintf("%v", haystack)
		n := fmt.Sprintf("%v", needle)
		if strings.Contains(hs, n) {
			return fmt.Sprintf("%q in %q", n, hs), true
		}
		return fmt.Sprintf("%q not in %q", n, hs), false
	}
}
