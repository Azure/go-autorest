package autorest

import (
	"net/http"
	"net/url"
	"testing"
)

func TestSasNewSasAuthorizerEmptyToken(t *testing.T) {
	auth, err := NewSASTokenAuthorizer("")
	if err == nil {
		t.Fatalf("azure: SASTokenAuthorizer#NewSASTokenAuthorizer didn't return an error")
	}

	if auth != nil {
		t.Fatalf("azure: SASTokenAuthorizer#NewSASTokenAuthorizer returned an authorizer")
	}
}

func TestSasNewSasAuthorizerEmptyTokenWithWhitespace(t *testing.T) {
	auth, err := NewSASTokenAuthorizer("  ")
	if err == nil {
		t.Fatalf("azure: SASTokenAuthorizer#NewSASTokenAuthorizer didn't return an error")
	}

	if auth != nil {
		t.Fatalf("azure: SASTokenAuthorizer#NewSASTokenAuthorizer returned an authorizer")
	}
}

func TestSasNewSasAuthorizerValidToken(t *testing.T) {
	auth, err := NewSASTokenAuthorizer("abc123")
	if err != nil {
		t.Fatalf("azure: SASTokenAuthorizer#NewSASTokenAuthorizer returned an error")
	}

	if auth == nil {
		t.Fatalf("azure: SASTokenAuthorizer#NewSASTokenAuthorizer didn't return an authorizer")
	}
}

func TestSasAuthorizerRequest(t *testing.T) {
	testData := []struct{
		name string
		token string
		input string
		expected string
	}{
		{
			name: "empty querystring without a prefix",
			token: "abc123",
			input: "https://example.com/foo/bar",
			expected: "https://example.com/foo/bar?abc123",
		},
		{
			name: "empty querystring with a prefix",
			token: "?abc123",
			input: "https://example.com/foo/bar",
			expected: "https://example.com/foo/bar?abc123",
		},
		{
			name: "existing querystring without a prefix",
			token: "abc123",
			input: "https://example.com/foo/bar?hello=world",
			expected: "https://example.com/foo/bar?hello=world&abc123",
		},
		{
			name: "existing querystring with a prefix",
			token: "?abc123",
			input: "https://example.com/foo/bar?hello=world",
			expected: "https://example.com/foo/bar?hello=world&abc123",
		},
	}

	for _, v := range testData {
		t.Logf("[DEBUG] Testing Case %q..", v.name)
		auth, err := NewSASTokenAuthorizer(v.token)
		if err != nil {
			t.Fatalf("azure: SASTokenAuthorizer#WithAuthorization expected %q but got an error", v.expected)
		}
		url, _ := url.ParseRequestURI(v.input)
		httpReq := &http.Request{
			URL: url,
		}

		req, err := Prepare(httpReq, auth.WithAuthorization())
		if err != nil {
			t.Fatalf("azure: SASTokenAuthorizer#WithAuthorization returned an error (%v)", err)
		}

		if req.RequestURI != v.expected {
			t.Fatalf("azure: SASTokenAuthorizer#WithAuthorization failed to set QueryString header - got %q but expected %q", req.RequestURI, v.expected)
		}

		if req.Header.Get(http.CanonicalHeaderKey("Authorization")) != "" {
			t.Fatal("azure: SASTokenAuthorizer#WithAuthorization set an Authorization header when it shouldn't!")
		}
	}
}