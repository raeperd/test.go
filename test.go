package test

import (
	"reflect"
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
