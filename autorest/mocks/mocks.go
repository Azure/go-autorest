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
	s             string
	b             []byte
	isOpen        bool
	closeAttempts int
}

func NewBody(s string) *Body {
	return (&Body{s: s}).reset()
}

func (body *Body) Read(b []byte) (n int, err error) {
	if !body.IsOpen() {
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
	if body.isOpen {
		body.isOpen = false
		body.closeAttempts += 1
	}
	return nil
}

func (body *Body) CloseAttempts() int {
	return body.closeAttempts
}

func (body *Body) IsOpen() bool {
	return body.isOpen
}

func (body *Body) reset() *Body {
	body.isOpen = true
	body.b = []byte(body.s)
	return body
}

type Sender struct {
	attempts      int
	content       string
	reuseResponse bool
	resp          *http.Response
	status        string
	statusCode    int
	emitErrors    int
	err           error
}

func NewSender() *Sender {
	return &Sender{status: "200 OK", statusCode: 200}
}

func (c *Sender) Do(r *http.Request) (*http.Response, error) {
	c.attempts += 1

	if !c.reuseResponse || c.resp == nil {
		resp := NewResponse()
		resp.Request = r
		resp.Body = NewBody(c.content)
		resp.Status = c.status
		resp.StatusCode = c.statusCode
		c.resp = resp
	} else {
		c.resp.Body.(*Body).reset()
	}

	if c.emitErrors > 0 || c.emitErrors < 0 {
		c.emitErrors -= 1
		if c.err == nil {
			return c.resp, fmt.Errorf("Faux Error")
		} else {
			return c.resp, c.err
		}
	} else {
		return c.resp, nil
	}
}

func (c *Sender) Attempts() int {
	return c.attempts
}

func (c *Sender) EmitErrors(emit int) {
	c.emitErrors = emit
}

func (c *Sender) SetError(err error) {
	c.err = err
}

func (c *Sender) ClearError() {
	c.SetError(nil)
}

func (c *Sender) EmitContent(s string) {
	c.content = s
}

func (c *Sender) EmitStatus(status string, code int) {
	c.status = status
	c.statusCode = code
}

func (c *Sender) ReuseResponse(reuseResponse bool) {
	c.reuseResponse = reuseResponse
}

type T struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
