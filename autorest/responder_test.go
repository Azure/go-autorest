package autorest

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/go-autorest/autorest/mocks"
)

func ExampleWithErrorUnlessOK() {
	r := mocks.NewResponse()
	r.Request = mocks.NewRequest()

	// Respond and leave the response body open (for a subsequent responder to close)
	err := Respond(r,
		WithErrorUnlessOK(),
		ByClosingIfError())

	if err == nil {
		fmt.Printf("%s of %s returned HTTP 200", r.Request.Method, r.Request.URL)

		// Complete handling the response and close the body
		Respond(r,
			ByClosing())
	}
	// Output: GET of https://microsoft.com/a/b/c/ returned HTTP 200
}

func ExampleByUnmarshallingJSON() {
	c := `
	{
		"name" : "Rob Pike",
		"age"  : 42
	}
	`

	type V struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	v := &V{}

	Respond(mocks.NewResponseWithContent(c),
		ByUnmarshallingJSON(v),
		ByClosing())

	fmt.Printf("%s is %d years old\n", v.Name, v.Age)
	// Output: Rob Pike is 42 years old
}

func ExampleByUnmarshallingXML() {
	c := `<?xml version="1.0" encoding="UTF-8"?>
	<Person>
	  <Name>Rob Pike</Name>
	  <Age>42</Age>
	</Person>`

	type V struct {
		Name string `xml:"Name"`
		Age  int    `xml:"Age"`
	}

	v := &V{}

	Respond(mocks.NewResponseWithContent(c),
		ByUnmarshallingXML(v),
		ByClosing())

	fmt.Printf("%s is %d years old\n", v.Name, v.Age)
	// Output: Rob Pike is 42 years old
}

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

func TestCreateResponderRunsDecoratorsInOrder(t *testing.T) {
	s := ""

	d := func(n int) RespondDecorator {
		return func(r Responder) Responder {
			return ResponderFunc(func(resp *http.Response) error {
				err := r.Respond(resp)
				if err == nil {
					s += fmt.Sprintf("%d", n)
				}
				return err
			})
		}
	}

	p := CreateResponder(d(1), d(2), d(3))
	err := p.Respond(&http.Response{})
	if err != nil {
		t.Errorf("autorest: Respond failed (%v)", err)
	}

	if s != "123" {
		t.Errorf("autorest: CreateResponder invoked decorators in an incorrect order; expected '123', received '%s'", s)
	}
}

func TestByIgnoring(t *testing.T) {
	r := mocks.NewResponse()

	Respond(r,
		(func() RespondDecorator {
			return func(r Responder) Responder {
				return ResponderFunc(func(r2 *http.Response) error {
					r1 := mocks.NewResponse()
					if !reflect.DeepEqual(r1, r2) {
						t.Errorf("autorest: ByIgnoring modified the HTTP Response -- received %v, expected %v", r2, r1)
					}
					return nil
				})
			}
		})(),
		ByIgnoring(),
		ByClosing())
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

func TestByClosingAcceptsNilResponse(t *testing.T) {
	r := mocks.NewResponse()

	Respond(r,
		(func() RespondDecorator {
			return func(r Responder) Responder {
				return ResponderFunc(func(resp *http.Response) error {
					resp.Body.Close()
					r.Respond(nil)
					return nil
				})
			}
		})(),
		ByClosing())
}

func TestByClosingAcceptsNilBody(t *testing.T) {
	r := mocks.NewResponse()

	Respond(r,
		(func() RespondDecorator {
			return func(r Responder) Responder {
				return ResponderFunc(func(resp *http.Response) error {
					resp.Body.Close()
					resp.Body = nil
					r.Respond(resp)
					return nil
				})
			}
		})(),
		ByClosing())
}

func TestByClosingClosesEvenAfterErrors(t *testing.T) {
	var e error

	r := mocks.NewResponse()
	Respond(r,
		withErrorRespondDecorator(&e),
		ByClosing())

	if r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: ByClosing did not close the response body after an error occurred")
	}
}

func TestByClosingClosesReturnsNestedErrors(t *testing.T) {
	var e error

	r := mocks.NewResponse()
	err := Respond(r,
		withErrorRespondDecorator(&e),
		ByClosing())

	if err == nil || !reflect.DeepEqual(e, err) {
		t.Errorf("autorest: ByClosing failed to return a nested error")
	}
}

func TestByClosingIfErrorAcceptsNilResponse(t *testing.T) {
	var e error

	r := mocks.NewResponse()

	Respond(r,
		withErrorRespondDecorator(&e),
		(func() RespondDecorator {
			return func(r Responder) Responder {
				return ResponderFunc(func(resp *http.Response) error {
					resp.Body.Close()
					r.Respond(nil)
					return nil
				})
			}
		})(),
		ByClosingIfError())
}

func TestByClosingIfErrorAcceptsNilBody(t *testing.T) {
	var e error

	r := mocks.NewResponse()

	Respond(r,
		withErrorRespondDecorator(&e),
		(func() RespondDecorator {
			return func(r Responder) Responder {
				return ResponderFunc(func(resp *http.Response) error {
					resp.Body.Close()
					resp.Body = nil
					r.Respond(resp)
					return nil
				})
			}
		})(),
		ByClosingIfError())
}

func TestByClosingIfErrorClosesIfAnErrorOccurs(t *testing.T) {
	var e error

	r := mocks.NewResponse()
	Respond(r,
		withErrorRespondDecorator(&e),
		ByClosingIfError())

	if r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: ByClosingIfError did not close the response body after an error occurred")
	}
}

func TestByClosingIfErrorDoesNotClosesIfNoErrorOccurs(t *testing.T) {
	r := mocks.NewResponse()
	Respond(r,
		ByClosingIfError())

	if !r.Body.(*mocks.Body).IsOpen() {
		t.Errorf("autorest: ByClosingIfError closed the response body even though no error occurred")
	}
}

func Test_ByUnmarshallingBool(t *testing.T) {
	m := map[string]bool{
		"true":  true,
		"True":  true,
		"false": false,
		"False": false,
	}
	for s, expected := range m {
		var b bool
		r := mocks.NewResponseWithContent(s)
		err := Respond(r,
			ByUnmarshallingBool(&b),
			ByClosing())
		if err != nil {
			t.Errorf("autorest: ByUnmarshallingBool returned an unexpected error parsing %v -- %v", s, err)
		}
		if b != expected {
			t.Errorf("autorest: ByUnmarshallingBool failed to correctly unmarshall -- expected %v, received %v", expected, b)
		}
	}
}

func Test_ByUnmarshallingBoolFailsWithInvalidString(t *testing.T) {
	var b bool
	r := mocks.NewResponseWithContent("Not a Boolean")
	err := Respond(r,
		ByUnmarshallingBool(&b),
		ByClosing())
	if err == nil {
		t.Errorf("autorest: ByUnmarshallingBool failed to return an error for an invalid string")
	}
}

func Test_ByUnmarshallingFloat32(t *testing.T) {
	for _, ch := range []byte{'e', 'E', 'f', 'g', 'G'} {
		var f1, f2 float32
		f2 = float32(123456.789)
		r := mocks.NewResponseWithContent(strconv.FormatFloat(float64(f2), ch, -1, 64))
		err := Respond(r,
			ByUnmarshallingFloat32(&f1),
			ByClosing())
		if err != nil {
			t.Errorf("autorest: ByUnmarshallingFloat32 returned an unexpected error parsing format %c -- %v", ch, err)
		}
		if f1 != f2 {
			t.Errorf("autorest: ByUnmarshallingFloat32 failed to correctly unmarshall -- expected %v, received %v", f2, f1)
		}
	}
}

func Test_ByUnmarshallingFloat32FailsWithInvalidString(t *testing.T) {
	var f float32
	r := mocks.NewResponseWithContent("Not a Float32")
	err := Respond(r,
		ByUnmarshallingFloat32(&f),
		ByClosing())
	if err == nil {
		t.Errorf("autorest: ByUnmarshallingFloat32 failed to return an error for an invalid string")
	}
}

func Test_ByUnmarshallingFloat64(t *testing.T) {
	for _, ch := range []byte{'e', 'E', 'f', 'g', 'G'} {
		var f1, f2 float64
		f2 = float64(123456.789)
		r := mocks.NewResponseWithContent(strconv.FormatFloat(float64(f2), ch, -1, 64))
		err := Respond(r,
			ByUnmarshallingFloat64(&f1),
			ByClosing())
		if err != nil {
			t.Errorf("autorest: ByUnmarshallingFloat64 returned an unexpected error parsing format %c -- %v", ch, err)
		}
		if f1 != f2 {
			t.Errorf("autorest: ByUnmarshallingFloat64 failed to correctly unmarshall -- expected %v, received %v", f2, f1)
		}
	}
}

func Test_ByUnmarshallingFloat64FailsWithInvalidString(t *testing.T) {
	var f float64
	r := mocks.NewResponseWithContent("Not a Float64")
	err := Respond(r,
		ByUnmarshallingFloat64(&f),
		ByClosing())
	if err == nil {
		t.Errorf("autorest: ByUnmarshallingFloat64 failed to return an error for an invalid string")
	}
}

func Test_ByUnmarshallingInt32(t *testing.T) {
	for _, i2 := range []int32{-123456, 123456} {
		var i1 int32
		r := mocks.NewResponseWithContent(strconv.FormatInt(int64(i2), 10))
		err := Respond(r,
			ByUnmarshallingInt32(&i1),
			ByClosing())
		if err != nil {
			t.Errorf("autorest: ByUnmarshallingInt32 returned an unexpected error -- %v", err)
		}
		if i1 != i2 {
			t.Errorf("autorest: ByUnmarshallingInt32 failed to correctly unmarshall -- expected %v, received %v", i2, i1)
		}
	}
}

func Test_ByUnmarshallingInt32FailsWithInvalidString(t *testing.T) {
	var i int32
	r := mocks.NewResponseWithContent("Not an Integer")
	err := Respond(r,
		ByUnmarshallingInt32(&i),
		ByClosing())
	if err == nil {
		t.Errorf("autorest: ByUnmarshallingInt32 failed to return an error for an invalid string")
	}
}

func Test_ByUnmarshallingInt64(t *testing.T) {
	for _, i2 := range []int64{-123456, 123456} {
		var i1 int64
		r := mocks.NewResponseWithContent(strconv.FormatInt(int64(i2), 10))
		err := Respond(r,
			ByUnmarshallingInt64(&i1),
			ByClosing())
		if err != nil {
			t.Errorf("autorest: ByUnmarshallingInt64 returned an unexpected error -- %v", err)
		}
		if i1 != i2 {
			t.Errorf("autorest: ByUnmarshallingInt64 failed to correctly unmarshall -- expected %v, received %v", i2, i1)
		}
	}
}

func Test_ByUnmarshallingInt64FailsWithInvalidString(t *testing.T) {
	var i int64
	r := mocks.NewResponseWithContent("Not an Integer")
	err := Respond(r,
		ByUnmarshallingInt64(&i),
		ByClosing())
	if err == nil {
		t.Errorf("autorest: ByUnmarshallingInt64 failed to return an error for an invalid string")
	}
}

func Test_ByUnmarshallingString(t *testing.T) {
	var s1, s2 string
	s2 = "Expected string"
	r := mocks.NewResponseWithContent(s2)
	err := Respond(r,
		ByUnmarshallingString(&s1),
		ByClosing())
	if err != nil {
		t.Errorf("autorest: ByUnmarshallingString returned an unexpected error -- %v", err)
	}
	if s1 != s2 {
		t.Errorf("autorest: ByUnmarshallingString failed to correctly unmarshall -- expected %v, received %v", s2, s1)
	}
}

func TestByUnmarshallingJSON(t *testing.T) {
	v := &mocks.T{}
	r := mocks.NewResponseWithContent(jsonT)
	err := Respond(r,
		ByUnmarshallingJSON(v),
		ByClosing())
	if err != nil {
		t.Errorf("autorest: ByUnmarshallingJSON failed (%v)", err)
	}
	if v.Name != "Rob Pike" || v.Age != 42 {
		t.Errorf("autorest: ByUnmarshallingJSON failed to properly unmarshal")
	}
}

func TestByUnmarshallingJSONIncludesJSONInErrors(t *testing.T) {
	v := &mocks.T{}
	j := jsonT[0 : len(jsonT)-2]
	r := mocks.NewResponseWithContent(j)
	err := Respond(r,
		ByUnmarshallingJSON(v),
		ByClosing())
	if err == nil || !strings.Contains(err.Error(), j) {
		t.Errorf("autorest: ByUnmarshallingJSON failed to return JSON in error (%v)", err)
	}
}

func TestByUnmarshallingXML(t *testing.T) {
	v := &mocks.T{}
	r := mocks.NewResponseWithContent(xmlT)
	err := Respond(r,
		ByUnmarshallingXML(v),
		ByClosing())
	if err != nil {
		t.Errorf("autorest: ByUnmarshallingXML failed (%v)", err)
	}
	if v.Name != "Rob Pike" || v.Age != 42 {
		t.Errorf("autorest: ByUnmarshallingXML failed to properly unmarshal")
	}
}

func TestByUnmarshallingXMLIncludesXMLInErrors(t *testing.T) {
	v := &mocks.T{}
	x := xmlT[0 : len(xmlT)-2]
	r := mocks.NewResponseWithContent(x)
	err := Respond(r,
		ByUnmarshallingXML(v),
		ByClosing())
	if err == nil || !strings.Contains(err.Error(), x) {
		t.Errorf("autorest: ByUnmarshallingXML failed to return XML in error (%v)", err)
	}
}

func TestRespondAcceptsNullResponse(t *testing.T) {
	err := Respond(nil)
	if err != nil {
		t.Errorf("autorest: Respond returned an unexpected error when given a null Response (%v)", err)
	}
}

func TestWithErrorUnlessStatusCode(t *testing.T) {
	r := mocks.NewResponse()
	r.Request = mocks.NewRequest()
	r.Status = "400 BadRequest"
	r.StatusCode = http.StatusBadRequest

	err := Respond(r,
		WithErrorUnlessStatusCode(http.StatusBadRequest, http.StatusUnauthorized, http.StatusInternalServerError),
		ByClosingIfError())

	if err != nil {
		t.Errorf("autorest: WithErrorUnlessStatusCode returned an error (%v) for an acceptable status code (%s)", err, r.Status)
	}
}

func TestWithErrorUnlessStatusCodeEmitsErrorForUnacceptableStatusCode(t *testing.T) {
	r := mocks.NewResponse()
	r.Request = mocks.NewRequest()
	r.Status = "400 BadRequest"
	r.StatusCode = http.StatusBadRequest

	err := Respond(r,
		WithErrorUnlessStatusCode(http.StatusOK, http.StatusUnauthorized, http.StatusInternalServerError),
		ByClosingIfError())

	if err == nil {
		t.Errorf("autorest: WithErrorUnlessStatusCode failed to return an error for an unacceptable status code (%s)", r.Status)
	}
}

func TestWithErrorUnlessOK(t *testing.T) {
	r := mocks.NewResponse()
	r.Request = mocks.NewRequest()

	err := Respond(r,
		WithErrorUnlessOK(),
		ByClosingIfError())

	if err != nil {
		t.Errorf("autorest: WithErrorUnlessOK returned an error for OK status code (%v)", err)
	}
}

func TestWithErrorUnlessOKEmitsErrorIfNotOK(t *testing.T) {
	r := mocks.NewResponse()
	r.Request = mocks.NewRequest()
	r.Status = "400 BadRequest"
	r.StatusCode = http.StatusBadRequest

	err := Respond(r,
		WithErrorUnlessOK(),
		ByClosingIfError())

	if err == nil {
		t.Errorf("autorest: WithErrorUnlessOK failed to return an error for a non-OK status code (%v)", err)
	}
}

func TestExtractHeader(t *testing.T) {
	r := mocks.NewResponse()
	v := []string{"v1", "v2", "v3"}
	mocks.SetResponseHeaderValues(r, mocks.TestHeader, v)

	if !reflect.DeepEqual(ExtractHeader(mocks.TestHeader, r), v) {
		t.Errorf("autorest: ExtractHeader failed to retrieve the expected header -- expected [%s]%v, received [%s]%v",
			mocks.TestHeader, v, mocks.TestHeader, ExtractHeader(mocks.TestHeader, r))
	}
}

func TestExtractHeaderHandlesMissingHeader(t *testing.T) {
	var v []string
	r := mocks.NewResponse()

	if !reflect.DeepEqual(ExtractHeader(mocks.TestHeader, r), v) {
		t.Errorf("autorest: ExtractHeader failed to handle a missing header -- expected %v, received %v",
			v, ExtractHeader(mocks.TestHeader, r))
	}
}

func TestExtractHeaderValue(t *testing.T) {
	r := mocks.NewResponse()
	v := "v1"
	mocks.SetResponseHeader(r, mocks.TestHeader, v)

	if ExtractHeaderValue(mocks.TestHeader, r) != v {
		t.Errorf("autorest: ExtractHeader failed to retrieve the expected header -- expected [%s]%v, received [%s]%v",
			mocks.TestHeader, v, mocks.TestHeader, ExtractHeaderValue(mocks.TestHeader, r))
	}
}

func TestExtractHeaderValueHandlesMissingHeader(t *testing.T) {
	r := mocks.NewResponse()
	v := ""

	if ExtractHeaderValue(mocks.TestHeader, r) != v {
		t.Errorf("autorest: ExtractHeader failed to retrieve the expected header -- expected [%s]%v, received [%s]%v",
			mocks.TestHeader, v, mocks.TestHeader, ExtractHeaderValue(mocks.TestHeader, r))
	}
}

func TestExtractHeaderValueRetrievesFirstValue(t *testing.T) {
	r := mocks.NewResponse()
	v := []string{"v1", "v2", "v3"}
	mocks.SetResponseHeaderValues(r, mocks.TestHeader, v)

	if ExtractHeaderValue(mocks.TestHeader, r) != v[0] {
		t.Errorf("autorest: ExtractHeader failed to retrieve the expected header -- expected [%s]%v, received [%s]%v",
			mocks.TestHeader, v[0], mocks.TestHeader, ExtractHeaderValue(mocks.TestHeader, r))
	}
}
