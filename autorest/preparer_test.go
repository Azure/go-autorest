package autorest

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// Create a chain of three Preparers that append segments to the URL path.
func ExampleCreatePreparer() {
	u, _ := url.Parse("https://microsoft.com/")
	req := &http.Request{URL: u}

	d1 := func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.URL.Path += "a/"
			return p.Prepare(r)
		})
	}
	d2 := func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.URL.Path += "b/"
			return p.Prepare(r)
		})
	}
	d3 := func(p Preparer) Preparer {
		return PreparerFunc(func(r *http.Request) (*http.Request, error) {
			r.URL.Path += "c/"
			return p.Prepare(r)
		})
	}

	r, err := CreatePreparer(d1, d2, d3).Prepare(req)
	if err != nil {
		fmt.Println("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/
}

// Create a request from a path with parameters
func ExampleCreatePreparerWithPathAndParameters() {
	params := map[string]interface{}{
		"param1": "a",
		"param2": "c",
	}
	r, err := Prepare(&http.Request{},
		WithURL("https://microsoft.com/", "/{param1}/b/{param2}/"),
		WithPathParameters(params))
	if err != nil {
		fmt.Println("ERROR: %v\n", err)
	} else {
		fmt.Println(r.URL)
	}
	// Output: https://microsoft.com/a/b/c/
}

// Ensure the default Preparer does not modify the http.Request
func TestCreatePreparer(t *testing.T) {
	r1 := &http.Request{}
	p := CreatePreparer()
	r2, err := p.Prepare(r1)
	if err != nil {
		t.Errorf("autorest: CreatePreparer failed (%v)\n", err)
	}
	if !reflect.DeepEqual(r1, r2) {
		t.Errorf("autorest: CreatePreparer without decorators modified the request\n")
	}
}
