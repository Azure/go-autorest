/*
	This package provides mocks and helpers used in testing.
*/
package mocks

import (
	"fmt"
	"io"
	"net/http"
)

type Body struct {
	s      string
	b      []byte
	isOpen bool
}

func NewBody(s string) *Body {
	return &Body{s: s, b: []byte(s), isOpen: true}
}

func (body *Body) Read(b []byte) (n int, err error) {
	if !body.isOpen {
		return 0, fmt.Errorf("ERROR: Body has been closed\n")
	}
	if len(body.b) == 0 {
		return 0, io.EOF
	}
	n = copy(b, body.b)
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

type Client struct {
	attempts   int
	content    string
	emitErrors int
	status     string
	statusCode int
	err        error
}

func NewClient() *Client {
	return &Client{status: "200 OK", statusCode: 200}
}

func (c *Client) Do(r *http.Request) (*http.Response, error) {
	c.attempts += 1

	resp := NewResponse()
	resp.Request = r
	resp.Body = NewBody(c.content)
	resp.Status = c.status
	resp.StatusCode = c.statusCode

	if c.emitErrors > 0 || c.emitErrors < 0 {
		c.emitErrors -= 1
		if c.err == nil {
			return resp, fmt.Errorf("Faux Error")
		} else {
			return resp, c.err
		}
	} else {
		return resp, nil
	}
}

func (c *Client) Attempts() int {
	return c.attempts
}

func (c *Client) EmitErrors(emit int) {
	c.emitErrors = emit
}

func (c *Client) SetError(err error) {
	c.err = err
}

func (c *Client) ClearError() {
	c.SetError(nil)
}

func (c *Client) EmitContent(s string) {
	c.content = s
}

func (c *Client) EmitStatus(status string, code int) {
	c.status = status
	c.statusCode = code
}

type T struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
