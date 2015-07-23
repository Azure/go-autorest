package mocks

import (
	"fmt"
	"net/http"
)

func NewRequest() *http.Request {
	return NewRequestWithContent("")
}

func NewRequestWithContent(c string) *http.Request {
	r, _ := http.NewRequest("GET", "https://microsoft.com/a/b/c/", NewBody(c))
	return r
}

func NewRequestForURL(u string) *http.Request {
	r, err := http.NewRequest("GET", u, NewBody(""))
	if err != nil {
		panic(fmt.Sprintf("mocks: ERROR (%v) parsing testing URL %s", err, u))
	}
	return r
}

func NewResponse() *http.Response {
	return NewResponseWithContent("")
}

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

func NewResponseWithStatus(s string, c int) *http.Response {
	resp := NewResponse()
	resp.Status = s
	resp.StatusCode = c
	return resp
}

func AddResponseHeader(resp *http.Response, h string, v string) {
	if resp.Header == nil {
		resp.Header = make(http.Header)
	}
	resp.Header.Add(h, v)
}
