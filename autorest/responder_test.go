package autorest

import (
	"reflect"
	"testing"

	"github.com/azure/go-autorest/autorest/mocks"
)

const (
	jsonT = `
    {
      "name":"Bill Gates",
      "age":42
    }`
)

func TestCreateResponderDoesNotModify(t *testing.T) {
	r1 := mocks.NewResponse()
	r2 := mocks.NewResponse()
	p := CreateResponder()
	err := p.Respond(r1)
	if err != nil {
		t.Errorf("autorest: CreateResponder failed (%v)", err)
	}
	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("autorest: CreateResponder without decorators modified the response")
	}
}

func TestByClosing(t *testing.T) {
	r := mocks.NewResponse()
	err := Respond(r, ByClosing())
	if err != nil {
		t.Errorf("autorest: ByClosing failed (%v)", err)
	}
	if r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: ByClosing did not close the response body")
	}
}

func TestThatByUnmarshallingJsonAndClosingCloses(t *testing.T) {
	v := &mocks.T{}
	r := mocks.NewResponseWithContent(jsonT)
	err := Respond(r, ByUnmarshallingJsonAndClosing(v))
	if err != nil {
		t.Errorf("autorest: ByUnmarshallingJsonAndClosing failed (%v)", err)
	}
	if r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: ByUnmarshallingJsonAndClosing did not close the response body")
	}
}

func TestThatByUnmarhallingJsonAndClosingUnmarshals(t *testing.T) {
	v := &mocks.T{}
	r := mocks.NewResponseWithContent(jsonT)
	err := Respond(r, ByUnmarshallingJsonAndClosing(v))
	if err != nil {
		t.Errorf("autorest: ByUnmarshallingJsonAndClosing failed (%v)", err)
	}
	if v.Name != "Bill Gates" || v.Age != 42 {
		t.Errorf("autorest: ByUnmarshallingJsonAndClosing failed to properly unmarshal")
	}
}
