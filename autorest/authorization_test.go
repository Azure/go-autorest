package autorest

import (
	"github.com/Azure/go-autorest/autorest/mocks"
	"reflect"
	"testing"
)

func TestWithAuthorizer(t *testing.T) {
	r1 := mocks.NewRequest()

	na := &NullAuthorizer{}
	r2, err := Prepare(r1,
		na.WithAuthorization())
	if err != nil {
		t.Fatalf("autorest: NullAuthorizer#WithAuthorization returned an unexpected error (%v)", err)
	} else if !reflect.DeepEqual(r1, r2) {
		t.Fatalf("autorest: NullAuthorizer#WithAuthorization modified the request -- received %v, expected %v", r2, r1)
	}
}
