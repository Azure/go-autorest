package auth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetTokenFromAuthFile(t *testing.T) {
	os.Setenv("AZURE_AUTH_LOCATION", getCredsPath())
	auth, err := GetTokenFromAuthFile("https://management.azure.com")
	if err != nil {
		t.Fatalf("GetTokenFromAuthFile failed, got error %v", err)
	}

	if auth.BaseURI != "https://management.azure.com/" {
		t.Fatalf("auth.BaseURI not set correctly, expected 'https://management.azure.com/', got '%s'", auth.BaseURI)
	}

	expectedFile := map[string]string{
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

	if areMapsEqual(expectedFile, auth.File) == false {
		t.Fatalf("auth.File not set correctly, expected %v, got %v", expectedFile, auth.File)
	}

	if auth.Authorizer == nil {
		t.Fatalf("auth.Authorizer not set correctly, got nil")
	}
}

func getCredsPath() string {
	gopath := os.Getenv("GOPATH")
	return filepath.Join(gopath, "src", "github.com", "Azure", "go-autorest", "autorest", "azure", "auth", "mycreds.json")
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
