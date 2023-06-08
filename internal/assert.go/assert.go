package assert

import (
	"strings"
	"testing"
)

func Equal[T comparable](t testing.TB, actual, expected T) {
	t.Helper()

	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func StringContains(t testing.TB, actualString, expectedSubString string) {
	t.Helper()

	if !strings.Contains(actualString, expectedSubString) {
		t.Errorf("got: %q; expected to contain: %q", actualString, expectedSubString)
	}
}
