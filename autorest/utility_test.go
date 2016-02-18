package autorest

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/go-autorest/autorest/mocks"
)

const (
	jsonT = `
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

func Test_readBool(t *testing.T) {
	m := map[string]bool{
		"true":  true,
		"True":  true,
		"false": false,
		"False": false,
	}
	for s, expected := range m {
		b, err := readBool(strings.NewReader(s))
		if err != nil {
			t.Errorf("autorest: readBool returned an unexpected error -- %v", err)
		}
		if b != expected {
			t.Errorf("autorest: readBool failed to parse value -- expected %v, received %v", expected, b)
		}
	}
}

func Test_readBoolFailsWithInvalidString(t *testing.T) {
	_, err := readBool(strings.NewReader("Not a Boolean"))
	if err == nil {
		t.Errorf("autorest: readBool failed to return an error for an invalid string")
	}
}

func Test_readBoolFailsWithBadReader(t *testing.T) {
	r := mocks.NewBody("Some content")
	r.Close()
	_, err := readBool(r)
	if err == nil {
		t.Errorf("autorest: readBool failed to return an error with an invalid io.Reader")
	}
}

func Test_readFloat32(t *testing.T) {
	for _, ch := range []byte{'e', 'E', 'f', 'g', 'G'} {
		f1 := float32(123456.789)
		f2, err := readFloat32(strings.NewReader(strconv.FormatFloat(float64(f1), ch, -1, 32)))
		if err != nil {
			t.Errorf("autorest: readFloat32 returned an unexpected error for fmt %c -- %v", ch, err)
		}
		if f1 != f2 {
			t.Errorf("autorest: readFloat32 failed to parse value using fmt %v -- expected %v, received %v", ch, f2, f1)
		}
	}
}

func Test_readFloat32FailsWithInvalidString(t *testing.T) {
	_, err := readFloat32(strings.NewReader("Not a Float32"))
	if err == nil {
		t.Errorf("autorest: readFloat32 failed to return an error for an invalid string")
	}
}

func Test_readFloat32FailsWithBadReader(t *testing.T) {
	r := mocks.NewBody("Some content")
	r.Close()
	_, err := readFloat32(r)
	if err == nil {
		t.Errorf("autorest: readFloat32 failed to return an error with an invalid io.Reader")
	}
}

func Test_readFloat64(t *testing.T) {
	for _, ch := range []byte{'e', 'E', 'f', 'g', 'G'} {
		f1 := float64(123456.789)
		f2, err := readFloat64(strings.NewReader(strconv.FormatFloat(float64(f1), ch, -1, 64)))
		if err != nil {
			t.Errorf("autorest: readFloat64 returned an unexpected error for fmt %c -- %v", ch, err)
		}
		if f1 != f2 {
			t.Errorf("autorest: readFloat64 failed to parse value using fmt %v -- expected %v, received %v", ch, f2, f1)
		}
	}
}

func Test_readFloat64FailsWithInvalidString(t *testing.T) {
	_, err := readFloat64(strings.NewReader("Not a Float64"))
	if err == nil {
		t.Errorf("autorest: readFloat64 failed to return an error for an invalid string")
	}
}

func Test_readFloat64FailsWithBadReader(t *testing.T) {
	r := mocks.NewBody("Some content")
	r.Close()
	_, err := readFloat64(r)
	if err == nil {
		t.Errorf("autorest: readFloat64 failed to return an error with an invalid io.Reader")
	}
}

func Test_readInt32(t *testing.T) {
	for _, i1 := range []int32{-123, 123} {
		i2, err := readInt32(strings.NewReader(strconv.FormatInt(int64(i1), 10)))
		if err != nil {
			t.Errorf("autorest: readFloat64 returned an unexpected error for %v -- %v", i1, err)
		}
		if i1 != i2 {
			t.Errorf("autorest: readFloat64 failed to parse value -- expected %v, received %v", i2, i1)
		}
	}
}

func Test_readInt32FailsWithInvalidString(t *testing.T) {
	_, err := readInt32(strings.NewReader("Not an Integer"))
	if err == nil {
		t.Errorf("autorest: readInt32 failed to return an error for an invalid string")
	}
}

func Test_readInt32FailsWithBadReader(t *testing.T) {
	r := mocks.NewBody("Some content")
	r.Close()
	_, err := readInt32(r)
	if err == nil {
		t.Errorf("autorest: readInt32 failed to return an error with an invalid io.Reader")
	}
}

func Test_readInt64(t *testing.T) {
	for _, i1 := range []int64{-123, 123} {
		i2, err := readInt64(strings.NewReader(strconv.FormatInt(int64(i1), 10)))
		if err != nil {
			t.Errorf("autorest: readFloat64 returned an unexpected error for %v -- %v", i1, err)
		}
		if i1 != i2 {
			t.Errorf("autorest: readFloat64 failed to parse value -- expected %v, received %v", i2, i1)
		}
	}
}

func Test_readInt64FailsWithInvalidString(t *testing.T) {
	_, err := readInt64(strings.NewReader("Not an Integer"))
	if err == nil {
		t.Errorf("autorest: readInt64 failed to return an error for an invalid string")
	}
}

func Test_readInt64FailsWithBadReader(t *testing.T) {
	r := mocks.NewBody("Some content")
	r.Close()
	_, err := readInt64(r)
	if err == nil {
		t.Errorf("autorest: readInt64 failed to return an error with an invalid io.Reader")
	}
}

func Test_readString(t *testing.T) {
	s1 := "The string to return"
	s2, err := readString(strings.NewReader(s1))
	if err != nil {
		t.Errorf("autorest: readString returned an error reading %v -- %v", s1, err)
	}
	if s1 != s2 {
		t.Errorf("autorest: readString failed to read a string -- expected %v, received %v", s2, s1)
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
	return WithHeader(headerAuthorization, mocks.TestAuthorizationHeader)
}

type mockFailingAuthorizer struct{}

func (mfa mockFailingAuthorizer) WithAuthorization() PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			return r, fmt.Errorf("ERROR: mockFailingAuthorizer returned expected error")
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
