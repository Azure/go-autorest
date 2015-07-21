package mocks

import (
	"net/http"
)

func NewRequest() (*http.Request, error) {
	return NewRequestWithContent("")
}

func NewRequestWithContent(c string) (*http.Request, error) {
	return http.NewRequest("GET", "https://microsoft.com/a/b/c/", newBody(c))
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
		Body:       newBody(c),
	}
}
