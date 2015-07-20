package autorest

import (
	"reflect"
	"testing"
)

// ensureValueStrings test
func TestEnsureStrings(t *testing.T) {
	m := map[string]interface{}{
		"string": "string",
		"int":    42,
		"nil":    nil,
	}
	r := map[string]string{
		"string": "string",
		"int":    "42",
		"nil":    "",
	}
	v := ensureValueStrings(m)
	if !reflect.DeepEqual(v, r) {
		t.Errorf("autorest: ensureValueStrings returned %v\n", v)
	}
}
