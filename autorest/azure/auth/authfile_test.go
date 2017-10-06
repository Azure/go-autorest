package auth

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var (
	expectedFile = map[string]string{
		"clientId":                       "client-id-123",
		"clientSecret":                   "client-secret-456",
		"subscriptionId":                 "sub-id-789",
		"tenantId":                       "tenant-id-123",
		"activeDirectoryEndpointUrl":     "https://login.microsoftonline.com",
		"resourceManagerEndpointUrl":     "https://management.azure.com/",
		"activeDirectoryGraphResourceId": "https://graph.windows.net/",
		"sqlManagementEndpointUrl":       "https://management.core.windows.net:8443/",
		"galleryEndpointUrl":             "https://gallery.azure.com/",
		"managementEndpointUrl":          "https://management.core.windows.net/",
	}
)

func TestGetTokenFromAuthFile(t *testing.T) {
	os.Setenv("AZURE_AUTH_LOCATION", filepath.Join(getCredsPath(), "credsutf16le.json"))
	auth, err := GetTokenFromAuthFile("https://management.azure.com")
	if err != nil {
		t.Logf("GetTokenFromAuthFile failed, got error %v", err)
		t.Fail()
	}

	if auth.BaseURI != "https://management.azure.com/" {
		t.Logf("auth.BaseURI not set correctly, expected 'https://management.azure.com/', got '%s'", auth.BaseURI)
		t.Fail()
	}

	if areMapsEqual(expectedFile, auth.File) == false {
		t.Logf("auth.File not set correctly, expected %v, got %v", expectedFile, auth.File)
		t.Fail()
	}

	if auth.Authorizer == nil {
		t.Log("auth.Authorizer not set correctly, got nil")
		t.Fail()
	}
}

func TestDecodeAndUnmarshal(t *testing.T) {
	tests := []string{
		"credsutf8.json",
		"credsutf16le.json",
		"credsutf16be.json",
	}
	creds := getCredsPath()
	for _, test := range tests {
		b, err := ioutil.ReadFile(filepath.Join(creds, test))
		if err != nil {
			t.Logf("error reading file '%s': %s", test, err)
			t.Fail()
		}
		decoded, err := decode(b)
		if err != nil {
			t.Logf("error decoding file '%s': %s", test, err)
			t.Fail()
		}
		var got map[string]string
		err = json.Unmarshal(decoded, &got)
		if err != nil {
			t.Logf("error unmarshaling file '%s': %s", test, err)
			t.Fail()
		}
		if areMapsEqual(expectedFile, got) == false {
			t.Logf("unmarshaled map expected %v, got %v", expectedFile, got)
			t.Fail()
		}
	}
}

func getCredsPath() string {
	gopath := os.Getenv("GOPATH")
	return filepath.Join(gopath, "src", "github.com", "Azure", "go-autorest", "testdata")
}

func areMapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k := range a {
		if a[k] != b[k] {
			return false
		}
	}
	return true
}
