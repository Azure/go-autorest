package autorest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/Azure/go-autorest/autorest/mocks"
)

const (
	testAuthorizationHeader = "BEARER SECRETTOKEN"
	testBadURL              = ""
	jsonT                   = `
    {
      "name":"Rob Pike",
      "age":42
    }`
	xmlT = `<?xml version="1.0" encoding="UTF-8"?>
	<Person>
		<Name>Rob Pike</Name>
		<Age>42</Age>
	</Person>`
)

func TestNewDecoderCreatesJSONDecoder(t *testing.T) {
	d := NewDecoder(EncodedAsJSON, strings.NewReader(jsonT))
	_, ok := d.(*json.Decoder)
	if d == nil || !ok {
		t.Error("autorest: NewDecoder failed to create a JSON decoder when requested")
	}
}

func TestNewDecoderCreatesXMLDecoder(t *testing.T) {
	d := NewDecoder(EncodedAsXML, strings.NewReader(xmlT))
	_, ok := d.(*xml.Decoder)
	if d == nil || !ok {
		t.Error("autorest: NewDecoder failed to create an XML decoder when requested")
	}
}

func TestNewDecoderReturnsNilForUnknownEncoding(t *testing.T) {
	d := NewDecoder("unknown", strings.NewReader(xmlT))
	if d != nil {
		t.Error("autorest: NewDecoder created a decoder for an unknown encoding")
	}
}

func TestCopyAndDecodeDecodesJSON(t *testing.T) {
	_, err := CopyAndDecode(EncodedAsJSON, strings.NewReader(jsonT), &mocks.T{})
	if err != nil {
		t.Errorf("autorest: CopyAndDecode returned an error with valid JSON - %v", err)
	}
}

func TestCopyAndDecodeDecodesXML(t *testing.T) {
	_, err := CopyAndDecode(EncodedAsXML, strings.NewReader(xmlT), &mocks.T{})
	if err != nil {
		t.Errorf("autorest: CopyAndDecode returned an error with valid XML - %v", err)
	}
}

func TestCopyAndDecodeReturnsJSONDecodingErrors(t *testing.T) {
	_, err := CopyAndDecode(EncodedAsJSON, strings.NewReader(jsonT[0:len(jsonT)-2]), &mocks.T{})
	if err == nil {
		t.Errorf("autorest: CopyAndDecode failed to return an error with invalid JSON")
	}
}

func TestCopyAndDecodeReturnsXMLDecodingErrors(t *testing.T) {
	_, err := CopyAndDecode(EncodedAsXML, strings.NewReader(xmlT[0:len(xmlT)-2]), &mocks.T{})
	if err == nil {
		t.Errorf("autorest: CopyAndDecode failed to return an error with invalid XML")
	}
}

func TestCopyAndDecodeAlwaysReturnsACopy(t *testing.T) {
	b, _ := CopyAndDecode(EncodedAsJSON, strings.NewReader(jsonT), &mocks.T{})
	if b.String() != jsonT {
		t.Errorf("autorest: CopyAndDecode failed to return a valid copy of the data - %v", b.String())
	}
}

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

func doEnsureBodyClosed(t *testing.T) SendDecorator {
	return func(s Sender) Sender {
		return SenderFunc(func(r *http.Request) (*http.Response, error) {
			resp, err := s.Do(r)
			if resp != nil && resp.Body != nil && resp.Body.(*mocks.Body).IsOpen() {
				t.Error("autorest: Expected Body to be closed -- it was left open")
			}
			return resp, err
		})
	}
}

type mockAuthorizer struct{}

func (ma mockAuthorizer) WithAuthorization() PrepareDecorator {
	return WithHeader(headerAuthorization, testAuthorizationHeader)
}

type mockFailingAuthorizer struct{}

func (mfa mockFailingAuthorizer) WithAuthorization() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			return r, fmt.Errorf("ERROR: mockFailingAuthorizer returned expected error")
		})
	}
}

func withMessage(output *string, msg string) SendDecorator {
	return func(s Sender) Sender {
		return SenderFunc(func(r *http.Request) (*http.Response, error) {
			resp, err := s.Do(r)
			if err == nil {
				*output += msg
			}
			return resp, err
		})
	}
}

type mockInspector struct {
	wasInvoked bool
}

func (mi *mockInspector) WithInspection() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			mi.wasInvoked = true
			return p.Prepare(r)
		})
	}
}

func (mi *mockInspector) ByInspecting() RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			mi.wasInvoked = true
			return r.Respond(resp)
		})
	}
}

func withErrorRespondDecorator(e *error) RespondDecorator {
	return func(r Responder) Responder {
		return ResponderFunc(func(resp *http.Response) error {
			err := r.Respond(resp)
			if err != nil {
				return err
			}
			*e = fmt.Errorf("autorest: Faux Respond Error")
			return *e
		})
	}
}
