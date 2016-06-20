package testutil

import (
	"reflect"
	"testing"
)

// Util is a wrapper for a test instance
type Util struct {
	*testing.T
}

// Wrap wraps a testing instance
func Wrap(t *testing.T) *Util {
	return &Util{t}
}

// AssertEqual tests that the provide values are equal.
func (t *Util) AssertEqual(expected interface{}, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Error: expected '%v' got '%v'", expected, got)
	}
}

// AssertNil tests that the provided value is nil
func (t *Util) AssertNil(got interface{}) {
	if got != nil {
		t.Fatalf("Expected nil, got %v", got)
	}
}
