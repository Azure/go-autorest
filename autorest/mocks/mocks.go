package mocks

import (
	"fmt"
)

type T struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Body struct {
	s      string
	b      []byte
	isOpen bool
}

func newBody(s string) *Body {
	return &Body{s: s, b: []byte(s), isOpen: true}
}

func (body *Body) Read(p []byte) (n int, err error) {
	if !body.isOpen {
		return 0, fmt.Errorf("ERROR: Body has been closed\n")
	}
	n = copy(p, body.b)
	body.b = body.b[n:]
	return n, nil
}

func (body *Body) Close() error {
	body.isOpen = false
	return nil
}

func (body *Body) IsOpen() bool {
	return body.isOpen
}
