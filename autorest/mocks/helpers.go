package mocks

import (
	"fmt"
	"net/http"
)

// NewRequest instantiates a new request.
func NewRequest() *http.Request {
	return NewRequestWithContent("")
}

// NewRequestWithContent instantiates a new request using the passed string for the body content.
func NewRequestWithContent(c string) *http.Request {
	r, _ := http.NewRequest("GET", "https://microsoft.com/a/b/c/", NewBody(c))
	return r
}

// NewRequestForURL instantiates a new request using the passed URL.
func NewRequestForURL(u string) *http.Request {
	r, err := http.NewRequest("GET", u, NewBody(""))
	if err != nil {
		panic(fmt.Sprintf("mocks: ERROR (%v) parsing testing URL %s", err, u))
	}
	return r
}

// NewResponse instantiates a new response.
func NewResponse() *http.Response {
	return NewResponseWithContent("")
}

// NewResponseWithContent instantiates a new response with the passed string as the body content.
func NewResponseWithContent(c string) *http.Response {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       NewBody(c),
		Request:    NewRequest(),
	}
}

// NewResponseWithStatus instantiates a new response using the passed string and integer as the
// status and status code.
func NewResponseWithStatus(s string, c int) *http.Response {
	resp := NewResponse()
	resp.Status = s
	resp.StatusCode = c
	return resp
}

// SetResponseHeader adds a header to the passed response.
func SetResponseHeader(resp *http.Response, h string, v string) {
	if resp.Header == nil {
		resp.Header = make(http.Header)
	}
	resp.Header.Set(h, v)
}
