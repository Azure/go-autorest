package autorest

import (
	"reflect"
	"testing"
)

func TestContainsIntFindsValue(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := 5
	if !containsInt(ints, v) {
		t.Errorf("autorest: containsInt failed to find %v in %v", v, ints)
	}
}

func TestContainsIntDoesNotFindValue(t *testing.T) {
	ints := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	v := 42
	if containsInt(ints, v) {
		t.Errorf("autorest: containsInt unexpectedly found %v in %v", v, ints)
	}
}

func TestEscapeStrings(t *testing.T) {
	m := map[string]string{
		"string": "a long string with = odd characters",
		"int":    "42",
		"nil":    "",
	}
	r := map[string]string{
		"string": "a+long+string+with+%3D+odd+characters",
		"int":    "42",
		"nil":    "",
	}
	v := escapeValueStrings(m)
	if !reflect.DeepEqual(v, r) {
		t.Errorf("autorest: ensureValueStrings returned %v\n", v)
	}
}

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
