package autorest

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/azure/go-autorest/autorest/mocks"
)

// PrepareDecorators wrap and invoke a Preparer. The decorator may invoke the passed Preparer
// before, after, or within its own processing.
func ExamplePrepareDecorator(path string) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			if r.URL == nil {
				return r, fmt.Errorf("ERROR: URL is not set")
			}
			r.URL.Path += path
			return p.Prepare(r)
		})
	}
}

func ExamplePrepareDecorator_post(path string) PrepareDecorator {
	return func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r, err := p.Prepare(r)
			if err != nil {
				return r, err
			}
			r.Header.Add(http.CanonicalHeaderKey("ContentType"), "application/json")
			return r, nil
		})
	}
}

// Create a sequence of three Preparers that build up the URL path.
func ExampleCreatePreparer() {
	p := CreatePreparer(WithBaseURL("https://microsoft.com/"), WithPath("/a/b/c/"))
	r, err := p.Prepare(&http.Request{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/
}

// Create, and then chain, separate Preparers
func ExampleCreatePreparer_chained() {
	params := map[string]interface{}{
		"param1": "a",
		"param2": "c",
	}

	p1 := CreatePreparer(WithBaseURL("https://microsoft.com/"), WithPath("/{param1}/b/{param2}/"))
	p2 := CreatePreparer(WithPathParameters(params))

	r, err := p1.Prepare(&http.Request{})
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}

	r, err = p2.Prepare(r)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/
}

// Create and prepare an http.Request in one call
func ExamplePrepare() {
	r, err := Prepare(&http.Request{},
		AsGet(),
		WithBaseURL("https://microsoft.com/"),
		WithPath("a/b/c/"))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("%s %s", r.Method, r.URL)
	}
	// Output: GET https://microsoft.com/a/b/c/
}

// Create a request for a supplied base URL and path
func ExampleWithBaseURL() {
	r, err := Prepare(&http.Request{},
		WithBaseURL("https://microsoft.com/a/b/c/"))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/
}

// Create a request with a custom HTTP header
func ExampleWithHeader() {
	r, err := Prepare(&http.Request{},
		WithBaseURL("https://microsoft.com/a/b/c/"),
		WithHeader("x-foo", "bar"))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Printf("Header %s=%s\n", "x-foo", r.Header.Get("x-foo"))
	}
	// Output: Header x-foo=bar
}

// Create a request from a path with parameters
func ExampleWithPathParameters() {
	params := map[string]interface{}{
		"param1": "a",
		"param2": "c",
	}
	r, err := Prepare(&http.Request{},
		WithBaseURL("https://microsoft.com/"),
		WithPath("/{param1}/b/{param2}/"),
		WithPathParameters(params))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/
}

// Create a request with query parameters
func ExampleWithQueryParameters() {
	params := map[string]interface{}{
		"q1": "value1",
		"q2": "value2",
	}
	r, err := Prepare(&http.Request{},
		WithBaseURL("https://microsoft.com/"),
		WithPath("/a/b/c/"),
		WithQueryParameters(params))
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/?q1=value1&q2=value2
}

func TestCreatePreparerDoesNotModify(t *testing.T) {
	r1 := &http.Request{}
	p := CreatePreparer()
	r2, err := p.Prepare(r1)
	if err != nil {
		t.Errorf("autorest: CreatePreparer failed (%v)", err)
	}
	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("autorest: CreatePreparer without decorators modified the request")
	}
}

func TestCreatePreparerRunsDecoratorsInOrder(t *testing.T) {
	p := CreatePreparer(WithBaseURL("https://microsoft.com/"), WithPath("1"), WithPath("2"), WithPath("3"))
	r, err := p.Prepare(&http.Request{})
	if err != nil {
		t.Errorf("autorest: CreatePreparer failed (%v)", err)
	}
	if r.URL.String() != "https://microsoft.com/1/2/3" {
		t.Errorf("autorest: CreatePreparer failed to run decorators in order")
	}
}

func TestAsJson(t *testing.T) {
	r, err := mocks.NewRequest()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	r, err = Prepare(r, AsJson())
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	if r.Header.Get(headerContentType) != mimeTypeJson {
		t.Errorf("autorest: WithBearerAuthorization failed to add header (%s=%s)", "Authorization", r.Header.Get("Authorization"))
	}
}

func TestWithBearerAuthorization(t *testing.T) {
	r, err := mocks.NewRequest()
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	r, err = Prepare(r, WithBearerAuthorization("SOME-TOKEN"))
	if err != nil {
		fmt.Printf("ERROR: %v", err)
	}
	if r.Header.Get(headerAuthorization) != "Bearer SOME-TOKEN" {
		t.Errorf("autorest: WithBearerAuthorization failed to add header (%s=%s)", headerAuthorization, r.Header.Get(headerAuthorization))
	}
}

func TestWithHeaderAllocatesHeaders(t *testing.T) {
	r, err := Prepare(&http.Request{}, WithHeader("x-foo", "bar"))
	if err != nil {
		t.Errorf("autorest: WithHeader failed (%v)", err)
	}
	if r.Header.Get("x-foo") != "bar" {
		t.Errorf("autorest: WithHeader failed to add header (%s=%s)", "x-foo", r.Header.Get("x-foo"))
	}
}

func TestModifyingExistingRequest(t *testing.T) {
	r, err := http.NewRequest("GET", "https://bing.com/", nil)
	if err != nil {
		t.Errorf("autorest: Creating a new request failed (%v)", err)
	}
	r, err = Prepare(r, WithPath("search"), WithQueryParameters(map[string]interface{}{"q": "golang"}))
	if err != nil {
		t.Errorf("autorest: Preparing an existing request returned an error (%v)", err)
	}
	if r.URL.String() != "https://bing.com/search?q=golang" {
		t.Errorf("autorest: Preparing an existing request failed (%s)", r.URL)
	}
}
