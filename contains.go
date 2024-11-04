package test

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

// Contains calls t.Fatalf if needle is not contained in the string or []byte haystack.
func Contains[byteseq ~string | ~[]byte](t testing.TB, needle string, haystack byteseq) {
	t.Helper()
	if !contains(haystack, needle) {
		t.Fatalf("%q not in %q", needle, haystack)
	}
}

// NotContains calls t.Fatalf if needle is contained in the string or []byte haystack.
func NotContains[byteseq ~string | ~[]byte](t testing.TB, needle string, haystack byteseq) {
	t.Helper()
	if contains(haystack, needle) {
		t.Fatalf("%q in %q", needle, haystack)
	}
}

func contains[byteseq ~string | ~[]byte](haystack byteseq, needle string) bool {
	rv := reflect.ValueOf(haystack)
	switch rv.Kind() {
	case reflect.String:
		return strings.Contains(rv.String(), needle)
	case reflect.Slice:
		return bytes.Contains(rv.Bytes(), []byte(needle))
	default:
		panic("unreachable")
	}
}
