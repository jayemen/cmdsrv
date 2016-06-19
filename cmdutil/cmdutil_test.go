package cmdutil

import (
	"reflect"
	"testing"
	"time"
)

func assertEqual(t *testing.T, expected interface{}, got interface{}) {
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Error: expected '%v' got '%v'", expected, got)
	}
}

func assertNil(t *testing.T, got interface{}) {
	if got != nil {
		t.Fatalf("Expected nil, got %v", got)
	}
}

func TestExecuteCommand(t *testing.T) {
	cmd := MakeCmdCache(1*time.Second, "echo", "-n", "this is a test")
	output, err := cmd.Run()
	assertNil(t, err)
	assertEqual(t, "this is a test", string(output))
}
