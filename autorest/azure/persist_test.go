package azure

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"
)

const MockTokenJSON string = `{
	"access_token": "accessToken",
	"refresh_token": "refreshToken",
	"expires_in": "1000",
	"expires_on": "2000",
	"not_before": "3000",
	"resource": "resource",
	"token_type": "type"
}`

var TestToken = Token{
	AccessToken:  "accessToken",
	RefreshToken: "refreshToken",
	ExpiresIn:    "1000",
	ExpiresOn:    "2000",
	NotBefore:    "3000",
	Resource:     "resource",
	Type:         "type",
}

func TestLoadToken(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "testloadtoken")
	if err != nil {
		t.Errorf("unexpected error when creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(MockTokenJSON))
	if err != nil {
		t.Errorf("unexpected error when writing temp test file: %v", err)
	}

	expectedToken := TestToken
	actualToken, err := LoadToken(f.Name())
	if err != nil {
		t.Errorf("unexpected error loading token from file: %v", err)
	}

	if *actualToken != expectedToken {
		t.Errorf("failed to decode properly actual(%v) expected(%v)", *actualToken, expectedToken)
	}
}

func TestLoadTokenFailsBadPath(t *testing.T) {
	_, err := LoadToken("/tmp/this_file_should_never_exist_really")
	if err == nil {
		// expected failed to open file...
		t.Errorf("failed to get error")
	}
}

func TestLoadTokenFailsBadJson(t *testing.T) {
	gibberishJSON := strings.Replace(MockTokenJSON, "expires_on", ";:\"gibberish", -1)

	f, err := ioutil.TempFile(os.TempDir(), "testloadtokenfailsbadjson")
	if err != nil {
		t.Errorf("unexpected error when creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	_, err = f.Write([]byte(gibberishJSON))
	if err != nil {
		t.Errorf("unexpected error when writing temp test file: %v", err)
	}

	_, err = LoadToken(f.Name())
	if err == nil {
		// expected failed to decode contents of file...
		t.Errorf("failed to get error")
	}
}

func token() *Token {
	var token Token
	json.Unmarshal([]byte(MockTokenJSON), &token)
	return &token
}

func TestSaveToken(t *testing.T) {
	f, err := ioutil.TempFile(os.TempDir(), "testloadtoken")
	if err != nil {
		t.Errorf("unexpected error when creating temp file: %v", err)
	}
	defer os.Remove(f.Name())

	err = SaveToken(f.Name(), *token())
	if err != nil {
		t.Errorf("unexpected error saving token to file: %v", err)
	}

	var actualToken Token
	var expectedToken Token

	json.Unmarshal([]byte(MockTokenJSON), expectedToken)

	contents, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Errorf("!!")
	}
	json.Unmarshal(contents, actualToken)

	if !reflect.DeepEqual(actualToken, expectedToken) {
		t.Errorf("token was not serialized correctly")
	}
}

func TestSaveTokenFailsNoPermission(t *testing.T) {
	err := SaveToken("/usr/thiswontwork/atall", *token())
	if err == nil {
		// expected failed to decode contents of file...
		t.Errorf("failed to get error")
	}
}

func TestSaveTokenFailsCantCreate(t *testing.T) {
	err := SaveToken("/thiswontwork", *token())
	if err == nil {
		// expected failed to decode contents of file...
		t.Errorf("failed to get error")
	}
}
