package autorest

import (
	"testing"
)

func TestVersion(t *testing.T) {
	v := "7.2.3"
	if Version() != v {
		t.Fatalf("autorest: Version failed to return the expected version -- expected %s, received %s",
			v, Version())
	}
}
